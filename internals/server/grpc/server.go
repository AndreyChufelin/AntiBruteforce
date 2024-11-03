package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/ratelimiter"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	pbratelimter.UnimplementedRatelimiterServer
	logger  *slog.Logger
	server  *grpc.Server
	limiter *ratelimiter.Limiter
	port    string
}

func NewGRPC(logger *slog.Logger, limiter *ratelimiter.Limiter, port string) *Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor(logger)))
	return &Server{
		logger:  logger,
		server:  grpcServer,
		limiter: limiter,
		port:    port,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed start tcp server: %w", err)
	}

	s.logger.Info("grpc server started", slog.String("addr", l.Addr().String()))
	pbratelimter.RegisterRatelimiterServer(s.server, s)

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stoping grpc server")
	done := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("grpc server stopped gracefully")
		return nil
	case <-ctx.Done():
		s.logger.Warn("context done, forcing server stop")
		s.server.Stop()
		return fmt.Errorf("stop operation canceled: %w", ctx.Err())
	}
}

func (s *Server) Allow(ctx context.Context, request *pbratelimter.AllowRequest) (*pbratelimter.AllowResponse, error) {
	ok, err := s.limiter.ReqAllowed(ctx, request.Login, request.Password, request.Ip)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pbratelimter.AllowResponse{Ok: ok}, nil
}

func LoggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		date := time.Now()
		resp, err := handler(ctx, req)

		latency := time.Since(date)

		p, ok := peer.FromContext(ctx)
		ip := "unknown"
		if ok {
			ip = p.Addr.String()
		}

		md, ok := metadata.FromIncomingContext(ctx)
		userAgent := "unknown"
		if ok {
			if userAgents, exists := md["user-agent"]; exists && len(userAgents) > 0 {
				userAgent = userAgents[0]
			}
		}

		logger.Info("GRPC request handled",
			"ip", ip,
			"method", info.FullMethod,
			"date", date.Format(time.RFC822Z),
			"userAgent", userAgent,
			"latency", latency.Milliseconds(),
			"status", status.Code(err).String(),
		)
		return resp, err
	}
}

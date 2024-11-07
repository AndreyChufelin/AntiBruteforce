package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Server struct {
	pbratelimter.UnimplementedRatelimiterServer
	pbiplist.UnimplementedIPListServiceServer
	logger  *slog.Logger
	server  *grpc.Server
	limiter Limiter
	iplist  IPList
	port    string
}

type IPList interface {
	WhitelistAdd(ctx context.Context, ip string) error
	WhitelistDelete(ctx context.Context, ip string) error
	BlacklistAdd(ctx context.Context, ip string) error
	BlacklistDelete(ctx context.Context, ip string) error
}

type Limiter interface {
	ReqAllowed(ctx context.Context, login, password, ip string) (bool, error)
	ClearReq(ctx context.Context, login, ip string) error
}

func NewGRPC(logger *slog.Logger, limiter Limiter, iplist IPList, port string) *Server {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor(logger)))
	return &Server{
		logger:  logger,
		server:  grpcServer,
		limiter: limiter,
		port:    port,
		iplist:  iplist,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed start tcp server: %w", err)
	}

	s.logger.Info("grpc server started", slog.String("addr", l.Addr().String()))
	pbratelimter.RegisterRatelimiterServer(s.server, s)
	pbiplist.RegisterIPListServiceServer(s.server, s)

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

func validateIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("wrong ip")
	}
	return nil
}

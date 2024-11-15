package grpcserver

import (
	"context"
	"errors"

	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Allow(ctx context.Context, request *pbratelimter.AllowRequest) (*pbratelimter.AllowResponse, error) {
	logg := s.logger.With("handler", "Allow")
	err := validateIP(request.Ip)
	if err != nil {
		logg.Warn("invalid argument", "ip", request.Ip)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ok, err := s.limiter.ReqAllowed(ctx, request.Login, request.Password, request.Ip)
	if err != nil {
		logg.Error("failed to allow request", "login", request.Login, "ip", request.Ip, "err", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pbratelimter.AllowResponse{Ok: ok}, nil
}

func (s *Server) Clear(ctx context.Context, request *pbratelimter.ClearRequest) (*pbratelimter.Empty, error) {
	logg := s.logger.With("handler", "Clear")

	if request.Ip != "" {
		err := validateIP(request.Ip)
		if err != nil {
			logg.Warn("invalid argument", "ip", request.Ip)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	if request.Login != "" {
		err := s.limiter.ClearReq(ctx, storage.LoginBucket, request.Login)
		if err != nil {
			if errors.Is(err, storage.ErrBucketNotExist) {
				logg.Warn("failed to clear bucket: login doesn't exist", "login", request.Login)
				return nil, status.Error(codes.NotFound, "no bucket with this login")
			}

			logg.Error("failed to clear login request", "login", request.Login, "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	if request.Ip != "" {
		err := s.limiter.ClearReq(ctx, storage.IPBucket, request.Ip)
		if err != nil {
			if errors.Is(err, storage.ErrBucketNotExist) {
				logg.Warn("failed to clear bucket: ip doesn't exist", "ip", request.Ip)
				return nil, status.Error(codes.NotFound, "no bucket with this ip")
			}

			logg.Error("failed to clear ip request", "ip", request.Ip, "err", err)
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &pbratelimter.Empty{}, nil
}

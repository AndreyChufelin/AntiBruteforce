package grpcserver

import (
	"context"
	"errors"

	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
	pbiplist "github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) WhitelistAdd(ctx context.Context, request *pbiplist.ListRequest) (*pbiplist.Empty, error) {
	logg := s.logger.With("handler", "WhitelistAdd")
	err := validateSubnet(request.Subnet)
	if err != nil {
		logg.Warn("invalid argument", "subnet", request.Subnet)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.iplist.WhitelistAdd(ctx, request.Subnet)
	if err != nil {
		logg.Error("failed to add to whitelist", "err", err)
		return nil, status.Error(codes.Internal, "intenal server error")
	}
	return &pbiplist.Empty{}, nil
}

func (s *Server) WhitelistDelete(ctx context.Context, request *pbiplist.ListRequest) (*pbiplist.Empty, error) {
	logg := s.logger.With("handler", "WhitelistDelete")
	err := validateSubnet(request.Subnet)
	if err != nil {
		logg.Warn("invalid argument", "subnet", request.Subnet)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.iplist.WhitelistDelete(ctx, request.Subnet)
	if err != nil {
		if errors.Is(err, storage.ErrSubnetNotExist) {
			logg.Warn("subnet doesn't exist")
			return nil, status.Error(codes.NotFound, "subnet doesn't exist")
		}
		logg.Error("failed to delete from whitelist", "err", err)
		return nil, status.Error(codes.Internal, "intenal server error")
	}
	return &pbiplist.Empty{}, nil
}

func (s *Server) BlacklistAdd(ctx context.Context, request *pbiplist.ListRequest) (*pbiplist.Empty, error) {
	logg := s.logger.With("handler", "BlacklistAdd")
	err := validateSubnet(request.Subnet)
	if err != nil {
		logg.Warn("invalid argument", "subnet", request.Subnet)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.iplist.BlacklistAdd(ctx, request.Subnet)
	if err != nil {
		logg.Error("failed to add to blacklist", "err", err)
		return nil, status.Error(codes.Internal, "intenal server error")
	}
	return &pbiplist.Empty{}, nil
}

func (s *Server) BlacklistDelete(ctx context.Context, request *pbiplist.ListRequest) (*pbiplist.Empty, error) {
	logg := s.logger.With("handler", "BlacklistDelete")
	err := validateSubnet(request.Subnet)
	if err != nil {
		logg.Warn("invalid argument", "subnet", request.Subnet)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.iplist.BlacklistDelete(ctx, request.Subnet)
	if err != nil {
		if errors.Is(err, storage.ErrSubnetNotExist) {
			logg.Warn("subnet doesn't exist")
			return nil, status.Error(codes.NotFound, "subnet doesn't exist")
		}
		logg.Error("failed to delete from blacklist", "err", err)
		return nil, status.Error(codes.Internal, "intenal server error")
	}
	return &pbiplist.Empty{}, nil
}

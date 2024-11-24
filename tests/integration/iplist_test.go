//go:build integration

package integration

import (
	"github.com/AndreyChufelin/AntiBruteforce/pb/iplist"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *IntegrationSuite) TestWhitelistAdd() {
	subnet := "192.168.1.0/24"
	_, err := s.handlers.WhitelistAdd(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().NoError(err)

	var exists bool
	err = s.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM whitelist WHERE subnet=$1)", subnet)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *IntegrationSuite) TestWhitelistAddInvalidSubnet() {
	subnet := "192.168.1"
	_, err := s.handlers.WhitelistAdd(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "wrong subnet"))
}

func (s *IntegrationSuite) TestWhitelistDelete() {
	subnet := "192.168.1.0/24"
	_, err := s.db.Exec("INSERT INTO whitelist (subnet) VALUES ($1)", subnet)
	s.Require().NoError(err)

	_, err = s.handlers.WhitelistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().NoError(err)

	var exists bool
	err = s.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM whitelist WHERE subnet=$1)", subnet)
	s.Require().NoError(err)
	s.Require().False(exists)
}

func (s *IntegrationSuite) TestWhitelistDeleteNotExist() {
	subnet := "192.168.1.0/24"
	_, err := s.handlers.WhitelistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.NotFound, "subnet doesn't exist"))
}

func (s *IntegrationSuite) TestWhitelistDeleteInvalidSubnet() {
	subnet := "192.168.1"
	_, err := s.handlers.WhitelistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "wrong subnet"))
}

func (s *IntegrationSuite) TestBlacklistAdd() {
	subnet := "192.168.1.0/24"
	_, err := s.handlers.BlacklistAdd(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().NoError(err)

	var exists bool
	err = s.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM blacklist WHERE subnet=$1)", subnet)
	s.Require().NoError(err)
	s.Require().True(exists)
}

func (s *IntegrationSuite) TestBlacklistDeleteNotExist() {
	subnet := "192.168.1.0/24"
	_, err := s.handlers.BlacklistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.NotFound, "subnet doesn't exist"))
}

func (s *IntegrationSuite) TestBlacklistAddInvalidSubnet() {
	subnet := "192.168.1"
	_, err := s.handlers.BlacklistAdd(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "wrong subnet"))
}

func (s *IntegrationSuite) TestBlacklistDelete() {
	subnet := "192.168.1.0/24"
	_, err := s.db.Exec("INSERT INTO blacklist (subnet) VALUES ($1)", subnet)
	s.Require().NoError(err)

	_, err = s.handlers.BlacklistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().NoError(err)

	var exists bool
	err = s.db.Get(&exists, "SELECT EXISTS(SELECT 1 FROM blacklist WHERE subnet=$1)", subnet)
	s.Require().NoError(err)
	s.Require().False(exists)
}

func (s *IntegrationSuite) TestBlacklistDeleteInvalidSubnet() {
	subnet := "192.168.1"
	_, err := s.handlers.BlacklistDelete(testSuite.ctx, &iplist.ListRequest{Subnet: subnet})
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "wrong subnet"))
}

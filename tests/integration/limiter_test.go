//go:build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/iplist"
	"github.com/AndreyChufelin/AntiBruteforce/internals/ratelimiter"
	grpcserver "github.com/AndreyChufelin/AntiBruteforce/internals/server/grpc"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage/postgres"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage/redis"
	pbratelimter "github.com/AndreyChufelin/AntiBruteforce/pb/ratelimiter"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	redisdb "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	suite.Suite
	db              *sqlx.DB
	config          Config
	hanlders        *grpcserver.Server
	redis           *redisdb.Client
	limiterInterval time.Duration
}

var (
	testSuite IntegrationSuite
)

func TestMain(m *testing.M) {
	var err error
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	testSuite.config, err = LoadConfig("./config.toml")
	if err != nil {
		log.Fatal("failed loading config", err)
	}

	testSuite.db, err = sqlx.Connect("postgres",
		fmt.Sprintf(
			"user=%s dbname=%s sslmode=disable password=%s host=%s port=%s",
			testSuite.config.DB.User,
			testSuite.config.DB.Name,
			testSuite.config.DB.Password,
			testSuite.config.DB.Host,
			testSuite.config.DB.Port,
		),
	)
	redisAddr := net.JoinHostPort(testSuite.config.Redis.Host, testSuite.config.Redis.Port)
	testSuite.redis = redisdb.NewClient(&redisdb.Options{
		Addr:     redisAddr,
		Password: testSuite.config.Redis.Password,
		DB:       testSuite.config.Redis.DB,
	})
	if err := testSuite.redis.Ping(context.TODO()).Err(); err != nil {
		log.Fatal("failed to connect to redis", err)
	}

	redis := redis.NewRedis(testSuite.config.Redis.Host, testSuite.config.Redis.Port, testSuite.config.Redis.Password, testSuite.config.Redis.DB)
	if err := redis.Connect(context.TODO()); err != nil {
		logger.Error("Failed to start redis", "err", err)
		os.Exit(1)
	}

	postgres := postgres.New(testSuite.config.DB.User, testSuite.config.DB.Password, testSuite.config.DB.Name, testSuite.config.DB.Host, testSuite.config.DB.Port)
	if err := postgres.Connect(context.TODO()); err != nil {
		logger.Error("Failed to start postgres", "err", err)
		os.Exit(1)
	}

	iplist := iplist.NewIPList(logger, postgres)

	testSuite.limiterInterval = time.Duration(testSuite.config.Limiter.Interval) * time.Second
	limiter := ratelimiter.NewRateLimiter(logger, redis, ratelimiter.Options{
		Login:    testSuite.config.Limiter.Login,
		Password: testSuite.config.Limiter.Password,
		IP:       testSuite.config.Limiter.IP,
		Interval: testSuite.limiterInterval,
	}, iplist)

	testSuite.hanlders = grpcserver.NewGRPC(logger, limiter, iplist, testSuite.config.GRPC.Port)

	code := m.Run()

	testSuite.db.Close()

	os.Exit(code)
}

func (s *IntegrationSuite) TearDownTest() {
	_, err := s.db.Exec("TRUNCATE TABLE whitelist, blacklist")
	if err != nil {
		log.Fatal("failed to clear postgres", err)
	}
	err = s.redis.FlushAll(context.TODO()).Err()
	if err != nil {
		log.Fatal("falied to clear redis", err)
	}
}

func (s *IntegrationSuite) TestAllowLogin() {
	for i := range s.config.Limiter.Login {
		password := fmt.Sprintf("pass%d", i)
		ip := fmt.Sprintf("127.0.0.%d", i)
		res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: password, Ip: ip})
		s.Require().NoError(err)
		s.Require().True(res.Ok, fmt.Sprintf("Request #%d", i))
	}

	res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().False(res.Ok)

	time.Sleep(s.limiterInterval)

	res, err = s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().True(res.Ok)
}

func (s *IntegrationSuite) TestAllowPassword() {
	for i := range s.config.Limiter.Password {
		login := fmt.Sprintf("user%d", i)
		ip := fmt.Sprintf("127.0.0.%d", i)
		res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: login, Password: "123456", Ip: ip})
		s.Require().NoError(err)
		s.Require().True(res.Ok, fmt.Sprintf("Request #%d", i))
	}

	res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().False(res.Ok)

	time.Sleep(s.limiterInterval)

	res, err = s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().True(res.Ok)
}

func (s *IntegrationSuite) TestAllowIP() {
	for i := range s.config.Limiter.Password {
		login := fmt.Sprintf("user%d", i)
		password := fmt.Sprintf("pass%d", i)
		res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: login, Password: password, Ip: "127.0.0.1"})
		s.Require().NoError(err)
		s.Require().True(res.Ok, fmt.Sprintf("Request #%d", i))
	}

	res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().False(res.Ok)

	time.Sleep(s.limiterInterval)

	res, err = s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().True(res.Ok)
}

func (s *IntegrationSuite) TestAllowWhitelist() {
	s.db.Exec("INSERT INTO whitelist (subnet) VALUES ($1)", "127.0.0.0/8")
	for range s.config.Limiter.Login + 1 {
		res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
		s.Require().NoError(err)
		s.Require().True(res.Ok)
	}
}

func (s *IntegrationSuite) TestAllowBlacklist() {
	s.db.Exec("INSERT INTO blacklist (subnet) VALUES ($1)", "127.0.0.0/8")
	res, err := s.hanlders.Allow(context.TODO(), &pbratelimter.AllowRequest{Login: "user", Password: "123456", Ip: "127.0.0.1"})
	s.Require().NoError(err)
	s.Require().False(res.Ok)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &testSuite)
}

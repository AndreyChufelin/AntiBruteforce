package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	conf "github.com/AndreyChufelin/AntiBruteforce/internals/config"
	"github.com/AndreyChufelin/AntiBruteforce/internals/iplist"
	"github.com/AndreyChufelin/AntiBruteforce/internals/ratelimiter"
	grpcserver "github.com/AndreyChufelin/AntiBruteforce/internals/server/grpc"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage/postgres"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage/redis"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := conf.LoadConfig(configFile)
	if err != nil {
		log.Fatal("failed loading config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	redis := redis.NewRedis(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB)
	if err := redis.Connect(ctx); err != nil {
		logger.Error("Failed to start redis", "err", err)
		os.Exit(1)
	}

	postgres := postgres.New(config.DB.User, config.DB.Password, config.DB.Name, config.DB.Host, config.DB.Port)
	if err := postgres.Connect(ctx); err != nil {
		logger.Error("Failed to start postgres", "err", err)
		os.Exit(1)
	}

	iplist := iplist.NewIPList(logger, postgres)

	limiter := ratelimiter.NewRateLimiter(logger, redis, ratelimiter.Options{
		Login:    config.Limiter.LoginLimit,
		Password: config.Limiter.PasswordLimit,
		IP:       config.Limiter.IPLimit,
		Interval: time.Duration(config.Limiter.Interval) * time.Second,
	}, iplist)

	server := grpcserver.NewGRPC(logger, limiter, iplist, config.GRPC.Port)

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Failed to start grpc server", "err", err)
			os.Exit(1)
		}
	}()

	defer cancel()
	<-ctx.Done()
	logger.Info("Stopping service")

	ctxStop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := redis.Close(ctxStop); err != nil {
		logger.Error("Failed to stop redis", "err", err)
	}
	if err := postgres.Close(ctxStop); err != nil {
		logger.Error("Failed to stop postgres", "err", err)
	}
	if err := server.Stop(ctxStop); err != nil {
		logger.Error("Failed to stop grpc server", "err", err)
	}
}

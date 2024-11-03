package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/ratelimiter"
	grpcserver "github.com/AndreyChufelin/AntiBruteforce/internals/server/grpc"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage/redis"
)

func main() {
	config, err := LoadConfig("/etc/antibruteforce/config.toml")
	if err != nil {
		log.Fatal("failed loading config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	redis := redis.NewRedis(config.Redis.Host, config.Redis.Port, config.Redis.Password, config.Redis.DB)
	if err := redis.Start(ctx); err != nil {
		logger.Error("Failed to start redis", "err", err)
		os.Exit(1)
	}

	limiter := ratelimiter.NewRateLimiter(logger, redis, ratelimiter.Rates{
		Login:    config.Rates.Login,
		Password: config.Rates.Password,
		IP:       config.Rates.IP,
	})

	server := grpcserver.NewGRPC(logger, limiter, config.GRPC.Port)

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

	if err := redis.Stop(ctxStop); err != nil {
		logger.Error("Failed to stop redis", "err", err)
	}
	if err := server.Stop(ctxStop); err != nil {
		logger.Error("Failed to stop grpc server", "err", err)
	}
}

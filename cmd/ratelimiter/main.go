package main

import (
	"context"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	redis := redis.NewRedis()
	if err := redis.Start(ctx); err != nil {
		logger.Error("Failed to start redis", "err", err)
		os.Exit(1)
	}

	limiter := ratelimiter.NewRateLimiter(logger, redis)

	server := grpcserver.NewGRPC(logger, limiter)

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

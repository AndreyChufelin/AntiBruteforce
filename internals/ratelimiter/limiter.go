package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
)

type Limiter struct {
	storage Storage
	logger  *slog.Logger
	rates   Rates
}

type Rates struct {
	Login    int
	Password int
	IP       int
}

//go:generate mockery --name Storage
type Storage interface {
	UpdateBucket(ctx context.Context, bucketType storage.BucketType, key string, limit int, period time.Duration) error
}

func NewRateLimiter(logger *slog.Logger, storage Storage, rates Rates) *Limiter {
	return &Limiter{
		storage: storage,
		logger:  logger,
		rates:   rates,
	}
}

func (r *Limiter) ReqAllowed(ctx context.Context, login, password, ip string) (bool, error) {
	logg := r.logger.With("op", "ReqAllowed")
	err := r.storage.UpdateBucket(ctx, storage.LoginBucket, login, r.rates.Login, time.Minute)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by login", "login", login)
			return false, nil
		}
		logg.Error("failed to update login bucket", "login", login, "err", err)
		return false, fmt.Errorf("failed to update login bucket %s: %w", login, err)
	}

	err = r.storage.UpdateBucket(ctx, storage.PasswordBucket, password, r.rates.Password, time.Minute)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by password")
			return false, nil
		}
		logg.Error("failed to update password bucket", "err", err)
		return false, fmt.Errorf("failed to update password bucket: %w", err)
	}

	err = r.storage.UpdateBucket(ctx, storage.IPBucket, ip, r.rates.IP, time.Minute)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by ip", "ip", ip)
			return false, nil
		}
		logg.Error("failed to update ip bucket", "ip", ip, "err", err)
		return false, fmt.Errorf("failed to update ip bucket %s: %w", ip, err)
	}

	return true, nil
}

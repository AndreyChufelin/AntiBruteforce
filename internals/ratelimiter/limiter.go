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
	iplist  IPList
	logger  *slog.Logger
	options Options
}

type Options struct {
	Login    int
	Password int
	IP       int
	Interval time.Duration
}

//go:generate mockery --name Storage
type Storage interface {
	UpdateBucket(ctx context.Context, bucketType storage.BucketType, key string, limit int, period time.Duration) error
	ClearBucket(ctx context.Context, bucketType storage.BucketType, key string) error
}

//go:generate mockery --name IPList
type IPList interface {
	WhitelistCheckSubnet(ctx context.Context, ip string) (bool, error)
	BlacklistCheckSubnet(ctx context.Context, ip string) (bool, error)
}

func NewRateLimiter(logger *slog.Logger, storage Storage, options Options, iplist IPList) *Limiter {
	return &Limiter{
		storage: storage,
		logger:  logger,
		options: options,
		iplist:  iplist,
	}
}

func (l *Limiter) ReqAllowed(ctx context.Context, login, password, ip string) (bool, error) {
	logg := l.logger.With("op", "ReqAllowed")
	ipstatus, err := l.isIPAllowed(ctx, ip)
	if err != nil {
		logg.Error("failed to check ip", "ip", ip, "err", err)
		return false, fmt.Errorf("failed to check ip %s: %w", ip, err)
	}
	if ipstatus == ipAllowed {
		logg.Info("request allowed from whitelist", "ip", ip)
		return true, nil
	}
	if ipstatus == ipRejected {
		logg.Info("request blocked from blacklist", "ip", ip)
		return false, nil
	}

	err = l.storage.UpdateBucket(ctx, storage.LoginBucket, login, l.options.Login, l.options.Interval)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by login", "login", login)
			return false, nil
		}
		logg.Error("failed to update login bucket", "login", login, "err", err)
		return false, fmt.Errorf("failed to update login bucket %s: %w", login, err)
	}

	err = l.storage.UpdateBucket(ctx, storage.PasswordBucket, password, l.options.Password, l.options.Interval)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by password")
			return false, nil
		}
		logg.Error("failed to update password bucket", "err", err)
		return false, fmt.Errorf("failed to update password bucket: %w", err)
	}

	err = l.storage.UpdateBucket(ctx, storage.IPBucket, ip, l.options.IP, l.options.Interval)
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

type ipstatus int

const (
	ipAllowed ipstatus = iota
	ipRejected
	ipNone
)

func (l *Limiter) isIPAllowed(ctx context.Context, ip string) (ipstatus, error) {
	inWhitelist, err := l.iplist.WhitelistCheckSubnet(ctx, ip)
	if err != nil {
		return ipNone, fmt.Errorf("failed to check whitelist: %w", err)
	}
	if inWhitelist {
		return ipAllowed, nil
	}
	inBlacklist, err := l.iplist.BlacklistCheckSubnet(ctx, ip)
	if err != nil {
		return ipNone, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if inBlacklist {
		return ipRejected, nil
	}

	return ipNone, nil
}

func (l *Limiter) ClearReq(ctx context.Context, bucketType storage.BucketType, key string) error {
	logg := l.logger.With("op", "ClearReq")

	err := l.storage.ClearBucket(ctx, bucketType, key)
	if err != nil {
		logg.Error("failed to clear bucket", "type", bucketType, "key", key, "err", err)
		return fmt.Errorf("failed to clear bucket: %w", err)
	}

	return nil
}

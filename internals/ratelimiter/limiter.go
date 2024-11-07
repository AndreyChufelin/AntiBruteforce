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
	ClearBucket(ctx context.Context, bucketType storage.BucketType, key string) error
}

//go:generate mockery --name IPList
type IPList interface {
	WhitelistCheckIP(ctx context.Context, ip string) (bool, error)
	BlacklistCheckIP(ctx context.Context, ip string) (bool, error)
}

func NewRateLimiter(logger *slog.Logger, storage Storage, rates Rates, iplist IPList) *Limiter {
	return &Limiter{
		storage: storage,
		logger:  logger,
		rates:   rates,
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

	err = l.storage.UpdateBucket(ctx, storage.LoginBucket, login, l.rates.Login, time.Minute)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by login", "login", login)
			return false, nil
		}
		logg.Error("failed to update login bucket", "login", login, "err", err)
		return false, fmt.Errorf("failed to update login bucket %s: %w", login, err)
	}

	err = l.storage.UpdateBucket(ctx, storage.PasswordBucket, password, l.rates.Password, time.Minute)
	if err != nil {
		if errors.Is(err, storage.ErrBucketFull) {
			logg.Warn("request rejected by password")
			return false, nil
		}
		logg.Error("failed to update password bucket", "err", err)
		return false, fmt.Errorf("failed to update password bucket: %w", err)
	}

	err = l.storage.UpdateBucket(ctx, storage.IPBucket, ip, l.rates.IP, time.Minute)
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
	inWhitelist, err := l.iplist.WhitelistCheckIP(ctx, ip)
	if err != nil {
		return ipNone, fmt.Errorf("failed to check whitelist: %w", err)
	}
	if inWhitelist {
		return ipAllowed, nil
	}
	inBlacklist, err := l.iplist.BlacklistCheckIP(ctx, ip)
	if err != nil {
		return ipNone, fmt.Errorf("failed to check blacklist: %w", err)
	}
	if inBlacklist {
		return ipRejected, nil
	}

	return ipNone, nil
}

func (l *Limiter) ClearReq(ctx context.Context, login, ip string) error {
	logg := l.logger.With("op", "ClearReq")
	if login != "" {
		err := l.storage.ClearBucket(ctx, storage.LoginBucket, login)
		if err != nil {
			logg.Error("failed to clear bucket", "login", login, "err", err)
			return fmt.Errorf("failed to clear bucket: %w", err)
		}
	}
	if ip != "" {
		err := l.storage.ClearBucket(ctx, storage.IPBucket, ip)
		if err != nil {
			logg.Error("failed to clear bucket", "ip", ip, "err", err)
			return fmt.Errorf("failed to clear bucket: %w", err)
		}
	}

	return nil
}

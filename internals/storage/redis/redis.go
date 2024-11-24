package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	client   *redis.Client
	addr     string
	password string
	db       int
}

func NewRedis(host, port, password string, db int) *Storage {
	return &Storage{
		addr:     net.JoinHostPort(host, port),
		password: password,
		db:       db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	s.client = redis.NewClient(&redis.Options{
		Addr:     s.addr,
		Password: s.password,
		DB:       s.db,
	})

	if err := s.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.client == nil {
		return fmt.Errorf("no client initialized")
	}

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		errCh <- s.client.Close()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to close client: %w", err)
		}
		return nil
	case <-ctx.Done():
		return fmt.Errorf("stop operation canceled: %w", ctx.Err())
	}
}

func (s *Storage) UpdateBucket(
	ctx context.Context,
	bucketType storage.BucketType,
	value string, limit int,
	period time.Duration,
) error {
	key := fmt.Sprintf("%s:%s", bucketType, value)
	txf := func(tx *redis.Tx) error {
		rTime, err := tx.Time(ctx).Result()
		if err != nil {
			return fmt.Errorf("failed to get redis time: %w", err)
		}
		now := rTime.UnixNano()

		requestCost := period.Nanoseconds() / int64(limit)
		ttl := period.Nanoseconds()

		err = tx.SetNX(ctx, key, 0, period).Err()
		if err != nil {
			return fmt.Errorf("failed check if bucket exists: %w", err)
		}

		tat, err := tx.Get(ctx, key).Int64()
		if err != nil {
			return fmt.Errorf("failed to get bucket: %w", err)
		}

		newTat := tat + requestCost
		if now > tat {
			newTat = now + requestCost
		}

		timeUntilNextRequest := newTat - now
		if timeUntilNextRequest > ttl {
			return storage.ErrBucketFull
		}

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			err := pipe.Set(ctx, key, newTat, period).Err()
			return err
		})
		return err
	}

	err := retry(ctx, 10, 10*time.Millisecond, func() error {
		return s.client.Watch(ctx, txf, key)
	})
	if err != nil {
		return err
	}

	return nil
}

func retry(ctx context.Context, maxRetries int, baseDelay time.Duration, fn func() error) error {
	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		if !errors.Is(err, redis.TxFailedErr) {
			return err
		}

		delay := baseDelay * time.Duration(1<<i)
		if delay > time.Second {
			delay = time.Second
		}
		select {
		case <-time.After(delay):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return fmt.Errorf("reached maximum retries")
}

func (s *Storage) ClearBucket(ctx context.Context, bucketType storage.BucketType, value string) error {
	key := fmt.Sprintf("%s:%s", bucketType, value)
	result, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}
	if result == 0 {
		return storage.ErrBucketNotExist
	}

	return nil
}

package redis

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestUpdateBucket(t *testing.T) {
	rdb := &Storage{}
	s := miniredis.RunT(t)
	rdb.client = redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	key := "user"
	limit := 10
	period := time.Minute
	for i := range limit {
		err := rdb.UpdateBucket(context.Background(), storage.LoginBucket, key, limit, period)
		require.NoError(t, err, fmt.Sprintf("call #%d", i+1))
	}
}

func TestUpdateBucketTooManyCalls(t *testing.T) {
	rdb := &Storage{}
	s := miniredis.RunT(t)
	rdb.client = redis.NewClient(&redis.Options{
		Addr:     s.Addr(),
		Password: "",
		DB:       0,
	})

	key := "user"
	limit := 10
	period := time.Second
	for i := range limit {
		err := rdb.UpdateBucket(context.Background(), storage.LoginBucket, key, limit, period)
		require.NoError(t, err, fmt.Sprintf("call #%d", i+1))
	}
	err := rdb.UpdateBucket(context.Background(), storage.LoginBucket, key, limit, period)
	require.ErrorIs(t, err, storage.ErrBucketFull)

	// Reset bucket by setting `tat` to now
	setKey := fmt.Sprintf("%s:%s", storage.LoginBucket, key)
	now := strconv.FormatInt(time.Now().UnixNano(), 10)
	s.Set(setKey, now)

	// Verify the bucket accepts a new request
	err = rdb.UpdateBucket(context.Background(), storage.LoginBucket, key, limit, period)
	require.NoError(t, err)
}

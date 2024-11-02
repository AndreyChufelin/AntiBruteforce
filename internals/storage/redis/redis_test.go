package redis

import (
	"context"
	"fmt"
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

	key := "login:user"
	limit := 10
	period := time.Minute
	for i := range limit {
		err := rdb.UpdateBucket(context.TODO(), storage.LoginBucket, key, limit, period)
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

	key := "login:user"
	limit := 10
	period := time.Second
	for i := range limit {
		err := rdb.UpdateBucket(context.TODO(), storage.LoginBucket, key, limit, period)
		require.NoError(t, err, fmt.Sprintf("call #%d", i+1))
	}
	err := rdb.UpdateBucket(context.TODO(), storage.LoginBucket, key, limit, period)
	require.ErrorIs(t, err, storage.ErrBucketFull)

	time.Sleep(period / time.Duration(limit))

	err = rdb.UpdateBucket(context.TODO(), storage.LoginBucket, key, limit, period)
	require.NoError(t, err)
}

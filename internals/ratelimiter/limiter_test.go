package ratelimiter

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/AndreyChufelin/AntiBruteforce/internals/ratelimiter/mocks"
	"github.com/AndreyChufelin/AntiBruteforce/internals/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func SetupLimiter(tb testing.TB, storage Storage, iplist IPList) *Limiter {
	tb.Helper()
	b := &bytes.Buffer{}
	logger := slog.New(slog.NewJSONHandler(b, nil))
	return NewRateLimiter(logger, storage, Options{Login: 10, Password: 100, IP: 1000, Interval: time.Minute}, iplist)
}

func TestReqAllowed(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockStorage.On("UpdateBucket", mock.Anything, storage.LoginBucket, "user", 10, time.Minute).Return(nil)
	mockStorage.On("UpdateBucket", mock.Anything, storage.PasswordBucket, "password", 100, time.Minute).Return(nil)
	mockStorage.On("UpdateBucket", mock.Anything, storage.IPBucket, "127.0.0.1", 1000, time.Minute).Return(nil)

	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)
	mockIPList.On("BlacklistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.True(t, allowed)
	require.NoError(t, err)
}

func TestReqAllowedInWhitelist(t *testing.T) {
	mockStorage := mocks.NewStorage(t)

	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(true, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.True(t, allowed)
	require.NoError(t, err)
}

func TestReqAllowedInBlacklist(t *testing.T) {
	mockStorage := mocks.NewStorage(t)

	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)
	mockIPList.On("BlacklistCheckSubnet", mock.Anything, "127.0.0.1").Return(true, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.False(t, allowed)
	require.NoError(t, err)
}

func TestReqAllowedLoginTooMany(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockStorage.On("UpdateBucket", mock.Anything, storage.LoginBucket, "user", 10, time.Minute).
		Return(storage.ErrBucketFull)
	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)
	mockIPList.On("BlacklistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.False(t, allowed)
	require.NoError(t, err)
}

func TestReqAllowedPasswordTooMany(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockStorage.On("UpdateBucket", mock.Anything, storage.LoginBucket, "user", 10, time.Minute).Return(nil)
	mockStorage.On("UpdateBucket", mock.Anything, storage.PasswordBucket, "password", 100, time.Minute).
		Return(storage.ErrBucketFull)
	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)
	mockIPList.On("BlacklistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.False(t, allowed)
	require.NoError(t, err)
}

func TestReqAllowedIPTooMany(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockStorage.On("UpdateBucket", mock.Anything, storage.LoginBucket, "user", 10, time.Minute).Return(nil)
	mockStorage.On("UpdateBucket", mock.Anything, storage.PasswordBucket, "password", 100, time.Minute).Return(nil)
	mockStorage.On("UpdateBucket", mock.Anything, storage.IPBucket, "127.0.0.1", 1000, time.Minute).
		Return(storage.ErrBucketFull)
	mockIPList := mocks.NewIPList(t)
	mockIPList.On("WhitelistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)
	mockIPList.On("BlacklistCheckSubnet", mock.Anything, "127.0.0.1").Return(false, nil)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	allowed, err := limiter.ReqAllowed(context.Background(), "user", "password", "127.0.0.1")
	require.False(t, allowed)
	require.NoError(t, err)
}

func TestClearReq(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockStorage.On("ClearBucket", mock.Anything, storage.LoginBucket, "user").Return(nil)
	mockIPList := mocks.NewIPList(t)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	err := limiter.ClearReq(context.Background(), storage.LoginBucket, "user")
	require.NoError(t, err)
}

func TestClearReqError(t *testing.T) {
	mockStorage := mocks.NewStorage(t)
	mockErr := fmt.Errorf("not exist")
	mockStorage.On("ClearBucket", mock.Anything, storage.LoginBucket, "user").Return(mockErr)
	mockIPList := mocks.NewIPList(t)

	limiter := SetupLimiter(t, mockStorage, mockIPList)

	err := limiter.ClearReq(context.Background(), storage.LoginBucket, "user")
	require.ErrorIs(t, err, mockErr)
}

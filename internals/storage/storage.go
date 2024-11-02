package storage

import "errors"

type BucketType string

const (
	LoginBucket    BucketType = "login"
	PasswordBucket BucketType = "password"
	IPBucket       BucketType = "ip"
)

var ErrBucketFull = errors.New("bucket is full")

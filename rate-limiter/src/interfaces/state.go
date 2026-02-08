package interfaces

import "time"

type LimiterState struct {
	TokenBucket TokenBucketState
	LeakyBucket LeakyBucketState
}

type TokenBucketState struct {
	Tokens         int
	LastRefillTime time.Time
}

type LeakyBucketState struct {
	Water        int
	LastLeakTime time.Time
}

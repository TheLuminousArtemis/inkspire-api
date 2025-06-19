package ratelimiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, ip string) (bool, time.Duration, error)
}

type Config struct {
	RequestsPerTimeFrame int
	Timeframe            time.Duration
	Enabled              bool
}

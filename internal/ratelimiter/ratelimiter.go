package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	Timeframe            time.Duration
	Enabled              bool
}

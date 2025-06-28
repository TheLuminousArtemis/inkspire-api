package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/inkspire/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
	RedisRateLimit interface {
		// GetCount(ctx context.Context, key string) (int, error)
		Increment(ctx context.Context, key string) (int, error)
		TTL(ctx context.Context, key string) (time.Duration, error)
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users:          &UserRedisStorage{rdb},
		RedisRateLimit: &RateLimitRedisStore{rdb},
	}
}

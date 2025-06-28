package cache

import (
	"context"
	"time"

	"github.com/theluminousartemis/inkspire/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users:          &MockUserStore{},
		RedisRateLimit: &MockRateLimitStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	return &store.User{}, nil
}

func (m *MockUserStore) Set(ctx context.Context, user *store.User) error {
	return nil
}

type MockRateLimitStore struct {
	count int
}

// func (m *MockRateLimitStore) GetCount(ctx context.Context, key string) (int, error) {
// 	return m.count, nil
// }

func (m *MockRateLimitStore) Increment(ctx context.Context, key string) (int, error) {
	m.count++
	return m.count, nil
}

func (m *MockRateLimitStore) TTL(ctx context.Context, key string) (time.Duration, error) {
	return time.Second, nil
}

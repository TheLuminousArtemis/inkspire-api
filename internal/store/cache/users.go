package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/theluminousartemis/socialnews/internal/store"
)

type UserRedisStorage struct {
	rdb *redis.Client
}

var UserTimeExp time.Duration = 24 * time.Hour

func (r *UserRedisStorage) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%v", userID)
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}
func (r *UserRedisStorage) Set(ctx context.Context, user *store.User) error {
	if user.ID == 0 {
		return errors.New("ID must be set for user to be stored in cache")
	}
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	r.rdb.SetEx(ctx, cacheKey, json, UserTimeExp).Err()
	return nil
}

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/domain"
)

// Define Cache interface
type Repository interface {
	Set(key string, user *domain.User) error
	Get(key string) (*domain.User, error)
}

// Define struct to implement the Cache interface
type repository struct {
	client  *redis.Client
	expires time.Duration
}

// Define constructor to inject redis properties
func NewRedis(host string, db int, expires time.Duration) Repository {
	return &repository{
		client: redis.NewClient(&redis.Options{
			Addr:     host,
			Password: "",
			DB:       db,
		}),
		expires: expires,
	}
}

// Set github user details if it does not exist in the cache
func (r *repository) Set(key string, user *domain.User) error {
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Set expiration
	ctx := context.Background()
	r.client.Set(ctx, key, json, r.expires*time.Minute)
	return nil
}

// Get github user details if it exists in the cache
func (r *repository) Get(key string) (*domain.User, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return &domain.User{}, nil
	}

	// Marshal value to user
	var user *domain.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

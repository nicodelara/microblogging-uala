package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisCacheRepository crea una nueva instancia del repositorio de caché
func NewRedisCacheRepository(client *redis.Client, ttl time.Duration) *redisCacheRepository {
	return &redisCacheRepository{
		client: client,
		ttl:    ttl,
	}
}

func (r *redisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisCacheRepository) Set(ctx context.Context, key string, value string) error {
	return r.client.Set(ctx, key, value, r.ttl).Err()
}

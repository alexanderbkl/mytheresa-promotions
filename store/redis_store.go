// store/redis_store.go
package store

import (
	"context"
	"encoding/json"

	"github.com/alexanderbkl/mytheresa-promotions/models"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		Client: client,
	}
}

func (s *RedisStore) GetProducts(ctx context.Context) ([]models.Product, error) {
	keys, err := s.Client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for _, key := range keys {
		val, err := s.Client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var product models.Product
		err = json.Unmarshal([]byte(val), &product)
		if err != nil {
			continue
		}

		products = append(products, product)
	}
	return products, nil
}

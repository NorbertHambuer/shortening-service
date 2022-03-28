package cache

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisCache struct{
	Client *redis.Client
}

// NewRedisCache creates a new redis client and returns a new *RedisCache that contains the client
func NewRedisCache(addr, port, pass string) (*RedisCache, error){
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pass,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{Client: client}, nil
}

// SetShortUrl saves a short url code and url into the cache
func (c *RedisCache) SetShortUrl(code, url string) error{
	return c.Client.Set(code, url, 0).Err()
}

// GetShortUrl fetches the url with the given code from the cache
func (c *RedisCache) GetShortUrl(code string) (string, error){
	return c.Client.Get(code).Result()
}
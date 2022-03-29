package cache

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisCache struct {
	Client *redis.Client
	Active bool
}

// NewRedisCache creates a new redis client and returns a new *RedisCache that contains the client
func NewRedisCache(addr, port, pass string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", addr, port),
		Password: pass,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return &RedisCache{Client: client, Active: false}, err
	}

	return &RedisCache{Client: client, Active: true}, nil
}

// SetShortUrl saves a short url code and url into the cache
func (c *RedisCache) SetShortUrl(code, url string) error {
	// if cache is not active
	if !c.Active {
		return nil
	}

	err := c.Client.Set(code, url, 0).Err()
	if err != nil {
		// disable cache
		c.Active = false
	}

	return err
}

// GetShortUrl fetches the url with the given code from the cache
func (c *RedisCache) GetShortUrl(code string) (string, error) {
	// if cache is not active
	if !c.Active {
		return "", nil
	}

	url, err := c.Client.Get(code).Result()
	if err != nil {
		// disable cache
		c.Active = false
	}

	return url, err
}

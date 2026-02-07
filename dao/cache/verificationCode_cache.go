package cache

import (
	"context"
	"fmt"
	"time"
)

type verificationCodeCache struct {
	cache *RedisClient
}

func NewVerificationCodeCache(cache *RedisClient) VerificationCodeCache {
	return &verificationCodeCache{
		cache: cache,
	}
}

func (c *verificationCodeCache) SaveVerificationCode(ctx context.Context, email, code string) error {
	key := fmt.Sprintf("verification:%s", email)
	return c.cache.Set(key, code, c.cache.RandExp(5*time.Minute))
}

func (c *verificationCodeCache) GetVerificationCode(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("verification:%s", email)
	var code string
	err := c.cache.Get(key, &code)
	return code, err
}

func (c *verificationCodeCache) DeleteVerificationCode(ctx context.Context, email string) error {
	key := fmt.Sprintf("verification:%s", email)
	return c.cache.Delete(key)
}

func (c *verificationCodeCache) SetRateLimit(ctx context.Context, email string) error {
	key := fmt.Sprintf("ratelimit:%s", email)
	return c.cache.Set(key, "1", 1*time.Minute)
}

func (c *verificationCodeCache) CheckRateLimit(ctx context.Context, email string) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", email)
	return c.cache.Exists(key), nil
}

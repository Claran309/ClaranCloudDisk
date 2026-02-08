package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	mu     sync.Mutex
}

func NewRedisClient(addr, password string, db int) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 100,
	})

	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

func (rc *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.Set(rc.ctx, key, data, expiration).Err()
}

func (rc *RedisClient) Get(key string, dest interface{}) error {
	data, err := rc.client.Get(rc.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (rc *RedisClient) Delete(key string) error {
	return rc.client.Del(rc.ctx, key).Err()
}

func (rc *RedisClient) SAdd(key string, members ...interface{}) error {
	return rc.client.SAdd(rc.ctx, key, members...).Err()
}

func (rc *RedisClient) SMembers(key string) ([]string, error) {
	return rc.client.SMembers(rc.ctx, key).Result()
}
func (rc *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	return rc.client.SIsMember(rc.ctx, key, member).Result()
}
func (rc *RedisClient) Expire(key string, expiration time.Duration) error {
	return rc.client.Expire(rc.ctx, key, expiration).Err()
}

// RandExp 防止缓存雪崩
func (rc *RedisClient) RandExp(base time.Duration) time.Duration {
	jitter := rand.Int63n(int64(base/5)) - int64(base/10)
	return time.Duration(int64(base) + jitter)
}

// Lock 获取分布式锁
func (rc *RedisClient) Lock(key string, expire time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	suc, err := rc.client.SetNX(rc.ctx, lockKey, "1", expire).Result()
	return suc, err
}

// Unlock 释放分布式锁
func (rc *RedisClient) Unlock(key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return rc.client.Del(rc.ctx, lockKey).Err()
}

// Clean 删除缓存
func (rc *RedisClient) Clean(keys ...string) error {
	return rc.client.Del(rc.ctx, keys...).Err()
}

func (rc *RedisClient) Exists(key string) bool {
	return rc.client.Exists(rc.ctx, key).Val() > 0
}

package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	SAdd(key string, members ...interface{}) error
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	Expire(key string, expiration time.Duration) error
	RandExp(base time.Duration) time.Duration
	Lock(key string, expire time.Duration) (bool, error)
	Unlock(key string) error
	Clean(keys ...string) error
	Exists(key string) bool
}

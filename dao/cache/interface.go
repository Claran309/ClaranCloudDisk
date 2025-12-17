package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	RandExp(base time.Duration) time.Duration
	Lock(key string, expire time.Duration) (bool, error)
	Unlock(key string) error
	Clean(keys ...string) error
	Exists(key string) bool
}

package util

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity float64
	tokens   float64
	rate     float64
	lastTime time.Time
	mu       sync.Mutex
}

func NewTokenBucket(capacity float64, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity: capacity,
		rate:     rate,
		tokens:   capacity,
		lastTime: time.Now(),
	}
}

// Allow 流出单个令牌检查
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN 流出多个令牌检查
func (tb *TokenBucket) AllowN(n float64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 更新桶内令牌
	now := time.Now()
	passedTime := now.Sub(tb.lastTime).Seconds()
	tb.tokens = tb.rate * passedTime
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastTime = now

	// 检查令牌
	if tb.tokens >= n {
		tb.tokens -= n
		return true
	}
	return false
}

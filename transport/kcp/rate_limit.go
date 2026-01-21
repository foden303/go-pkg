package kcp

import "sync"

type RateLimiter struct {
	rate     int64 // bytes per second
	tokens   int64
	lastTime int64
	mu       sync.Mutex
}

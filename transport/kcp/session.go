package kcp

import (
	"sync"
	"time"
)

type Session struct {
	mu            sync.RWMutex
	LastHeartbeat time.Time
	LastActivity  time.Time
	Data          interface{}
	ID            string
	Conn          *Conn
	IsAlive       bool
}

package kcp

import (
	"sync"
	"time"
)

// Pool manages a pool of KCP connections to multiple servers.
type Pool struct {
	cfg     *Config
	servers []*PoolConn
	mu      sync.RWMutex
}

// PoolConn represents a connection in the KCP connection pool.
type PoolConn struct {
	mu       sync.RWMutex
	Latency  time.Duration
	LastPing time.Time
	Addr     string
	Conn     *Conn
	FailCnt  int
	IsAlive  bool
}

// NewPool creates a new KCP connection pool with the given configuration.
func NewPool(cfg *Config) *Pool {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return &Pool{
		cfg:     cfg,
		servers: make([]*PoolConn, 0),
	}
}

// Add adds a new server connection to the pool.
func (p *Pool) Add(addr string) error {
	conn, err := Dial(addr, p.cfg)
	if err != nil {
		return err
	}

	pc := &PoolConn{
		Addr:     addr,
		Conn:     conn,
		IsAlive:  true,
		LastPing: time.Now(),
	}
	p.mu.Lock()
	p.servers = append(p.servers, pc)
	p.mu.Unlock()
	return nil
}

// GetBest returns the connection with the lowest latency that is currently alive.
func (p *Pool) GetBest() *Conn {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var best *PoolConn
	for _, pc := range p.servers {
		pc.mu.RLock()
		isActive := pc.IsAlive
		latency := pc.Latency
		pc.mu.RUnlock()
		if !isActive {
			continue
		}
		if best == nil || latency < best.Latency {
			best = pc
		}
	}
	if best == nil {
		return nil
	}
	return best.Conn
}

// HealthCheck performs a health check on all connections in the pool.
func (p *Pool) HealthCheck(timeout time.Duration) {
	p.mu.RLock()
	servers := make([]*PoolConn, len(p.servers))
	copy(servers, p.servers)
	p.mu.RUnlock()

	var wg sync.WaitGroup
	for _, pc := range servers {
		wg.Add(1)
		go func(pc *PoolConn) {
			defer wg.Done()
			latency, err := pc.Conn.Ping(timeout)
			pc.mu.Lock()
			if err != nil {
				pc.FailCnt++
				if pc.FailCnt >= 3 {
					pc.IsAlive = false
				}
			} else {
				pc.FailCnt = 0
				pc.Latency = latency
				pc.IsAlive = true
				pc.LastPing = time.Now()
			}
			pc.mu.Unlock()
		}(pc)
	}
	wg.Wait()
}

// Close closes all connections in the pool.
func (p *Pool) Close() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, pc := range p.servers {
		if pc.Conn != nil {
			pc.Conn.Close()
		}
	}
	p.servers = nil
}

// CheckService attempts to connect to the specified KCP service address within the given timeout.
func CheckService(addr string, timeout time.Duration) (time.Duration, error) {
	cfg := DefaultConfig()
	cfg.ConnTimeout = timeout
	cfg.MaxRetries = 1
	conn, err := Dial(addr, cfg)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return conn.Ping(timeout)
}

func CheckServices(addrs []string, timeout time.Duration) map[string]time.Duration {
	results := make(map[string]time.Duration)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			latency, err := CheckService(addr, timeout)
			if err == nil {
				mu.Lock()
				results[addr] = latency
				mu.Unlock()
			}
		}(addr)
	}
	wg.Wait()
	return results
}

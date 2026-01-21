package kcp

import "time"

const (
	_defaultSendWindowSize   = 1024
	_defaultRecvWindowSize   = 1024
	_defaultNoDelay          = 1
	_defaultInterval         = 10
	_defaultResend           = 2
	_defaultNoCongestion     = 1
	_defaultMTU              = 1400
	_defaultRateLimit        = 0 // unlimited
	_defaultMaxRetries       = 3
	_defaultRetryDelay       = 2 * time.Second
	_defaultConnTimeout      = 10 * time.Second
	_defaultReadTimeout      = 30 * time.Second
	_defaultWriteTimeout     = 30 * time.Second
	_defaultHeartbeatTimeout = 10 * time.Second
	_defaultACKNoDelay       = true
	_defaultWriteDelay       = false
)

type Config struct {
	ConnTimeout      time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	HeartbeatTimeout time.Duration
	RetryDelay       time.Duration

	SendWindowSize int
	RecvWindowSize int
	NoDelay        int
	Interval       int
	Resend         int
	NoCongestion   int
	MTU            int
	RateLimit      int
	MaxRetries     int

	ACKNoDelay bool
	WriteDelay bool
}

// DefaultConfig returns the default KCP configuration.
func DefaultConfig() *Config {
	return &Config{
		ConnTimeout:      _defaultConnTimeout,
		ReadTimeout:      _defaultReadTimeout,
		WriteTimeout:     _defaultWriteTimeout,
		HeartbeatTimeout: _defaultHeartbeatTimeout,
		RetryDelay:       _defaultRetryDelay,

		SendWindowSize: _defaultSendWindowSize,
		RecvWindowSize: _defaultRecvWindowSize,
		NoDelay:        _defaultNoDelay,
		Interval:       _defaultInterval,
		Resend:         _defaultResend,
		NoCongestion:   _defaultNoCongestion,
		MTU:            _defaultMTU,
		RateLimit:      _defaultRateLimit,
		MaxRetries:     _defaultMaxRetries,

		ACKNoDelay: _defaultACKNoDelay,
		WriteDelay: _defaultWriteDelay,
	}
}

// HighSpeedConfig returns a configuration optimized for high speed.
func HighSpeedConfig() *Config {
	cfg := DefaultConfig()
	cfg.SendWindowSize = 2048
	cfg.RecvWindowSize = 2048
	cfg.Interval = 5
	return cfg
}

// LowLatencyConfig returns a configuration optimized for low latency.
func LowLatencyConfig() *Config {
	cfg := DefaultConfig()
	cfg.Interval = 5
	return cfg
}

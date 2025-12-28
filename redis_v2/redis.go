package redis_v2

import (
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisType represents the type of Redis deployment.
type RedisType string

const (
	// ClusterType indicates a Redis Cluster deployment.
	ClusterType RedisType = "cluster"
	//  NodeType indicates a single Redis Node deployment.s
	NodeType RedisType = "node"
)

func (r RedisType) String() string {
	return string(r)
}

const (
	// Nil is an alias for redis.Nil
	Nil = redis.Nil
	// Default configuration values
	_defaultShowThreshold    = 100 * time.Microsecond
	_defaultPingTimeout      = 1 * time.Second
	_defaultReadWriteTimeout = 3 * time.Second
	_defaultQueryTimeout     = 5 * time.Second
)

var (
	// ErrNilNode indicates that the Redis node is nil.
	ErrNilNode = errors.New("nil redis node")
)

type (
	// Option defines a function type for configuring Redis.
	Option func(r *Redis)
	// IntPair represents a key with an integer score.
	Pair struct {
		Key   string
		Score int64
	}
	// FloatPair represents a key with a float score.
	FloatPair struct {
		Key   string
		Score float64
	}
	// Redis represents a Redis client configuration.
	Redis struct {
		Addr  string
		Type  RedisType
		User  string
		Pwd   string
		tls   bool
		hooks []redis.Hook
	}
	// RedisNode is an alias for redis.Cmdable
	RedisNode interface {
		redis.Cmdable
	}
	// GeoLocation is an alias for redis.GeoLocation
	GeoLocation = redis.GeoLocation
	// GeoRadiusQuery is an alias for redis.GeoRadiusQuery
	GeoRadiusQuery = redis.GeoRadiusQuery
	// GeoPos is an alias for redis.GeoPos
	GeoPos = redis.GeoPos
	// Pipeliner is an alias for redis.Pipeliner
	Pipeliner = redis.Pipeliner
	// Z is an alias for redis.Z
	Z = redis.Z
	// ZStore is an alias for redis.ZStore
	ZStore = redis.ZStore
	// IntCmd is an alias for redis.IntCmd
	IntCmd = redis.IntCmd
	// FloatCmd is an alias for redis.FloatCmd
	FloatCmd = redis.FloatCmd
	// StringCmd is an alias for redis.StringCmd
	StringCmd = redis.StringCmd
	// Script is an alias for redis.Script
	Script = redis.Script
	// Hook is an alias for redis.Hook
	Hook = redis.Hook
	// DialHook is an alias for redis.DialHook
	DialHook = redis.DialHook
	// ProcessHook is an alias for redis.ProcessHook
	ProcessHook = redis.ProcessHook
	// ProcessPipelineHook is an alias for redis.ProcessPipelineHook
	ProcessPipelineHook = redis.ProcessPipelineHook
	// Cmder is an alias for redis.Cmder
	Cmder = redis.Cmder
)

// New creates a new Redis instance with the given address and options.
func New(addr string, opts ...Option) *Redis {
	return newRedis(addr, opts...)
}

// newRedis is an internal function to create a new Redis instance.
func newRedis(addr string, opts ...Option) *Redis {
	r := &Redis{
		Addr: addr,
		Type: NodeType,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// NewScript creates a new Redis script with the given script content.
func NewScript(script string) *Script {
	return redis.NewScript(script)
}

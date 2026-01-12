package config

type Observer func(string, Value)

type Config interface {
	Load() error
	Scan(v any) error
	Value(key string) Value
	Watch(key string, o Observer) error
	Close() error
}

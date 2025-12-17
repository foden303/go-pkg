package logger

import (
	"go-pkg/logger/graylog"
	"go-pkg/logger/sentry"
)

type Option func(*Config)

func EnableGraylog(url string) Option {
	return func(cfg *Config) {
		cfg.enableGraylog = true
		cfg.graylogCore = graylog.New(graylog.Config{
			Address:    url,
			ServerName: cfg.ServerName,
			Version:    cfg.Version,
		}, newDefaultConfig())
	}
}

func EnableSentry(url string) Option {
	return func(cfg *Config) {
		cfg.enableSentry = true
		cfg.sentryCore = sentry.New(sentry.Config{
			DSN:         url,
			Version:     cfg.Version,
			ServerName:  cfg.ServerName,
			Environment: cfg.Environment,
		}, newDefaultConfig())
	}
}

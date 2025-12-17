package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type impl struct {
	logger *zap.Logger
}

// New creates a new Monitor instance
func New(cfg Config, opts ...Option) Logger {
	// Set default meta tags
	if cfg.ExtraTags == nil {
		cfg.ExtraTags = make(map[string]interface{})
	}
	cfg.ExtraTags["server.name"] = cfg.ServerName
	cfg.ExtraTags["environment"] = cfg.Environment
	cfg.ExtraTags["version"] = cfg.Version

	// Options
	for _, opt := range opts {
		opt(&cfg)
	}

	// Create default zap core configuration
	zapCfg := newDefaultConfig()

	var cores []zapcore.Core
	// Console output
	if cfg.Environment != EnvProd.String() {
		consoleCore := newConsoleCore(zapCfg)
		cores = append(cores, consoleCore)
	}
	// Graylog output
	if cfg.enableGraylog {
		cores = append(cores, cfg.graylogCore)
	}

	// Sentry output
	if cfg.enableSentry {
		cores = append(cores, cfg.sentryCore)
	}

	// Create underlying zap core
	core := zapcore.NewTee(cores...)
	logger := zap.New(core).With(zapFields(cfg.ExtraTags)...)

	return &impl{
		logger: logger,
	}
}

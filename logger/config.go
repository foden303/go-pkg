package logger

import (
	"go.uber.org/zap/zapcore"
)

// Config holds Monitor configuration
type Config struct {
	ServerName  string
	ServerType  string
	Environment string
	Version     string
	ExtraTags   map[string]interface{}

	enableSentry  bool
	sentryCore    zapcore.Core
	enableGraylog bool
	graylogCore   zapcore.Core
}

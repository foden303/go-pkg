package sentry

import (
	"log"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

type impl struct {
	levelEnabler zapcore.LevelEnabler
	encoder      zapcore.Encoder
}

func New(cfg Config, zapCfg zapcore.EncoderConfig) zapcore.Core {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Release:     cfg.Version,
		ServerName:  cfg.ServerName,
		Environment: cfg.Environment,
		SampleRate:  1,

		Debug:            true,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1,
	})
	if err != nil {
		log.Fatalf("sentry initialization failed: %v\n", err)
	}
	encoder := zapcore.NewJSONEncoder(zapCfg)
	return &impl{
		levelEnabler: zapcore.ErrorLevel,
		encoder:      encoder,
	}
}

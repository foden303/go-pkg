package graylog

import (
	"log"

	"github.com/Graylog2/go-gelf/gelf"
	"go.uber.org/zap/zapcore"
)

type impl struct {
	writer       *gelf.Writer
	levelEnabler zapcore.LevelEnabler
	encoder      zapcore.Encoder
	cfg          Config
}

func New(cfg Config, zapCfg zapcore.EncoderConfig) zapcore.Core {
	// Create a graylog connection (GELF UDP)
	writer, err := gelf.NewWriter(cfg.Address)
	if err != nil {
		log.Fatalf("graylog initialization failed: %v\n", err)
	}

	zapCfg.EncodeLevel = encodeLevel
	encoder := zapcore.NewJSONEncoder(zapCfg)
	return &impl{
		levelEnabler: zapcore.InfoLevel,
		encoder:      encoder,
		writer:       writer,
		cfg:          cfg,
	}
}

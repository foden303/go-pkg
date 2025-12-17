package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDefaultConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // to show {"level": "info"} or to show {"level": "INFO"}
		EncodeTime:     zapcore.ISO8601TimeEncoder,  // to format time as {"timestamp":"2021-06-21T09:25:51.230+08:00"}
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func newConsoleCore(cfg zapcore.EncoderConfig) zapcore.Core {
	encoder := zapcore.NewJSONEncoder(cfg)
	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zap.InfoLevel)
}

func zapFields(tags map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(tags))
	for k, v := range tags {
		fields = append(fields, zap.String(sanitizeKey(k), fmt.Sprintf("%v", v)))
	}
	return fields
}

package sentry

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

func mapZapToSentryLevel(lv zapcore.Level) sentry.Level {
	switch lv {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelError
	}
}

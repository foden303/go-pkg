package graylog

import "go.uber.org/zap/zapcore"

func mapZapToGraylogLevel(level zapcore.Level) int32 {
	switch level {
	case zapcore.DebugLevel:
		return 7 // Debug
	case zapcore.InfoLevel:
		return 6 // Info
	case zapcore.WarnLevel:
		return 4 // Warning
	case zapcore.ErrorLevel:
		return 3 // Error
	case zapcore.DPanicLevel:
		return 2 // Critical
	case zapcore.PanicLevel:
		return 1 // Alert
	case zapcore.FatalLevel:
		return 0 // Emergency
	default:
		return 6 // Default: Info
	}
}

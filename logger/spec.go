package logger

import "go.uber.org/zap"

type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(err error, format string, args ...interface{})
	Tracef(message string, tags map[string]interface{})
	Write(p []byte) (n int, err error)

	With(tags map[string]interface{}) Logger
	Sync() error
	Zap() *zap.Logger
}

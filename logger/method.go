package logger

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	pkgErr "go-pkg/err"

	"go.uber.org/zap"
)

// Infof logs the message using the info level
func (i *impl) Infof(format string, args ...interface{}) {
	i.logger.Info(fmt.Sprintf(format, args...))
}

// Debugf logs the message using the debug level
func (i *impl) Debugf(format string, args ...interface{}) {
	i.logger.Debug(fmt.Sprintf(format, args...))
}

// Warnf logs the message using the warn level
func (i *impl) Warnf(format string, args ...interface{}) {
	i.logger.Warn(fmt.Sprintf(format, args...))
}

// Errorf logs the message using the error level
func (i *impl) Errorf(err error, msg string, args ...interface{}) {
	// Check kind of error
	var stdErr pkgErr.Error
	if errors.As(err, &stdErr) {
		if stdErr.HttpCode() == http.StatusInternalServerError && stdErr.Err() != nil {
			err = stdErr.Err()
		}
	}

	// Capture error.key, error.message and error.stack to log
	logger := i.logger.With(
		zap.String("error_kind", reflect.TypeOf(err).String()),
		zap.String("error_message", err.Error()),
	)
	if v, ok := err.(stackTracer); ok {
		stack := fmt.Sprintf("%+v", v.StackTrace())
		if len(stack) > 0 && stack[0] == '\n' {
			stack = stack[1:]
		}
		logger = logger.With(zap.String("error_stack", stack))
	}

	if msg != "" {
		logger.Error(fmt.Sprintf(msg+". Err: %v", append(args, err)...))
	} else {
		logger.Error(fmt.Sprintf("Err: %v", err))
	}
}

// Tracef logs the message using the trace level.
func (i *impl) Tracef(message string, tags map[string]interface{}) {
	i.logger.Info(message, zapFields(tags)...)
}

func (i *impl) Write(p []byte) (n int, err error) {
	i.logger.Info(strings.TrimSpace(string(p)))
	return len(p), nil
}

func (i *impl) Sync() error {
	return i.logger.Sync()
}

func (i *impl) Zap() *zap.Logger {
	return i.logger
}

package logger

import pkgerr "github.com/pkg/errors"

// This is to extract the stack trace from pkgerrors
type stackTracer interface {
	StackTrace() pkgerr.StackTrace
}

package sentry

import (
	"encoding/json"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

func (i *impl) Enabled(level zapcore.Level) bool {
	return i.levelEnabler.Enabled(level)
}

func (i *impl) With(fields []zapcore.Field) zapcore.Core {
	clone := *i
	clone.encoder = i.encoder.Clone()
	for _, f := range fields {
		f.AddTo(clone.encoder)
	}
	return &clone
}

func (i *impl) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if i.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, i)
	}
	return checkedEntry
}

func (i *impl) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Clone encoder to encode all fields into a map
	enc := i.encoder.Clone()
	for _, f := range fields {
		f.AddTo(enc)
	}

	buf, err := enc.EncodeEntry(entry, nil)
	if err != nil {
		return err
	}

	// Build extras map
	var extras map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &extras); err != nil {
		return err
	}

	sentryEvent := sentry.NewEvent()
	sentryEvent.Level = mapZapToSentryLevel(entry.Level)
	sentryEvent.Message = entry.Message
	sentryEvent.Timestamp = entry.Time
	sentryEvent.Extra = extras

	sentry.CaptureEvent(sentryEvent)
	sentry.Flush(2 * time.Second)

	return nil
}

func (i *impl) Sync() error {
	sentry.Flush(2 * time.Second)
	return nil
}

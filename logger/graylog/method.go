package graylog

import (
	"encoding/json"
	"time"

	"github.com/Graylog2/go-gelf/gelf"
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
	extras["release_version"] = i.cfg.Version
	// ❌ Delete timestamp if exists and wrong type
	delete(extras, "timestamp")

	msg := &gelf.Message{
		Version:  i.cfg.Version,
		Host:     i.cfg.ServerName,
		Short:    entry.Message,
		TimeUnix: float64(time.Now().UnixNano()) / 1e9,
		Level:    mapZapToGraylogLevel(entry.Level),
		Extra:    extras,
	}
	return i.writer.WriteMessage(msg)
}

func (i *impl) Sync() error {
	return nil
}

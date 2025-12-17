package graylog

import "go.uber.org/zap/zapcore"

func encodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(int(mapZapToGraylogLevel(lv)))
}

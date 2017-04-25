package trclient

import (
	"github.com/drkaka/lg"
	"go.uber.org/zap/zapcore"
)

func p(msg string, fields ...zapcore.Field) {
	if lg.L(nil) != nil {
		lg.L(nil).Debug(msg, fields...)
	}
}

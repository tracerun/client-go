package trclient

import (
	"errors"

	"github.com/drkaka/lg"
	"go.uber.org/zap/zapcore"
)

var (
	// ErrRoute route information wrong
	ErrRoute = errors.New("route wrong")
)

func p(msg string, fields ...zapcore.Field) {
	if lg.L(nil) != nil {
		lg.L(nil).Debug(msg, fields...)
	}
}

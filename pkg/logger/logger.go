package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	Logger *zap.Logger
}

func NewLogger() (*ZapLogger, error) {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()
	return &ZapLogger{Logger: logger}, nil
}
func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}
func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}
func Now() time.Time {
	return time.Now()
}
func Since(t time.Time) time.Duration {
	return time.Since(t)
}
func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(5); v != nil {
		if tid, ok := v.(string); ok {
			return tid
		}
	}
	return ""
}
func (l *ZapLogger) WithTrace(ctx context.Context) *zap.Logger {
	traceID := GetTraceID(ctx)
	return l.Logger.With(zap.String("trace_id", traceID))
}

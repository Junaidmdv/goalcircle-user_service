package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ZapLogger struct {
	Logger *zap.Logger `json:"logger,omitempty"`
}

func NewLogger() (*ZapLogger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{Logger: z}, nil
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

package logctx

import (
	"context"

	"go.uber.org/zap"
)

// From return the zap.Logger from a context. If there is no configuration a default instance is returned.
func From(ctx context.Context) *zap.Logger {
	loggerGeneric := ctx.Value(defaultOptions.ctxKey)
	if loggerGeneric == nil {
		return defaultOptions.logger
	}
	logger, ok := loggerGeneric.(*zap.Logger)
	if !ok {
		return defaultOptions.logger
	}
	return logger
}

// WithLogger add the given zap.Logger and returning the new context.Context.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, defaultOptions.ctxKey, logger)
}

// WithFields returns a new context with the logger that includes the given fields.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return WithLogger(ctx, From(ctx).With(fields...))
}

// Error is a helper function to log an error with the given context.
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Error(msg, fields...)
}

// Info is a helper function to log an info with the given context.
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Info(msg, fields...)
}

// Debug is a helper function to log a debug with the given context.
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Debug(msg, fields...)
}

// Warn is a helper function to log a warning with the given context.
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Warn(msg, fields...)
}

// Fatal is a helper function to log a fatal with the given context.
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Fatal(msg, fields...)
}

// Panic is a helper function to log a panic with the given context.
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	From(ctx).Panic(msg, fields...)
}

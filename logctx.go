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

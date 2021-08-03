package logctx

import (
	"fmt"

	"go.uber.org/zap"
)

var (
	defaultOptions options
)

type options struct {
	logger *zap.Logger
	ctxKey ctxKey
}

func defaultOpts() options {
	return options{
		ctxKey: "jamillosantos/logctx/logging#logger",
	}
}

type ctxKey string

func init() {
	err := Initialize()
	if err != nil {
		panic(err)
	}
}

type Option = func(*options)

// WithDefaultLogger is the option that initializes the default logger instance.
func WithDefaultLogger(logger *zap.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// WithContextKey is the option that initializes the default context key.
func WithContextKey(key string) Option {
	return func(o *options) {
		o.ctxKey = ctxKey(key)
	}
}

// Initialize intiializes the default configuration for the logger.
func Initialize(opts ...Option) error {
	options := defaultOpts()
	for _, o := range opts {
		o(&options)
	}
	if options.logger == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return fmt.Errorf("failed initializing zap production: %w", err)
		}
		options.logger = logger
	}
	defaultOptions = options
	return nil
}

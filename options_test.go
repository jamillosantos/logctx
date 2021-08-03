package logctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithDefaultLogger(t *testing.T) {
	var opts options
	logger := zap.NewNop()
	WithDefaultLogger(logger)(&opts)
	assert.Equal(t, logger, opts.logger)
}

func TestWithContextKey(t *testing.T) {
	var opts options
	wantContextKey := "wantContextKey"
	WithContextKey(wantContextKey)(&opts)
	assert.Equal(t, ctxKey(wantContextKey), opts.ctxKey)
}

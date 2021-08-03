package logctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInitialize(t *testing.T) {
	logger := zap.NewNop()
	assert.NotEqual(t, logger, defaultOptions.logger)
	err := Initialize(WithDefaultLogger(logger))
	require.NoError(t, err)
	assert.Equal(t, logger, defaultOptions.logger)
}

func TestFrom(t *testing.T) {
	defaultOptions.logger = zap.NewNop()

	t.Run("return the default logger when there is no context", func(t *testing.T) {
		assert.Equal(t, defaultOptions.logger, From(context.TODO()))
	})

	t.Run("return the default logger when the context has the wrong logger type", func(t *testing.T) {
		ctx := context.WithValue(context.TODO(), defaultOptions.ctxKey, "this is not a logger")
		assert.Equal(t, defaultOptions.logger, From(ctx))
	})

	t.Run("return the logger within a context", func(t *testing.T) {
		l1 := zap.NewNop()
		ctx := WithLogger(context.TODO(), l1)
		assert.Equal(t, l1, From(ctx))
	})
}

func TestWithLogger(t *testing.T) {
	wantLogger := zap.NewNop()
	ctx := WithLogger(context.TODO(), wantLogger)
	gotLogger := From(ctx)
	require.Equal(t, wantLogger, gotLogger)
}

func TestWithFields(t *testing.T) {
	core, obs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	ctx := WithLogger(context.TODO(), logger)
	logger.Info("test")
	ctx = WithFields(ctx, zap.String("foo", "bar"))
	From(ctx).Info("test with fields")
	entries := obs.All()
	assert.Len(t, entries, 2)
	assert.Empty(t, entries[0].ContextMap())
	assert.Len(t, entries[1].ContextMap(), 1)
	assert.Contains(t, entries[1].ContextMap(), "foo")
	assert.NotContains(t, entries[1].ContextMap(), "bar")
}

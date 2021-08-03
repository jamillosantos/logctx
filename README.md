# logctx

Utility library for logging.

## Usage

Adding a logger into a context.Context:

```go
// ...
ctx = logctx.WithLogger(ctx, logger.With(zap.String("request_id", ulid.New().String())))
// Now the logger inside of `ctx` has `request_id` as field.
// ...
```

Extracting a logger from a context.Context:

```go
// ...
logger := logctx.From(ctx)
// Now you can use logger as you wish.
// ...
```

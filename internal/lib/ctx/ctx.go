package jctx

import (
	"context"
)

type contextKey string

func WithContext[T any](ctx context.Context, key contextKey, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

func FromContext[T any](ctx context.Context, key contextKey) (T, bool) {
	value, ok := ctx.Value(key).(T)
	return value, ok
}

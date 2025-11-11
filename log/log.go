package log

import (
	"context"
	"fmt"
)

type CtxKeyIgnore struct{}

var ctxKeyIgnore CtxKeyIgnore

// 根据上下文，不打日志
func WithRestLogIgnore(ctx context.Context) context.Context {
	v := ctx.Value(ctxKeyIgnore)
	if v != nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKeyIgnore, true)
}

func RestLogIgnoreFrmCtx(ctx context.Context) bool {
	v := ctx.Value(ctxKeyIgnore)
	if v == nil {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}

var _ = ILogger(&DefaultLogger{})

type ILogger interface {
	DebugCtx(ctx context.Context, msg string)
	InfoCtx(ctx context.Context, msg string)
	WarnCtx(ctx context.Context, msg string)
	ErrorCtx(ctx context.Context, msg string)
}

type DefaultLogger struct {
}

func (l *DefaultLogger) DebugCtx(ctx context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) InfoCtx(ctx context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) WarnCtx(ctx context.Context, msg string) {
	fmt.Println(msg)
}

func (l *DefaultLogger) ErrorCtx(ctx context.Context, msg string) {
	fmt.Println(msg)
}

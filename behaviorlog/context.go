package behaviorlog

import "context"

const LogContextName = "behaviorlog"

func (c *LogContext) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, LogContextName, c)
}
func FromCtx(ctx context.Context) *LogContext {
	if ctx == nil {
		return NewNopContext()
	}
	v := ctx.Value(LogContextName)
	if v == nil {
		return NewNopContext()
	}
	logContext, ok := v.(*LogContext)
	if !ok {
		return NewNopContext()
	}
	return logContext
}

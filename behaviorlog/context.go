package behaviorlog

import "context"

const LogContextName = "behaviorlog"

func (c *LogContext) ToCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, LogContextName, c)
}
func FromCtx(ctx context.Context) *LogContext {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(LogContextName)
	if v == nil {
		return nil
	}
	logContext, ok := v.(*LogContext)
	if !ok {
		return nil
	}
	return logContext
}

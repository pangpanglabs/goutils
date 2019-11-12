package ctxbase

import (
	"context"

	uuid "github.com/gofrs/uuid"
)

const ContextBaseName = "ContextBase"

type ContextBase struct {
	RequestID string
	ActionID  string
}

func FromCtx(ctx context.Context) *ContextBase {
	if ctx == nil {
		return nil
	}
	c, ok := ctx.Value(ContextBaseName).(*ContextBase)
	if !ok {
		return nil
	}
	return c
}

func NewID() string {
	u, err := uuid.NewV1()
	if err != nil {
		return ""
	}
	return u.String()
}

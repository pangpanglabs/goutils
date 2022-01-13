package ctxaws

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/session"
)

const ContextAWSName = "ContextAWS"

type ContextAWS struct {
	*session.Session
}

func New(sess *session.Session) *ContextAWS {

	return &ContextAWS{Session: sess}
}

func (aws *ContextAWS) NewSession(ctx context.Context) *session.Session {
	sess := aws.Session
	func(sess interface{}, ctx context.Context) {
		if s, ok := sess.(interface{ SetContext(context.Context) }); ok {
			s.SetContext(ctx)
		}
	}(sess, ctx)

	return sess
}

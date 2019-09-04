package cronjob

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/pangpanglabs/goutils/behaviorlog"
)

func Recover() Middleware {
	stackSize := 4 << 10 // 4 KB
	disableStackAll := false
	disablePrintStack := false

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context) error {

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, stackSize)
					length := runtime.Stack(stack, !disableStackAll)
					if !disablePrintStack {
						log.Printf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					behaviorlog.FromCtx(ctx).WithError(err)
				}
			}()

			return next(ctx)
		}
	}
}

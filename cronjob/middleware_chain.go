package cronjob

import "context"

type HandlerFunc func(ctx context.Context) error
type Middleware func(next HandlerFunc) HandlerFunc
type middlewareChain struct {
	middlewares []Middleware
}

func (c middlewareChain) run(f HandlerFunc) HandlerFunc {
	for i := range c.middlewares {
		f = c.middlewares[len(c.middlewares)-1-i](f)
	}

	return f
}

func (c *middlewareChain) append(middlewares ...Middleware) {
	newMiddlewares := make([]Middleware, 0, len(c.middlewares)+len(middlewares))
	newMiddlewares = append(newMiddlewares, c.middlewares...)
	newMiddlewares = append(newMiddlewares, middlewares...)

	c.middlewares = newMiddlewares
}

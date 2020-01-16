package cronjob

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type Cron struct {
	name  string
	cron  *cron.Cron
	chain middlewareChain
}

type Job struct {
	Name string
	c    *Cron
}

func New(name string) *Cron {
	c := Cron{
		name: name,
		cron: cron.New(),
	}

	c.chain.append(
		Recover(),
	)

	return &c
}

func (c *Cron) Use(middlewares ...Middleware) {
	c.chain.append(middlewares...)
}

func (c *Cron) AddJob(job string) *Job {
	return &Job{
		Name: job,
	}
}
func (j *Job) AddAction(action, spec string, f HandlerFunc) *Job {
	j.c.cron.AddFunc(spec, func() {
		ctx := context.Background()
		if err := j.c.chain.run(j.Name, action, f)(ctx); err != nil {
			logrus.WithError(err).Error("")
		}
	})
	return j
}
func (c *Cron) AddAction(action, spec string, f HandlerFunc) {
	c.cron.AddFunc(spec, func() {
		ctx := context.Background()
		if err := c.chain.run(c.name, action, f)(ctx); err != nil {
			logrus.WithError(err).Error("")
		}
	})
}

func (c *Cron) Start(address string) {
	go c.cron.Start()

	e := echo.New()
	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "pong")
	})
	e.GET("/entries", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, c.cron.Entries())
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	if err := e.Start(address); err != nil {
		log.Println(err)
	}
	c.cron.Stop()
}

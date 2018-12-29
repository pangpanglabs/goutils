package cronjob

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

type Cron struct {
	cron  *cron.Cron
	chain middlewareChain
}

func New() *Cron {
	return &Cron{
		cron: cron.New(),
	}
}
func Default(serviceName string, behaviorlogKafkaConfig echomiddleware.KafkaConfig) *Cron {
	c := Cron{
		cron: cron.New(),
	}

	c.chain.append(
		BehaviorLogger(serviceName, behaviorlogKafkaConfig),
		Recover(),
	)

	return &c
}
func (c *Cron) Use(middlewares ...Middleware) {
	c.chain.append(middlewares...)
}
func (c *Cron) AddFunc(spec string, f HandlerFunc) {
	c.cron.AddFunc(spec, func() {
		ctx := context.Background()

		if err := c.chain.run(f)(ctx); err != nil {
			logrus.WithError(err).Error("")
		}
	})
}

func (c *Cron) Start() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	c.cron.Start()

	return fmt.Errorf("Got signal: %v", <-quit)
}

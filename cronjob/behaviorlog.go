package cronjob

import (
	"context"
	"time"

	"github.com/labstack/gommon/random"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

func BehaviorLogger(serviceName string, config kafka.Config) Middleware {
	var producer *kafka.Producer
	if config.Brokers != nil && config.Topic != "" {
		if p, err := kafka.NewProducer(config.Brokers, config.Topic,
			kafka.WithDefault(),
			kafka.WithTLS(config.SSL)); err != nil {
			logrus.Error("Create Kafka Producer Error", err)
		} else {
			producer = p
		}
	}

	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context) (err error) {
			behaviorLogContext := behaviorlog.NewNopContext()
			behaviorLogContext.Producer = producer
			behaviorLogContext.Service = serviceName
			behaviorLogContext.RequestID = random.String(32)
			behaviorLogContext.ActionID = random.String(32)
			behaviorLogContext.Timestamp = time.Now()
			behaviorLogContext.RemoteIP = "127.0.0.1"
			behaviorLogContext.Host = "127.0.0.1"

			behaviorLogContext.Uri = ""
			behaviorLogContext.Method = ""
			behaviorLogContext.Path = ""
			behaviorLogContext.Params = map[string]interface{}{}
			behaviorLogContext.Referer = ""
			behaviorLogContext.UserAgent = ""
			behaviorLogContext.RequestLength = 0
			behaviorLogContext.BizAttr = map[string]interface{}{}
			behaviorLogContext.AuthToken = ""

			err = next(context.WithValue(ctx, behaviorlog.LogContextName, behaviorLogContext))
			if err != nil {
				behaviorLogContext.Err = err.Error()
			}

			behaviorLogContext.Write()

			return err
		}
	}
}

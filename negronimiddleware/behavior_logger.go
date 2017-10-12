package negronimiddleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

const BehaviorLoggerName = "BehaviorLogger"

type BehaviorLogger struct {
	serviceName string
	producer    *kafka.Producer
}

func NewBehaviorLogger(serviceName string, brokers []string, topic string) *BehaviorLogger {
	b := BehaviorLogger{serviceName: serviceName}
	if p, err := kafka.NewProducer(brokers, topic, func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	}); err != nil {
		logrus.Error("Create Kafka Producer Error", err)
	} else {
		b.producer = p
	}
	return &b
}
func (b *BehaviorLogger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	behaviorLogger := behaviorlog.New(b.serviceName, req, behaviorlog.KafkaProducer(b.producer))
	next(rw, req.WithContext(context.WithValue(req.Context(),
		BehaviorLoggerName, behaviorLogger,
	)))

	// behaviorLogger.Status = req.Response.StatusCode
	behaviorLogger.Write()
}

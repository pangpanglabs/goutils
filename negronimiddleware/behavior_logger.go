package negronimiddleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/sirupsen/logrus"
)

type BehaviorLogger struct {
	hostname    string
	serviceName string
	producer    *kafka.Producer
}

func NewBehaviorLogger(serviceName string, brokers []string, topic string, options ...func(*sarama.Config)) *BehaviorLogger {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.WithError(err).Error("Fail to get hostname")
	}

	b := BehaviorLogger{serviceName: serviceName, hostname: hostname}
	options = append(options, option)
	if p, err := kafka.NewProducer(brokers, topic, options...); err != nil {
		logrus.Error("Create Kafka Producer Error", err)
	} else {
		b.producer = p
	}

	return &b
}
func (b *BehaviorLogger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	behaviorLogger := behaviorlog.New(b.serviceName, req, behaviorlog.KafkaProducer(b.producer))
	behaviorLogger.Hostname = b.hostname

	next(rw, req.WithContext(context.WithValue(req.Context(),
		behaviorlog.LogContextName, behaviorLogger,
	)))

	// behaviorLogger.Status = req.Response.StatusCode
	behaviorLogger.Write()
}
func option(c *sarama.Config) {
	c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
	c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	c.Producer.Return.Successes = false
}

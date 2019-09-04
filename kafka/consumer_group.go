package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type ConsumerGroupHandler struct {
	messages chan *sarama.ConsumerMessage
}

func NewConsumerGroup(groupId string, brokers []string, topic string, options ...func(*sarama.Config)) (*ConsumerGroupHandler, error) {
	kafkaConfig := sarama.NewConfig()
	for _, option := range options {
		option(kafkaConfig)
	}

	if kafkaConfig.Version == sarama.MinVersion {
		kafkaConfig.Version = sarama.V0_10_2_0
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupId, kafkaConfig)
	if err != nil {
		return nil, err
	}

	handler := ConsumerGroupHandler{
		messages: make(chan *sarama.ConsumerMessage),
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		log.Println("Initiating shutdown of consumer group...")
		consumerGroup.Close()
		os.Exit(1)
	}()

	go func() {
		for {
			if err := consumerGroup.Consume(context.Background(), []string{topic}, &handler); err != nil {
				logrus.WithError(err).Error("Fail to consume kafka message")
				return
			}
		}
	}()

	logrus.WithField("topic", topic).Info("Start to consume")

	return &handler, nil
}
func (c *ConsumerGroupHandler) Messages() (<-chan *sarama.ConsumerMessage, error) {
	return c.messages, nil
}
func (c *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for m := range claim.Messages() {
		c.messages <- m
		//MarkOffset 并不是实时写入kafka，有可能在程序crash时丢掉未提交的offset
		sess.MarkMessage(m, "")
	}
	return nil
}

func (ConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	logrus.WithFields(logrus.Fields{
		"MemberID":     s.MemberID(),
		"GenerationID": s.GenerationID(),
		"Claims":       s.Claims(),
	}).Info("Setup consumer group")
	return nil
}
func (ConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	logrus.WithFields(logrus.Fields{
		"MemberID":     s.MemberID(),
		"GenerationID": s.GenerationID(),
		"Claims":       s.Claims(),
	}).Info("Cleanup consumer group")
	return nil
}

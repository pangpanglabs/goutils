package kafka

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
)

type Producer struct {
	topic    string
	producer sarama.AsyncProducer
}

func NewProducer(brokers []string, topic string, options ...func(*sarama.Config)) (*Producer, error) {
	kafkaConfig := sarama.NewConfig()
	for _, option := range options {
		option(kafkaConfig)
	}

	producer, err := sarama.NewAsyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			log.Printf("Failed to send log entry to kafka: %v\n", err)
		}
	}()

	return &Producer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (p *Producer) Send(v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}

package kafka

import (
	"log"
	"os"
	"os/signal"

	"sync"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	topic      string
	partitions []int32
	offset     int64
	consumer   sarama.Consumer
	closing    chan struct{}
}

func NewConsumer(brokers []string, topic string, partitions []int32, initialOffset int64, options ...func(*sarama.Config)) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	for _, option := range options {
		option(kafkaConfig)
	}
	kafkaConsumer, err := sarama.NewConsumer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		topic:      topic,
		partitions: partitions,
		offset:     initialOffset,
		consumer:   kafkaConsumer,
		closing:    make(chan struct{}),
	}

	return consumer, nil
}

func (c *Consumer) Close() {
	close(c.closing)
}

func (c *Consumer) Messages() (<-chan *sarama.ConsumerMessage, error) {
	var partitionList []int32
	if len(c.partitions) != 0 {
		partitionList = c.partitions
	} else if list, err := c.consumer.Partitions(c.topic); err != nil {
		return nil, err
	} else {
		partitionList = list
	}

	messages := make(chan *sarama.ConsumerMessage, 1024)
	wg := sync.WaitGroup{}
	closing := make(chan struct{})

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Kill, os.Interrupt)
		<-signals
		log.Println("Initiating shutdown of consumer...")
		close(closing)
	}()

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, c.offset)
		if err != nil {
			return nil, err
		}

		go func(pc sarama.PartitionConsumer) {
			<-closing
			pc.AsyncClose()
		}(pc)

		wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for message := range pc.Messages() {
				messages <- message
			}
		}(pc)
	}

	go func() {
		wg.Wait()
		log.Println("Done consuming topic:", c.topic)
		close(messages)
	}()

	return messages, nil
}

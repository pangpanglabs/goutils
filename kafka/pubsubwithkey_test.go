package kafka_test

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hillfolk/goutils/kafka"
	"github.com/hillfolk/goutils/test"
)

var (
	testBrokers      = []string{"localhost:9092"}
	testTopic        = "multi-partitions-test"
	testkeys         = []string{"1009348249", "2009348249", "3009348249"}
	testMessageCount = 1000
)

func TestSend(t *testing.T) {
	t.Run("send with partitions key", func(t *testing.T) {
		// given
		produceMessages := map[int]interface{}{}
		consumeMessages := map[int]interface{}{}

		producer, _ := kafka.NewProducer(testBrokers, testTopic, func(c *sarama.Config) {
			c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
			c.Producer.Compression = sarama.CompressionSnappy   // Compress messages
			c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
		})
		defer producer.Close()

		consumer, _ := kafka.NewConsumer(testBrokers, testTopic, nil, sarama.OffsetNewest)

		messages, _ := consumer.Messages()
		done := make(chan struct{})
		i := 0
		go func() {
			for m := range messages {
				t.Logf("[Receive] Offset:%d\tPartition:%d\tValue:%v\n", m.Offset, m.Partition, string(m.Value))

				var v map[string]interface{}
				d := json.NewDecoder(bytes.NewReader(m.Value))
				d.UseNumber()
				d.Decode(&v)

				idx, err := v["idx"].(json.Number).Int64()
				test.Ok(t, err)
				msg, err := v["msg"].(json.Number).Int64()
				test.Ok(t, err)
				consumeMessages[int(idx)] = int(msg)
				i++
				if i == testMessageCount*len(testkeys) {
					break
				}
			}

			consumer.Close()
			close(done)
		}()

		// when
		for k := 0; k < len(testkeys); k++ {
			testKey := testkeys[k]
			for i := 0; i < testMessageCount; i++ {
				idx := i + (k * testMessageCount)
				produceMessage := map[string]interface{}{
					"idx": idx,
					"key": testKey,
					"msg": rand.Int(),
				}
				err := producer.SendWithKey(produceMessage, testKey)
				test.Ok(t, err)
				t.Log("[Send]", produceMessage)
				produceMessages[idx] = produceMessage["msg"]
			}
		}

		// then
		<-done

		test.Equals(t, len(produceMessages), len(consumeMessages))
		for i := 0; i < len(produceMessages); i++ {
			test.Equals(t, produceMessages[i], consumeMessages[i])
		}
	})
}

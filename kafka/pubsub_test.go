package kafka_test

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/pangpanglabs/goutils/kafka"
	"github.com/pangpanglabs/goutils/test"
)

var (
	brokers      = []string{"steamer-01.srvs.cloudkafka.com:9093", "steamer-02.srvs.cloudkafka.com:9093", "steamer-03.srvs.cloudkafka.com:9093"}
	topic        = "7uf2-default"
	messageCount = 10
)

func TestPubSub(t *testing.T) {
	produceMessages := map[int]interface{}{}
	consumeMessages := map[int]interface{}{}

	producer, err := kafka.NewProducer(brokers, topic, func(c *sarama.Config) {
		if tlsConfig, err := createTlsConfiguration(); err == nil {
			c.Net.TLS.Config = tlsConfig
			c.Net.TLS.Enable = true
		}
	})
	test.Ok(t, err)

	consumer, err := kafka.NewConsumer(brokers, topic, nil, sarama.OffsetNewest, func(c *sarama.Config) {
		if tlsConfig, err := createTlsConfiguration(); err == nil {
			c.Net.TLS.Config = tlsConfig
			c.Net.TLS.Enable = true
		}
	})
	test.Ok(t, err)

	messages, err := consumer.Messages()
	test.Ok(t, err)

	closing := make(chan struct{})
	i := 0
	go func() {
		for m := range messages {
			var v map[string]interface{}
			d := json.NewDecoder(bytes.NewReader(m.Value))
			d.UseNumber()
			err := d.Decode(&v)
			test.Ok(t, err)

			t.Logf("[Receive] Offset:%d\tPartition:%d\tValue:%v\n", m.Offset, m.Partition, v)
			idx, err := v["idx"].(json.Number).Int64()
			test.Ok(t, err)
			msg, err := v["msg"].(json.Number).Int64()
			test.Ok(t, err)
			consumeMessages[int(idx)] = int(msg)
			i++
			if i == messageCount {
				break
			}
		}
		t.Log("closing")
		consumer.Close()
		closing <- struct{}{}
	}()

	// run producer
	for i := 0; i < messageCount; i++ {
		msg := map[string]interface{}{
			"idx": i,
			"msg": rand.Int(),
		}
		err := producer.Send(msg)
		test.Ok(t, err)
		t.Log("[Send]", msg)
		produceMessages[i] = msg["msg"]
	}

	<-closing
	t.Log("closed")

	for k, v := range produceMessages {
		m, ok := consumeMessages[k]
		if !ok {
			t.Fail()
		}
		test.Equals(t, v, m)
	}
}

func createTlsConfiguration() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair("cert.crt", "cert.key")
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile("cert.ca")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}, nil
}

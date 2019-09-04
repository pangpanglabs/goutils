package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Brokers []string
	Topic   string
	SSL     SslConfig
}

type SslConfig struct {
	Enable                                    bool
	ClientCertFile, ClientKeyFile, CACertFile string
}

func WithDefault() func(*sarama.Config) {
	return func(c *sarama.Config) {
		c.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
		c.Producer.Compression = sarama.CompressionGZIP     // Compress messages
		c.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	}

}

func WithTLS(config SslConfig) func(*sarama.Config) {
	return func(c *sarama.Config) {
		if config.Enable {
			tlsConfig, err := newTLSConfig(config.ClientCertFile, config.ClientKeyFile, config.CACertFile)
			if err != nil {
				logrus.Error("Unable new TLS config for kafka.", err)
			}
			c.Net.TLS.Enable = true
			c.Net.TLS.Config = tlsConfig
		}
	}
}

func newTLSConfig(clientCertFile, clientKeyFile, caCertFile string) (*tls.Config, error) {
	tlsConfig := tls.Config{}
	// Load client cert
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		return &tlsConfig, err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return &tlsConfig, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, err
}

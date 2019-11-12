package echomiddleware

import "github.com/hillfolk/goutils/kafka"

type KafkaConfig = kafka.Config

type ZipkinConfig struct {
	Collector struct {
		Url string
	}
	Addr, Service string
	Kafka         KafkaConfig
}

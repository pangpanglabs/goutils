package echomiddleware

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type ZipkinConfig struct {
	Collector struct {
		Url string
	}
	Addr, Service string
	Kafka         KafkaConfig
}

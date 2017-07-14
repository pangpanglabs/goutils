package echomiddleware

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type ZipkinConfig struct {
	Addr, Service string
	Kafka         KafkaConfig
}

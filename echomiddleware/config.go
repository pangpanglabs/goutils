package echomiddleware

type KafkaConfig struct {
	Brokers []string
	Topic   string
	SSL     struct {
		Enable                                    bool
		ClientCertFile, ClientKeyFile, CACertFile string
	}
}

type ZipkinConfig struct {
	Collector struct {
		Url string
	}
	Addr, Service string
	Kafka         KafkaConfig
}

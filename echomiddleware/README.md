# goutils/echomiddleware

## Getting Started

### Basic usage
A log is logged when the api is called
```golang
e := echo.New()
e.Pre(echomiddleware.ContextBase())
e.Use(echomiddleware.BehaviorLogger("xxx-Service", echomiddleware.KafkaConfig(
		echomiddleware.KafkaConfig{
			Brokers: []string{
				"127.0.0.1:9092",
			},
			Topic: "behaviorlog",
		},
	))
```
### Advanced usage
```golang
e := echo.New()
e.Pre(echomiddleware.ContextBase())
e.Use(echomiddleware.BehaviorLogger("xxx-Service", echomiddleware.KafkaConfig(
		echomiddleware.KafkaConfig{
			Brokers: []string{
				"127.0.0.1:9092",
			},
			Topic: "behaviorlog",
		},
	), func(logContext *behaviorlog.LogContext) {
		logContext.BodyHide = true//Optional: Available when performing scheduled tasks to save large amounts of data
	}))
```
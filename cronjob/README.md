# goutils/cronjob

```go
import github.com/pangpanglabs/goutils/cronjob
```

## Getting Started

```go
// Create default cronjob
// Includes some default middlewares
c := cronjob.Default(config.ServiceName, config.BehaviorLog.Kafka)

// Add ContextDB middleware
c.Use(cronjob.ContextDB(config.ServiceName, db, config.Database.Logger.Kafka))

// Add job
c.AddFunc("0 0 1 * * *", func(ctx context.Context) error {
	// Start everyday 1 o'clock
	return nil
})

// Start cron job
log.Println(c.Start())
```

### `cronjob.Default`

Create cron job with default middleware

- `cronjob.BehaviorLogger` - Behavior Log Middleware.
- `cronjob.Recover` - Panic Recovery Middleware.

### `cronjob.New`

If you want to create without middleware, use like this.
```go
cronjob.New()
```

## Custom Middleware

If you want to add your own middleware, use like this.

```go
c.Use(func(next cronjob.HandlerFunc) cronjob.HandlerFunc {
	return func(ctx context.Context) error {
		// your middleware logic
		return next(ctx)
	}
})
```
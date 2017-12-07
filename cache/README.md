# goutils/cache

Install `github.com/pangpanglabs/goutils/cache` package.
```golang
go get -u github.com/pangpanglabs/goutils/cache
```

## Getting Started

Use `cache.Cache` interface type
```golang
var cache cache.Cache
```

Create local cache(use `sync.Map` as a cache storage):
```golang
cache = cache.Local{ExpireTime: time.Hour * 24}
```

Create redis cache(default expire time is `time.Hour * 24`):
```golang
redisConn := "redis://127.0.0.1:6379"
cache = cache.NewRedis(redisConn)
```

Or, create redis cache with addiotional options:
```golang
redisConn := "redis://127.0.0.1:6379"
cache := cache.NewRedis(redisConn, func(redis *cache.Redis) {
        redis.ExpireTime = time.Hour * 1
})
```

Save to cache:
```golang
loadFromCache, err := cache.LoadOrStore(key, &target,
        func() (interface{}, error) {
                // load and return value
        })
```

Delete from cache:
```golang
cache.Delete(key)
```

# goutils/cache

Install `github.com/pangpanglabs/goutils/cache` package.
```golang
go get -u github.com/pangpanglabs/goutils/cache
```

## Getting Started

Use `cache.Cache` interface type
```golang
var mycache cache.Cache
```

Create local cache(use `sync.Map` as a cache storage):
```golang
mycache = cache.Local{ExpireTime: time.Hour * 24}
```

Create redis cache
```golang
redisConn := "redis://127.0.0.1:6379"
mycache = cache.NewRedis(redisConn)
```
> - default expire time: `time.Hour * 24`
> - default `Converter`: `JsonConverter`

Or, create redis cache with addiotional options:
```golang
redisConn := "redis://127.0.0.1:6379"
mycache = cache.NewRedis(redisConn,
        cache.WithExpireTime(time.Hour*3),
        cache.WithGobConverter(),
        func(redis *cache.Redis) {
                // setup additional options
        },
)
```
> If you want to use `GobConverter`, you have to identify the concrete type of a value using `gob.Register()` function.

Save to cache:
```golang
loadFromCache, err := cache.LoadOrStore(key, &target, func() (interface{}, error) {
        // return your own result
        return "value", nil
})
```

Delete from cache:
```golang
cache.Delete(key)
```

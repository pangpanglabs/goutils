package cache_test

import (
	"os/exec"
	"testing"
	"time"

	"github.com/pangpanglabs/goutils/cache"
	"github.com/pangpanglabs/goutils/test"
)

// Prerequisite: `docker run --name redis -p 6379:6379 -d redis:alpine``
func TestRedis(t *testing.T) {
	redisCache := cache.NewRedis("redis://127.0.0.1:6379")

	var v interface{}

	t.Run("LoadOrStore", func(t *testing.T) {
		loadFromCache, err := redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
			return "value", nil
		})
		test.Ok(t, err)
		test.Equals(t, v, "value")
		test.Equals(t, loadFromCache, false)

		t.Run("StopRedis", func(t *testing.T) {
			// stop redis
			err := exec.Command("docker", "stop", "redis").Run()
			test.Ok(t, err)
			time.Sleep(time.Second)

			// try get value
			loadFromCache, err := redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
				return "value", nil
			})
			test.Ok(t, err)
			test.Equals(t, loadFromCache, false)
		})

		t.Run("RunRedis", func(t *testing.T) {
			// start redis
			err := exec.Command("docker", "start", "redis").Run()
			test.Ok(t, err)
			time.Sleep(time.Second)

			// try get value
			loadFromCache, err := redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
				return "value", nil
			})
			test.Ok(t, err)
			test.Equals(t, v, "value")
			test.Equals(t, loadFromCache, true)
		})
	})

	t.Run("Store/Load", func(t *testing.T) {
		redisCache.Store("key", "value2")

		time.Sleep(time.Second)

		loadFromCache := redisCache.Load("key", &v)
		test.Equals(t, v, "value2")
		test.Equals(t, loadFromCache, true)

		redisCache.Delete("key")

		loadFromCache = redisCache.Load("key", &v)
		test.Equals(t, loadFromCache, false)
	})

}

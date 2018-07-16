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

	// try get value
	_, err := redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
		return "value", nil
	})
	test.Ok(t, err)
	test.Equals(t, v, "value")

	// stop redis
	err = exec.Command("docker", "stop", "redis").Run()
	test.Ok(t, err)
	time.Sleep(time.Second)

	// try get value
	loadFromCache, err := redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
		return "value", nil
	})
	test.Ok(t, err)
	test.Equals(t, loadFromCache, false)

	// start redis
	err = exec.Command("docker", "start", "redis").Run()
	test.Ok(t, err)
	time.Sleep(time.Second)

	// try get value
	loadFromCache, err = redisCache.LoadOrStore("key", &v, func() (interface{}, error) {
		return "value", nil
	})
	test.Ok(t, err)
	test.Equals(t, v, "value")
	test.Equals(t, loadFromCache, true)
}

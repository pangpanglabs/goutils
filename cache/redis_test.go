package cache_test

import (
	"os"
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

func TestRedis_String(t *testing.T) {
	uri := os.Getenv("REDIS_CONN")
	redis := cache.NewRedis(uri, func(redis *cache.Redis) {
		redis.ExpireTime = time.Second * 60
	})
	key := "test:ping"
	value := StringPoint("")
	_, err := redis.LoadOrStore(key, value, func() (interface{}, error) {
		return "pong", nil
	})
	test.Ok(t, err)
	test.Equals(t, StringPoint("pong"), value)
	time.Sleep(10e9)
}
func TestRedis_Delete(t *testing.T) {
	uri := os.Getenv("REDIS_CONN")
	redis := cache.NewRedis(uri, func(redis *cache.Redis) {
		redis.ExpireTime = time.Second * 30
	})
	key := "test:ping"
	err := redis.Delete(key)
	test.Ok(t, err)
	time.Sleep(10e9)
}

func TestRedis_Struct(t *testing.T) {
	uri := os.Getenv("REDIS_CONN")
	redis := cache.NewRedis(uri, func(redis *cache.Redis) {
		redis.ExpireTime = time.Second * 30
	})
	key := "test:star"
	value := new(StarDto)
	_, err := redis.LoadOrStore(key, value, func() (interface{}, error) {
		return StarDto{
			Name: "Kwone Sang Woo",
		}, nil
	})
	test.Ok(t, err)
	test.Equals(t, "Kwone Sang Woo", value.Name)
	time.Sleep(10e9)
}

type StarDto struct {
	Name string `json:"name"`
}

func StringPoint(flag string) *string {
	return func(b string) *string { return &b }(flag)
}

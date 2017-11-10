package cache_test

import (
	"testing"
	"time"

	"github.com/pangpanglabs/goutils/cache"
	"github.com/pangpanglabs/goutils/test"
)

func TestLocal(t *testing.T) {
	c := cache.Local{
		ExpireTime: time.Second,
	}

	var v map[string]interface{}
	loadFromCache, err := c.LoadOrStore("a", &v, func() (interface{}, error) {
		return map[string]interface{}{
			"AA": "BB",
			"11": 22,
		}, nil
	})
	test.Ok(t, err)
	test.Equals(t, loadFromCache, false)
	test.Equals(t, v["AA"], "BB")
	test.Equals(t, v["11"], 22)

	loadFromCache, err = c.LoadOrStore("a", &v, func() (interface{}, error) {
		return map[string]interface{}{
			"AA": "BB",
			"11": 22,
		}, nil
	})
	test.Ok(t, err)
	test.Equals(t, loadFromCache, true)
	test.Equals(t, v["AA"], "BB")
	test.Equals(t, v["11"], 22)

	<-time.After(time.Second * 2)

	loadFromCache, err = c.LoadOrStore("a", &v, func() (interface{}, error) {
		return map[string]interface{}{
			"AA": "BB",
			"11": 22,
		}, nil
	})
	test.Ok(t, err)
	test.Equals(t, loadFromCache, false)
	test.Equals(t, v["AA"], "BB")
	test.Equals(t, v["11"], 22)
}

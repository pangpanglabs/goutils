package cache

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Local struct {
	m          sync.Map
	ExpireTime time.Duration
}

func (c *Local) LoadOrStore(key string, value interface{}, getter func() (interface{}, error)) (loadFromCache bool, err error) {
	result, ok := c.m.Load(key)
	if ok {
		if err := writeTo(result, value); err != nil {
			return false, err
		}
		return true, nil
	}

	result, err = getter()
	if err != nil {
		return false, err
	}

	c.m.Store(key, result)

	go func() {
		if c.ExpireTime > 0 {
			<-time.After(c.ExpireTime)
			c.m.Delete(key)
			logrus.WithField("key", key).Info("DELETE FROM LOCAL CACHE")
		}

	}()

	if err := writeTo(result, value); err != nil {
		return false, err
	}
	return false, nil
}

func (c *Local) Delete(key string) error {
	c.m.Delete(key)
	return nil
}

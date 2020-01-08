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
func (c *Local) Load(key string, value interface{}) (ok bool) {
	result, ok := c.m.Load(key)
	if !ok {
		return false
	}
	if err := writeTo(result, value); err != nil {
		logrus.WithFields(logrus.Fields{
			"key":    key,
			"result": result,
		}).WithError(err).Error("Fail to write to value")
		return false
	}
	return true
}
func (c *Local) Store(key string, value interface{}) {
	c.m.Store(key, value)
}

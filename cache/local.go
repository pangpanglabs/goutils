package cache

import (
	"sync"
)

type Local struct {
	m sync.Map
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

	if err := writeTo(result, value); err != nil {
		return false, err
	}
	return false, nil
}

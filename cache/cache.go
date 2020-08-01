package cache

import (
	"errors"
	"reflect"
)

type Cache interface {
	LoadOrStore(key string, value interface{}, getter func() (interface{}, error)) (loadFromCache bool, err error)
	Delete(key string) error
	Load(key string, value interface{}) (ok bool)
	Store(key string, value interface{})
}

func writeTo(data, dest interface{}) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return errors.New("needs a pointer to a value")
	} else if value.Elem().Kind() == reflect.Ptr {
		return errors.New("a pointer to a pointer is not allowed")
	}

	if v := value.Elem(); v.CanSet() {
		v.Set(reflect.ValueOf(data))
		return nil
	}

	return errors.New("cannot set value")
}

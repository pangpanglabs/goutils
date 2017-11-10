package cache

import (
	"errors"
	"reflect"
)

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

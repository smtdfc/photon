package photon

import (
	"fmt"
	"reflect"
)

func assertNotNil(val any, name string) {
	if val == nil {
		panic(fmt.Sprintf("%s must not be nil (val == nil)", name))
	}

	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		if v.IsNil() {
			panic(fmt.Sprintf("%s must not be nil (value is nil)", name))
		}
	}
}

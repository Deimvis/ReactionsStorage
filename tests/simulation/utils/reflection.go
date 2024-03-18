package utils

import "reflect"

func AssertPtr(v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("not pointer")
	}
}

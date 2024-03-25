package utils

import (
	"reflect"
	"runtime"
)

func AssertPtr(v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		panic("not pointer")
	}
}

func GetFnName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

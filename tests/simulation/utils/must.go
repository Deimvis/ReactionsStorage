package utils

import (
	"errors"
	"fmt"
)

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func Mustf[T any](v T, err error, msgAndArgs ...interface{}) T {
	if err != nil {
		msg := messageFromMsgAndArgs(append(msgAndArgs, err)...)
		panic(errors.New(msg))
	}
	return v
}

// https://github.com/stretchr/testify/blob/v1.9.0/assert/assertions.go#L280
func messageFromMsgAndArgs(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		return fmt.Sprintf(msgAndArgs[0].(string), msgAndArgs[1:]...)
	}
	return ""
}

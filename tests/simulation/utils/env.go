package utils

import (
	"os"
	"strings"
)

func Getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func IsDebugEnv() bool {
	trueOpts := []string{"1", "true"}
	return Contains(trueOpts, strings.ToLower(os.Getenv("DEBUG")))
}

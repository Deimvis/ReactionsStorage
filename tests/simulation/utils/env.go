package utils

import (
	"os"
	"strings"
)

func IsDebugEnv() bool {
	trueOpts := []string{"1", "true"}
	return Contains(trueOpts, strings.ToLower(os.Getenv("DEBUG")))
}

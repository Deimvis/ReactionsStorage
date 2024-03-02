package thelpers

import (
	"log"
	"os"
)

// should be used only in TestMain

// fails if !expr
func Check(expr bool, format string, v ...any) {
	if !expr {
		log.Fatalf(format, v...)
	}
}

// fails if env variable is not set
func CheckEnv(envVar string) {
	_, present := os.LookupEnv(envVar)
	if !present {
		log.Fatalf("Env variable `%s` is not set", envVar)
	}
}

package util

import (
	"log"
	"os"
	"strconv"
)

func GetEnvVarOrPanic(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Can't find env variable %s", key)
	}

	return val
}

func GetEnvInt64VarOrPanic(key string) int64 {
	val := GetEnvVarOrPanic(key)
	valAsInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		log.Fatalf("Invalid env variable %s", key)
	}

	return valAsInt
}

func GetEnvIntVarOrPanic(key string) int {
	return int(GetEnvInt64VarOrPanic(key))
}

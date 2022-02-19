package config

import (
	"os"
	"strconv"
)

func GetStr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Error to get " + key)
	}

	return val
}

func GetInt(key string) int {
	val := os.Getenv(key)
	if val == "" {
		panic("Error to get" + key)
	}

	val_int, err := strconv.Atoi(val)
	if err != nil {
		panic("Error to get" + key)
	}

	return val_int
}

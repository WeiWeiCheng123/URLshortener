package config

import (
	"os"
	"strconv"
)

//Get env variable and return as a string
func GetStr(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("Error to get " + key)
	}

	return val
}

//Get env variable and return as a Int
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

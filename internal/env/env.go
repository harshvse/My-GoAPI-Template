package env

import (
	"log"
	"os"
	"strconv"
)

// read the enviornment variable and return the string other types are given below
func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)

	if err != nil {
		log.Printf("The Environment Variable %s is not an integer", key)
		return fallback
	}
	return valAsInt
}

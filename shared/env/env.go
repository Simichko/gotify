package env

import (
	"os"
	"strconv"
)

func GetEnvAsInt(key string, defaultVal int) int {
	valueStr := GetEnv(key, "")
	value, err := strconv.Atoi(valueStr)
	if err == nil {
		return value
	}

	return defaultVal
}

func GetEnv(key string, defaultVal string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return defaultVal
}

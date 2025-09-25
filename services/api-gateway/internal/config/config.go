package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port            int
	EmailServiceUrl string
}

const emailUrlDefault = "amqp://guest:guest@rabbitmq:5672/"

func New() Config {
	return Config{
		Port:            getEnvAsInt("PORT", 8080),
		EmailServiceUrl: getEnv("EMAIL_SERVICE_URL", emailUrlDefault),
	}
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

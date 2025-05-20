package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type env struct {
	Port int
}

func initEnvVariables() env {
	godotenv.Load()

	return env{
		Port: getIntEnv("PORT", 8080),
	}
}

func getIntEnv(key string, fallback int) int {
	if stringValue, ok := os.LookupEnv(key); ok {
		value, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return fallback
		}

		return int(value)
	}

	return fallback
}

var Env = initEnvVariables()

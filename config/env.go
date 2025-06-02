package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type env struct {
	ApiPort              int
	DbUser               string
	DbPassword           string
	DbAddress            string
	DbName               string
	JWTExpirationSeconds int
	JWTSecret            string
}

func initEnvVariables() env {
	godotenv.Load()

	return env{
		ApiPort:              getIntEnv("API_PORT", 8080),
		DbUser:               getEnv("DB_USER", "development"),
		DbPassword:           getEnv("DB_PASSWORD", "db12345"),
		DbName:               getEnv("DB_NAME", "help"),
		DbAddress:            fmt.Sprintf("%s:%d", getEnv("DB_HOST", "localhost"), getIntEnv("DB_PORT", 3306)),
		JWTExpirationSeconds: getIntEnv("JWT_EXPIRATION_SECONDS", 24*60*60),
		JWTSecret:            getEnv("JWT_SECRET", "4d065eb6-e9dc-4aed-8709-30dbbcd8b1a5"),
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

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

var Env = initEnvVariables()

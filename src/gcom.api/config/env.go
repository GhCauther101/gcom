package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port string

	DbUser                 string
	DbPassword             string
	DbAddress              string
	DbName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Host:                   getEnv("PUB_HOST", "http://localhost"),
		Port:                   getEnv("PUB_PORT", "8080"),
		DbUser:                 getEnv("DB_USER", "username"),
		DbPassword:             getEnv("DB_PASSWORD", "1qaz!QAZ"),
		DbAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DbName:                 getEnv("DB_NAME", "gcom"),
		JWTSecret:              getEnv("JWT_SECRET", "random_secret_escape_this_pattern"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", (3600 * 24 * 7)),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}

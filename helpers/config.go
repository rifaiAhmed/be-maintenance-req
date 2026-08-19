package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

var Env = map[string]string{}

func SetupConfig() {
	fileEnv, err := godotenv.Read(".env")
	if err == nil {
		Env = fileEnv
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	if value := Env[key]; value != "" {
		return value
	}
	return fallback
}

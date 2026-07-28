package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Carrega as variáveis do arquivo .env se ele existir.
	// Não retornamos erro fatal se o arquivo não existir, pois em ambientes
	// de produção/Docker as variáveis podem ser passadas diretamente no container.
	_ = godotenv.Load()
}

func GetToken(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("variável de ambiente não encontrada: ", key)
	}
	return value
}

func GetString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return fallback
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return intValue
}

func GetDuration(key string, fallback time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return duration
}

package config

import (
	"os"
)

// Config содержит конфигурацию приложения
type Config struct {
	Port   string
	DBPath string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	config := &Config{
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "bookshelf.db"),
	}
	return config
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

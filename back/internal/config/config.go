package config

import (
	"log"
	"os"
)

type Config struct {
	DB  DBConfig
	JWT JWTConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	TTL string
}

func Load() *Config {
	return &Config{
		DB:  loadDBConfig(),
		JWT: loadJWTConfig(),
	}
}

func loadDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "app_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	log.Printf("env %s not set, using default: %s", key, fallback)
	return fallback
}

func loadJWTConfig() JWTConfig {
	return JWTConfig{
		Secret: getEnv("JWT_SECRET", "change_me"),
		TTL:    getEnv("JWT_TTL", "24h"),
	}
}
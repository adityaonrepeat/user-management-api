package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultDatabaseURL = "postgres://postgres:postgres@localhost:5432/user_management?sslmode=disable"
	defaultServerPort  = "3000"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: getEnv("DATABASE_URL", defaultDatabaseURL),
		ServerPort:  getEnv("SERVER_PORT", defaultServerPort),
	}

	if _, err := strconv.Atoi(cfg.ServerPort); err != nil {
		return Config{}, fmt.Errorf("SERVER_PORT must be numeric, got %q", cfg.ServerPort)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

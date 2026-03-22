package config

import "os"

type Config struct {
	PostgresURL string
	Port        string
	JWTSecret   string
}

func Load() Config {
	cfg := Config{
		PostgresURL: "postgres://audit:auditpass@localhost:5432/auditdb?sslmode=disable",
		Port:        ":8081",
		JWTSecret:   "jwt-secret-change-in-production-must-be-32-chars!",
	}

	if v := os.Getenv("POSTGRES_URL"); v != "" {
		cfg.PostgresURL = v
	}
	if v := os.Getenv("API_PORT"); v != "" {
		cfg.Port = ":" + v
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWTSecret = v
	}

	return cfg
}

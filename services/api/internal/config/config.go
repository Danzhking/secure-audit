package config

import "os"

type Config struct {
	PostgresURL string
	Port        string
}

func Load() Config {
	cfg := Config{
		PostgresURL: "postgres://audit:auditpass@localhost:5432/auditdb?sslmode=disable",
		Port:        ":8081",
	}

	if v := os.Getenv("POSTGRES_URL"); v != "" {
		cfg.PostgresURL = v
	}
	if v := os.Getenv("API_PORT"); v != "" {
		cfg.Port = ":" + v
	}

	return cfg
}

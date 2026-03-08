package config

import "os"

type Config struct {
	RabbitURL   string
	PostgresURL string
}

func Load() Config {
	cfg := Config{
		RabbitURL:   "amqp://guest:guest@localhost:5672/",
		PostgresURL: "postgres://audit:auditpass@localhost:5432/auditdb?sslmode=disable",
	}

	if v := os.Getenv("RABBIT_URL"); v != "" {
		cfg.RabbitURL = v
	}
	if v := os.Getenv("POSTGRES_URL"); v != "" {
		cfg.PostgresURL = v
	}

	return cfg
}

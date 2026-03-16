package config

import (
	"os"
	"strings"
)

type Config struct {
	RabbitURL  string
	Port       string
	APIKeys    []string
	HMACSecret string
}

func Load() Config {
	cfg := Config{
		RabbitURL:  "amqp://guest:guest@localhost:5672/",
		Port:       ":8080",
		APIKeys:    []string{"service-key-auth-001", "service-key-files-002"},
		HMACSecret: "super-secret-hmac-key-change-in-production",
	}

	if v := os.Getenv("RABBIT_URL"); v != "" {
		cfg.RabbitURL = v
	}
	if v := os.Getenv("COLLECTOR_PORT"); v != "" {
		cfg.Port = ":" + v
	}
	if v := os.Getenv("API_KEYS"); v != "" {
		cfg.APIKeys = strings.Split(v, ",")
	}
	if v := os.Getenv("HMAC_SECRET"); v != "" {
		cfg.HMACSecret = v
	}

	return cfg
}

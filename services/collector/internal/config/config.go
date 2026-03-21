package config

import (
	"os"
	"strings"
)

type Config struct {
	RabbitURL  string
	Port       string
	TLSPort    string
	TLSCert    string
	TLSKey     string
	APIKeys    []string
	HMACSecret string
}

func Load() Config {
	cfg := Config{
		RabbitURL:  "amqp://guest:guest@localhost:5672/",
		Port:       ":8080",
		TLSPort:    ":8443",
		TLSCert:    "",
		TLSKey:     "",
		APIKeys:    []string{"service-key-auth-001", "service-key-files-002"},
		HMACSecret: "super-secret-hmac-key-change-in-production",
	}

	if v := os.Getenv("RABBIT_URL"); v != "" {
		cfg.RabbitURL = v
	}
	if v := os.Getenv("COLLECTOR_PORT"); v != "" {
		cfg.Port = ":" + v
	}
	if v := os.Getenv("TLS_CERT"); v != "" {
		cfg.TLSCert = v
	}
	if v := os.Getenv("TLS_KEY"); v != "" {
		cfg.TLSKey = v
	}
	if v := os.Getenv("API_KEYS"); v != "" {
		cfg.APIKeys = strings.Split(v, ",")
	}
	if v := os.Getenv("HMAC_SECRET"); v != "" {
		cfg.HMACSecret = v
	}

	return cfg
}

func (c Config) TLSEnabled() bool {
	return c.TLSCert != "" && c.TLSKey != ""
}

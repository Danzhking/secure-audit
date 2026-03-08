package config

type Config struct {
	PostgresURL string
	RabbitURL   string
}

func Load() Config {
	return Config{
		PostgresURL: "postgres://audit:auditpass@localhost:5432/auditdb?sslmode=disable",
		RabbitURL:   "amqp://guest:guest@localhost:5672/",
	}
}

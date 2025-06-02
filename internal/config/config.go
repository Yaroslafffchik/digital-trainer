package config

import (
	"os"
)

type Config struct {
	HTTPPort      string
	TCPPort       string
	UDPPort       string
	WebSocketPort string
	DatabaseURL   string
}

func Load() (*Config, error) {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", ":8080"),
		TCPPort:       getEnv("TCP_PORT", ":1234"),
		UDPPort:       getEnv("UDP_PORT", ":1235"),
		WebSocketPort: getEnv("WS_PORT", ":1236"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:CAXAPOK2005ya@localhost:5432/digital_trainer?sslmode=disable"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

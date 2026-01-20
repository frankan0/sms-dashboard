package config

import (
	"os"
)

type Config struct {
	Port         string
	APIToken     string // For Android SMS Post
	JWTSecret    string // For Dashboard Login
	DatabasePath string
	Secret       string
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		APIToken:     getEnv("API_TOKEN", "uyejxhhshfe"),
		JWTSecret:    getEnv("JWT_SECRET", "jwdesxxdf817D"),
		DatabasePath: getEnv("DB_PATH", "sms.db"),
		Secret:       getEnv("SECRET", "smsadmin"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

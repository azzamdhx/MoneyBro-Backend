package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	ResendAPIKey string
	FrontendURL  string
	Env          string
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		Env:          getEnv("ENV", "development"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		RedisURL:     getEnv("REDIS_URL", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		ResendAPIKey: getEnv("RESEND_API_KEY", ""),
		FrontendURL:  getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

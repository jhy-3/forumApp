package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	Environment string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret                   string
	JWTAccessTokenExpireMinutes int
	JWTRefreshTokenExpireDays   int

	CORSOrigins string
}

func NewConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "forumuser"),
		DBPassword: getEnv("DB_PASSWORD", "forumpass123"),
		DBName:     getEnv("DB_NAME", "forumdb"),

		JWTSecret:                   getEnv("JWT_SECRET", "your-secret-key"),
		JWTAccessTokenExpireMinutes: getEnvAsInt("JWT_ACCESS_TOKEN_EXPIRE_MINUTES", 15),
		JWTRefreshTokenExpireDays:   getEnvAsInt("JWT_REFRESH_TOKEN_EXPIRE_DAYS", 7),

		CORSOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}


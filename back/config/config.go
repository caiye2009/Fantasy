package config

import (
	"os"
)

type Config struct {
	// Database
	DatabaseDSN string

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Elasticsearch
	ESAddress   string
	ESAddresses []string
	ESUsername  string
	ESPassword  string

	// Casbin
	CasbinModelPath string

	// JWT
	JWTSecret string

	// Server
	ServerAddr string

	// Log
	LogLevel  string
	LogFormat string
}

func LoadConfig() *Config {
	return &Config{
		// PostgreSQL DSN
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=fantasy port=5432 sslmode=disable TimeZone=Asia/Shanghai"),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,

		// Elasticsearch
		ESAddress:   getEnv("ES_ADDRESS", "http://localhost:9200"),
		ESAddresses: []string{getEnv("ES_ADDRESS", "http://localhost:9200")},
		ESUsername:  getEnv("ES_USERNAME", "elastic"),
		ESPassword:  getEnv("ES_PASSWORD", "123456"),

		// Casbin
		CasbinModelPath: getEnv("CASBIN_MODEL_PATH", "./config/casbin_model.conf"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),

		// Server
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),

		// Log
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "console"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
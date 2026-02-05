package config

import (
	"os"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN string
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret   string
	Issuer   string
	Audience string
}

// Load 加载配置
func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Addr: getEnv("SERVER_ADDR", ":8085"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DB_DSN", "root:root@tcp(localhost:3306)/jimureport?charset=utf8mb4&parseTime=True&loc=Local"),
		},
		JWT: JWTConfig{
			Secret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Issuer:   getEnv("JWT_ISSUER", "jimureport"),
			Audience: getEnv("JWT_AUDIENCE", "jimureport"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

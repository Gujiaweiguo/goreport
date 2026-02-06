package config

import (
	"fmt"
	"os"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Cache    CacheConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Addr string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // seconds
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled  bool
	Addr     string
	Password string
	DB       int
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
			DSN:             getEnv("DB_DSN", "root:root@tcp(localhost:3306)/jimureport?charset=utf8mb4&parseTime=True&loc=Local"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: getIntEnv("DB_CONN_MAX_LIFETIME", 3600),
		},
		JWT: JWTConfig{
			Secret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Issuer:   getEnv("JWT_ISSUER", "jimureport"),
			Audience: getEnv("JWT_AUDIENCE", "jimureport"),
		},
		Cache: CacheConfig{
			Enabled:  getBoolEnv("CACHE_ENABLED", false),
			Addr:     getEnv("CACHE_ADDR", "localhost:6379"),
			Password: getEnv("CACHE_PASSWORD", ""),
			DB:       getIntEnv("CACHE_DB", 0),
		},
	}, nil
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

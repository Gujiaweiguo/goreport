package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clear all relevant environment variables
	clearConfigEnvVars(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Test server defaults
	if cfg.Server.Addr != ":8085" {
		t.Errorf("Server.Addr = %q, want :8085", cfg.Server.Addr)
	}

	// Test database defaults
	if cfg.Database.DSN != "root:root@tcp(localhost:3306)/goreport?charset=utf8mb4&parseTime=True&loc=Local" {
		t.Errorf("Database.DSN = %q, want default DSN", cfg.Database.DSN)
	}
	if cfg.Database.MaxOpenConns != 100 {
		t.Errorf("Database.MaxOpenConns = %d, want 100", cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns != 10 {
		t.Errorf("Database.MaxIdleConns = %d, want 10", cfg.Database.MaxIdleConns)
	}
	if cfg.Database.ConnMaxLifetime != 3600 {
		t.Errorf("Database.ConnMaxLifetime = %d, want 3600", cfg.Database.ConnMaxLifetime)
	}

	// Test JWT defaults
	if cfg.JWT.Issuer != "goreport" {
		t.Errorf("JWT.Issuer = %q, want goreport", cfg.JWT.Issuer)
	}
	if cfg.JWT.Audience != "goreport" {
		t.Errorf("JWT.Audience = %q, want goreport", cfg.JWT.Audience)
	}
	// JWT Secret should be generated if not set
	if cfg.JWT.Secret == "" {
		t.Error("JWT.Secret should not be empty")
	}
	if len(cfg.JWT.Secret) != 64 { // 32 bytes hex encoded = 64 chars
		t.Errorf("JWT.Secret length = %d, want 64", len(cfg.JWT.Secret))
	}

	// Test cache defaults
	if cfg.Cache.Enabled != false {
		t.Errorf("Cache.Enabled = %v, want false", cfg.Cache.Enabled)
	}
	if cfg.Cache.Addr != "localhost:6379" {
		t.Errorf("Cache.Addr = %q, want localhost:6379", cfg.Cache.Addr)
	}
	if cfg.Cache.Password != "" {
		t.Errorf("Cache.Password = %q, want empty", cfg.Cache.Password)
	}
	if cfg.Cache.DB != 0 {
		t.Errorf("Cache.DB = %d, want 0", cfg.Cache.DB)
	}
	if cfg.Cache.DefaultTTL != 3600 {
		t.Errorf("Cache.DefaultTTL = %d, want 3600", cfg.Cache.DefaultTTL)
	}
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	testCases := map[string]string{
		"SERVER_ADDR":          ":9090",
		"DB_DSN":               "user:pass@tcp(localhost:3307)/testdb",
		"DB_MAX_OPEN_CONNS":    "50",
		"DB_MAX_IDLE_CONNS":    "5",
		"DB_CONN_MAX_LIFETIME": "1800",
		"JWT_SECRET":           "test-secret-key",
		"JWT_ISSUER":           "test-issuer",
		"JWT_AUDIENCE":         "test-audience",
		"CACHE_ENABLED":        "true",
		"CACHE_ADDR":           "redis:6379",
		"CACHE_PASSWORD":       "redis-pass",
		"CACHE_DB":             "1",
		"CACHE_DEFAULT_TTL":    "7200",
	}

	// Set environment variables
	for key, value := range testCases {
		os.Setenv(key, value)
	}
	defer clearConfigEnvVars(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify server config
	if cfg.Server.Addr != ":9090" {
		t.Errorf("Server.Addr = %q, want :9090", cfg.Server.Addr)
	}

	// Verify database config
	if cfg.Database.DSN != "user:pass@tcp(localhost:3307)/testdb" {
		t.Errorf("Database.DSN = %q, want user:pass@tcp(localhost:3307)/testdb", cfg.Database.DSN)
	}
	if cfg.Database.MaxOpenConns != 50 {
		t.Errorf("Database.MaxOpenConns = %d, want 50", cfg.Database.MaxOpenConns)
	}
	if cfg.Database.MaxIdleConns != 5 {
		t.Errorf("Database.MaxIdleConns = %d, want 5", cfg.Database.MaxIdleConns)
	}
	if cfg.Database.ConnMaxLifetime != 1800 {
		t.Errorf("Database.ConnMaxLifetime = %d, want 1800", cfg.Database.ConnMaxLifetime)
	}

	// Verify JWT config
	if cfg.JWT.Secret != "test-secret-key" {
		t.Errorf("JWT.Secret = %q, want test-secret-key", cfg.JWT.Secret)
	}
	if cfg.JWT.Issuer != "test-issuer" {
		t.Errorf("JWT.Issuer = %q, want test-issuer", cfg.JWT.Issuer)
	}
	if cfg.JWT.Audience != "test-audience" {
		t.Errorf("JWT.Audience = %q, want test-audience", cfg.JWT.Audience)
	}

	// Verify cache config
	if cfg.Cache.Enabled != true {
		t.Errorf("Cache.Enabled = %v, want true", cfg.Cache.Enabled)
	}
	if cfg.Cache.Addr != "redis:6379" {
		t.Errorf("Cache.Addr = %q, want redis:6379", cfg.Cache.Addr)
	}
	if cfg.Cache.Password != "redis-pass" {
		t.Errorf("Cache.Password = %q, want redis-pass", cfg.Cache.Password)
	}
	if cfg.Cache.DB != 1 {
		t.Errorf("Cache.DB = %d, want 1", cfg.Cache.DB)
	}
	if cfg.Cache.DefaultTTL != 7200 {
		t.Errorf("Cache.DefaultTTL = %d, want 7200", cfg.Cache.DefaultTTL)
	}
}

func TestGetEnv(t *testing.T) {
	// Test with environment variable set
	os.Setenv("TEST_GET_ENV", "test-value")
	defer os.Unsetenv("TEST_GET_ENV")

	result := getEnv("TEST_GET_ENV", "default")
	if result != "test-value" {
		t.Errorf("getEnv() = %q, want test-value", result)
	}

	// Test with environment variable not set
	result = getEnv("NON_EXISTENT_ENV", "default-value")
	if result != "default-value" {
		t.Errorf("getEnv() = %q, want default-value", result)
	}
}

func TestGetIntEnv(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT_ENV", "42")
	defer os.Unsetenv("TEST_INT_ENV")

	result := getIntEnv("TEST_INT_ENV", 0)
	if result != 42 {
		t.Errorf("getIntEnv() = %d, want 42", result)
	}

	// Test with invalid integer
	os.Setenv("TEST_INT_ENV_INVALID", "not-a-number")
	defer os.Unsetenv("TEST_INT_ENV_INVALID")

	result = getIntEnv("TEST_INT_ENV_INVALID", 100)
	if result != 100 {
		t.Errorf("getIntEnv() with invalid value = %d, want 100 (default)", result)
	}

	// Test with environment variable not set
	result = getIntEnv("NON_EXISTENT_INT_ENV", 999)
	if result != 999 {
		t.Errorf("getIntEnv() = %d, want 999", result)
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		want         bool
	}{
		{"true value", "true", false, true},
		{"1 value", "1", false, true},
		{"false value", "false", true, false},
		{"0 value", "0", true, false},
		{"other value", "random", true, false},
		{"empty uses default", "", true, true},
		{"empty uses default false", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_BOOL_ENV", tt.envValue)
				defer os.Unsetenv("TEST_BOOL_ENV")
			} else {
				os.Unsetenv("TEST_BOOL_ENV")
			}

			result := getBoolEnv("TEST_BOOL_ENV", tt.defaultValue)
			if result != tt.want {
				t.Errorf("getBoolEnv() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestGenerateRandomSecret(t *testing.T) {
	// Test that secrets are generated
	secret1 := generateRandomSecret()
	if secret1 == "" {
		t.Error("generateRandomSecret() returned empty string")
	}

	// Test length (32 bytes = 64 hex chars)
	if len(secret1) != 64 {
		t.Errorf("generateRandomSecret() length = %d, want 64", len(secret1))
	}

	// Test that secrets are unique
	secret2 := generateRandomSecret()
	if secret1 == secret2 {
		t.Error("generateRandomSecret() generated identical secrets")
	}

	// Test that secrets are valid hex
	for i, c := range secret1 {
		if !isHexChar(c) {
			t.Errorf("generateRandomSecret() contains non-hex character at position %d: %c", i, c)
		}
	}
}

func TestConfigStruct(t *testing.T) {
	// Test that Config struct can be created and fields are accessible
	cfg := &Config{
		Server: ServerConfig{
			Addr: ":8080",
		},
		Database: DatabaseConfig{
			DSN:             "test-dsn",
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300,
		},
		JWT: JWTConfig{
			Secret:   "secret",
			Issuer:   "issuer",
			Audience: "audience",
		},
		Cache: CacheConfig{
			Enabled:    true,
			Addr:       "localhost:6379",
			Password:   "pass",
			DB:         1,
			DefaultTTL: 1800,
		},
	}

	if cfg.Server.Addr != ":8080" {
		t.Errorf("Config.Server.Addr = %q, want :8080", cfg.Server.Addr)
	}
	if cfg.Database.DSN != "test-dsn" {
		t.Errorf("Config.Database.DSN = %q, want test-dsn", cfg.Database.DSN)
	}
	if cfg.JWT.Secret != "secret" {
		t.Errorf("Config.JWT.Secret = %q, want secret", cfg.JWT.Secret)
	}
	if !cfg.Cache.Enabled {
		t.Error("Config.Cache.Enabled should be true")
	}
}

// Helper function to clear config-related environment variables
func clearConfigEnvVars(t *testing.T) {
	t.Helper()
	envVars := []string{
		"SERVER_ADDR",
		"DB_DSN",
		"DB_MAX_OPEN_CONNS",
		"DB_MAX_IDLE_CONNS",
		"DB_CONN_MAX_LIFETIME",
		"JWT_SECRET",
		"JWT_ISSUER",
		"JWT_AUDIENCE",
		"CACHE_ENABLED",
		"CACHE_ADDR",
		"CACHE_PASSWORD",
		"CACHE_DB",
		"CACHE_DEFAULT_TTL",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

// Helper to check if character is a valid hex character
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func TestLoad_JWTSecretFromEnvironment(t *testing.T) {
	clearConfigEnvVars(t)
	os.Setenv("JWT_SECRET", "my-custom-secret")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.JWT.Secret != "my-custom-secret" {
		t.Errorf("JWT.Secret = %q, want my-custom-secret", cfg.JWT.Secret)
	}
}

func TestLoad_JWTSecretAutoGenerated(t *testing.T) {
	clearConfigEnvVars(t)

	// Generate multiple secrets and ensure they are unique
	secrets := make(map[string]bool)
	for i := 0; i < 10; i++ {
		cfg, err := Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if secrets[cfg.JWT.Secret] {
			t.Error("generateRandomSecret() generated duplicate secret")
		}
		secrets[cfg.JWT.Secret] = true
	}
}

func TestGetIntEnv_LargeNumbers(t *testing.T) {
	os.Setenv("TEST_LARGE_INT", "999999")
	defer os.Unsetenv("TEST_LARGE_INT")

	result := getIntEnv("TEST_LARGE_INT", 0)
	if result != 999999 {
		t.Errorf("getIntEnv() = %d, want 999999", result)
	}
}

func TestGetIntEnv_NegativeNumbers(t *testing.T) {
	os.Setenv("TEST_NEGATIVE_INT", "-1")
	defer os.Unsetenv("TEST_NEGATIVE_INT")

	result := getIntEnv("TEST_NEGATIVE_INT", 0)
	if result != -1 {
		t.Errorf("getIntEnv() = %d, want -1", result)
	}
}

func TestGetBoolEnv_TrueVariants(t *testing.T) {
	trueValues := []string{"true", "TRUE", "True", "1"}
	for _, v := range trueValues {
		os.Setenv("TEST_BOOL_VARIANT", v)
		result := getBoolEnv("TEST_BOOL_VARIANT", false)
		// Only lowercase "true" and "1" should be true
		expected := v == "true" || v == "1"
		if result != expected {
			t.Errorf("getBoolEnv(%q) = %v, want %v", v, result, expected)
		}
	}
	os.Unsetenv("TEST_BOOL_VARIANT")
}

func TestLoad_EmptyStringEnvVars(t *testing.T) {
	clearConfigEnvVars(t)

	// Set some env vars to empty string (should use defaults)
	os.Setenv("SERVER_ADDR", "")
	os.Setenv("DB_DSN", "")
	defer os.Unsetenv("SERVER_ADDR")
	defer os.Unsetenv("DB_DSN")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Empty string should still use default
	if cfg.Server.Addr != ":8085" {
		t.Errorf("Server.Addr = %q, want :8085", cfg.Server.Addr)
	}
	if !strings.Contains(cfg.Database.DSN, "goreport") {
		t.Errorf("Database.DSN = %q, should contain 'goreport'", cfg.Database.DSN)
	}
}

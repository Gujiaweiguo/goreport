package datasource

import (
	"context"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
)

func TestProfileValidator(t *testing.T) {
	validator := NewProfileValidator()

	t.Run("ValidateMySQL", func(t *testing.T) {
		config := map[string]interface{}{
			"host":     "localhost",
			"port":     3306,
			"database": "test",
			"username": "root",
			"password": "password",
		}

		if err := validator.Validate("mysql", config); err != nil {
			t.Errorf("validation should succeed, got: %v", err)
		}
	})

	t.Run("ValidateMissingRequired", func(t *testing.T) {
		config := map[string]interface{}{
			"host": "localhost",
		}

		if err := validator.Validate("mysql", config); err == nil {
			t.Error("validation should fail for missing required fields")
		}
	})

	t.Run("ValidateInvalidMaxConnections", func(t *testing.T) {
		config := map[string]interface{}{
			"host":            "localhost",
			"port":            3306,
			"database":        "test",
			"max_connections": 101,
		}

		if err := validator.Validate("mysql", config); err == nil {
			t.Error("validation should fail for max_connections > 100")
		}
	})

	t.Run("ValidateInvalidQueryTimeout", func(t *testing.T) {
		config := map[string]interface{}{
			"host":                  "localhost",
			"port":                  3306,
			"database":              "test",
			"query_timeout_seconds": 301,
		}

		if err := validator.Validate("mysql", config); err == nil {
			t.Error("validation should fail for query_timeout_seconds > 300")
		}
	})

	t.Run("ValidateSSHTunnel", func(t *testing.T) {
		config := map[string]interface{}{
			"host":         "localhost",
			"port":         3306,
			"database":     "test",
			"ssh_host":     "bastion.example.com",
			"ssh_port":     22,
			"ssh_username": "user",
			"ssh_key":      "key-content",
		}

		if err := validator.Validate("mysql", config); err != nil {
			t.Errorf("SSH tunnel validation should succeed, got: %v", err)
		}
	})

	t.Run("ValidateSSHTunnelMissingPassword", func(t *testing.T) {
		config := map[string]interface{}{
			"host":         "localhost",
			"port":         3306,
			"database":     "test",
			"ssh_host":     "bastion.example.com",
			"ssh_port":     22,
			"ssh_username": "user",
		}

		if err := validator.Validate("mysql", config); err == nil {
			t.Error("validation should fail for SSH tunnel without password or key")
		}
	})

	t.Run("GetProfiles", func(t *testing.T) {
		profiles := validator.ListProfiles()

		if len(profiles) == 0 {
			t.Error("profiles list should not be empty")
		}

		mysqlProfile := findProfile(profiles, "mysql")
		if mysqlProfile.Name == "" {
			t.Error("mysql profile should exist")
		}

		if !mysqlProfile.SupportsSSH {
			t.Error("mysql profile should support SSH")
		}
	})
}

func findProfile(profiles []ConnectorProfile, name string) ConnectorProfile {
	for _, profile := range profiles {
		if profile.Name == name {
			return profile
		}
	}
	return ConnectorProfile{}
}

func TestConnectionBuilder(t *testing.T) {
	builder := NewConnectionBuilder()

	t.Run("BuildDSNWithoutSSH", func(t *testing.T) {
		ds := &models.DataSource{
			Host:     "localhost",
			Port:     3306,
			Database: "test",
			Username: "root",
			Password: "password",
		}

		dsn, tunnel, err := builder.BuildDSN(context.Background(), ds)
		if err != nil {
			t.Errorf("BuildDSN should succeed, got: %v", err)
		}

		if tunnel != nil {
			t.Error("tunnel should be nil for non-SSH connection")
		}

		if dsn == "" {
			t.Error("DSN should not be empty")
		}
	})

	t.Run("BuildDSNWithSSH", func(t *testing.T) {
		ds := &models.DataSource{
			Host:        "target-db.example.com",
			Port:        3306,
			Database:    "test",
			Username:    "root",
			Password:    "password",
			SSHHost:     "bastion.example.com",
			SSHPort:     22,
			SSHUsername: "user",
			SSHPassword: "ssh-password",
		}

		dsn, tunnel, err := builder.BuildDSN(context.Background(), ds)
		if err == nil {
			t.Fatal("BuildDSN should fail in test environment without reachable SSH bastion")
		}

		if tunnel != nil {
			t.Error("tunnel should be nil when SSH tunnel setup fails")
		}

		if dsn != "" {
			t.Error("dsn should be empty when SSH tunnel setup fails")
		}
	})
}

func TestSSHTunnel(t *testing.T) {
	t.Skip("Integration test - requires real SSH server")

	tunnel := NewSSHTunnel(&SSHTunnelConfig{
		Host:     "localhost",
		Port:     22,
		Username: "test",
		Password: "test",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	localAddr, err := tunnel.Connect(ctx, &SSHTunnelConfig{
		Host:     "localhost",
		Port:     22,
		Username: "test",
		Password: "test",
	}, "localhost", 3306)

	if err != nil {
		t.Skipf("SSH connection failed (expected in test env): %v", err)
		return
	}

	defer tunnel.Close()

	if localAddr == "" {
		t.Error("local address should not be empty")
	}
}

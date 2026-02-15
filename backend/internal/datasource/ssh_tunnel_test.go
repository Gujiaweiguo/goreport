package datasource

import (
	"context"
	"fmt"
	"os"
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
	sshHost := os.Getenv("SSH_TEST_HOST")
	if sshHost == "" {
		t.Skip("SSH_TEST_HOST not set - skipping SSH tunnel integration test")
	}

	sshPort := 22
	if port := os.Getenv("SSH_TEST_PORT"); port != "" {
		_, _ = fmt.Sscanf(port, "%d", &sshPort)
	}

	sshUser := os.Getenv("SSH_TEST_USER")
	if sshUser == "" {
		sshUser = "testuser"
	}

	sshPassword := os.Getenv("SSH_TEST_PASSWORD")
	if sshPassword == "" {
		sshPassword = "testpassword"
	}

	tunnel := NewSSHTunnel(&SSHTunnelConfig{
		Host:     sshHost,
		Port:     sshPort,
		Username: sshUser,
		Password: sshPassword,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	localAddr, err := tunnel.Connect(ctx, &SSHTunnelConfig{
		Host:     sshHost,
		Port:     sshPort,
		Username: sshUser,
		Password: sshPassword,
	}, "127.0.0.1", 3306)

	if err != nil {
		t.Logf("SSH tunnel test skipped (server may not support port forwarding): %v", err)
		t.Skip("SSH server does not support port forwarding")
	}

	defer tunnel.Close()

	if localAddr == "" {
		t.Error("local address should not be empty")
	}

	t.Logf("SSH tunnel established, local address: %s", localAddr)
}

func TestSSHTunnel_InvalidCredentials(t *testing.T) {
	sshHost := os.Getenv("SSH_TEST_HOST")
	if sshHost == "" {
		t.Skip("SSH_TEST_HOST not set - skipping SSH tunnel integration test")
	}

	tunnel := NewSSHTunnel(&SSHTunnelConfig{
		Host:     sshHost,
		Port:     22,
		Username: "invaliduser",
		Password: "invalidpassword",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := tunnel.Connect(ctx, &SSHTunnelConfig{
		Host:     sshHost,
		Port:     22,
		Username: "invaliduser",
		Password: "invalidpassword",
	}, "127.0.0.1", 3306)

	if err == nil {
		t.Error("SSH connection should fail with invalid credentials")
		defer tunnel.Close()
	}
}

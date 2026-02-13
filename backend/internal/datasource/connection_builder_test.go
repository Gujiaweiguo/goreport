package datasource

import (
	"context"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewConnectionBuilder(t *testing.T) {
	builder := NewConnectionBuilder()
	assert.NotNil(t, builder)
}

func TestConnectionBuilder_BuildDSN_NoSSH(t *testing.T) {
	builder := NewConnectionBuilder()
	ctx := context.Background()

	ds := &models.DataSource{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "testdb",
	}

	dsn, tunnel, err := builder.BuildDSN(ctx, ds)

	assert.NoError(t, err)
	assert.Nil(t, tunnel)
	assert.Contains(t, dsn, "root:password")
	assert.Contains(t, dsn, "localhost:3306")
	assert.Contains(t, dsn, "testdb")
}

func TestConnectionBuilder_BuildDSN_EmptyPassword(t *testing.T) {
	builder := NewConnectionBuilder()
	ctx := context.Background()

	ds := &models.DataSource{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "",
		Database: "testdb",
	}

	dsn, tunnel, err := builder.BuildDSN(ctx, ds)

	assert.NoError(t, err)
	assert.Nil(t, tunnel)
	assert.Contains(t, dsn, "root:@tcp")
}

func TestConnectionBuilder_BuildDSN_VariousPorts(t *testing.T) {
	builder := NewConnectionBuilder()
	ctx := context.Background()

	tests := []struct {
		port     int
		expected string
	}{
		{3306, "localhost:3306"},
		{3307, "localhost:3307"},
		{5432, "localhost:5432"},
		{80, "localhost:80"},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.port)), func(t *testing.T) {
			ds := &models.DataSource{
				Host:     "localhost",
				Port:     tt.port,
				Username: "root",
				Password: "password",
				Database: "testdb",
			}

			dsn, _, err := builder.BuildDSN(ctx, ds)

			assert.NoError(t, err)
			assert.Contains(t, dsn, tt.expected)
		})
	}
}

func TestParseHostPort(t *testing.T) {
	tests := []struct {
		name        string
		addr        string
		expectHost  string
		expectPort  int
		expectError bool
	}{
		{
			name:       "valid address",
			addr:       "localhost:3306",
			expectHost: "localhost",
			expectPort: 3306,
		},
		{
			name:       "valid with IP",
			addr:       "127.0.0.1:8080",
			expectHost: "127.0.0.1",
			expectPort: 8080,
		},
		{
			name:        "missing port",
			addr:        "localhost",
			expectError: true,
		},
		{
			name:        "invalid port",
			addr:        "localhost:abc",
			expectError: true,
		},
		{
			name:        "empty address",
			addr:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, err := parseHostPort(tt.addr)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectHost, host)
				assert.Equal(t, tt.expectPort, port)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{5, 5, 5},
		{0, 10, 0},
		{-1, 5, -1},
		{100, 50, 50},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConnectionBuilder_BuildDSN_SpecialCharacters(t *testing.T) {
	builder := NewConnectionBuilder()
	ctx := context.Background()

	ds := &models.DataSource{
		Host:     "localhost",
		Port:     3306,
		Username: "user@domain",
		Password: "p@ss:word",
		Database: "test_db",
	}

	dsn, _, err := builder.BuildDSN(ctx, ds)

	assert.NoError(t, err)
	assert.Contains(t, dsn, "user@domain")
	assert.Contains(t, dsn, "p@ss:word")
	assert.Contains(t, dsn, "test_db")
}

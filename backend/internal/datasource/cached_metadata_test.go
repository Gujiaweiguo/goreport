package datasource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDatabaseFromDSN(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		expected string
	}{
		{
			name:     "standard DSN with params",
			dsn:      "root:password@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True",
			expected: "mydb",
		},
		{
			name:     "DSN without params",
			dsn:      "root:password@tcp(localhost:3306)/mydb",
			expected: "mydb",
		},
		{
			name:     "DSN with underscore database",
			dsn:      "user:pass@tcp(host:3306)/my_database?params",
			expected: "my_database",
		},
		{
			name:     "DSN with hyphen database",
			dsn:      "user:pass@tcp(host:3306)/my-database",
			expected: "my-database",
		},
		{
			name:     "DSN without database",
			dsn:      "root:password@tcp(localhost:3306)",
			expected: "",
		},
		{
			name:     "empty DSN",
			dsn:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDatabaseFromDSN(tt.dsn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewCachedMetadataService(t *testing.T) {
	service := NewCachedMetadataService(nil)
	assert.NotNil(t, service)
	assert.Nil(t, service.cache)
}

func TestFieldInfo_Struct(t *testing.T) {
	field := FieldInfo{
		Name:     "id",
		Type:     "int",
		Nullable: false,
		Comment:  "Primary key",
	}

	assert.Equal(t, "id", field.Name)
	assert.Equal(t, "int", field.Type)
	assert.False(t, field.Nullable)
	assert.Equal(t, "Primary key", field.Comment)
}

func TestExtractDatabaseFromDSN_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		expected string
	}{
		{
			name:     "DSN with numbers in database",
			dsn:      "user:pass@tcp(host:3306)/db123",
			expected: "db123",
		},
		{
			name:     "DSN with uppercase database",
			dsn:      "user:pass@tcp(host:3306)/MYDATABASE",
			expected: "MYDATABASE",
		},
		{
			name:     "DSN with empty database",
			dsn:      "user:pass@tcp(host:3306)/",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDatabaseFromDSN(tt.dsn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

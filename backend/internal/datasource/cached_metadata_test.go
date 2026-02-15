package datasource

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDSNForCachedMetadata() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

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

func TestCachedMetadataService_GetTables_CacheHit(t *testing.T) {
	cfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cfg)
	service := NewCachedMetadataService(c)

	ctx := context.Background()
	tenantID := "tenant-1"
	datasourceID := "ds-1"
	dsn := getTestDSNForCachedMetadata()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	tables, err := service.GetTables(ctx, tenantID, datasourceID, dsn)

	require.NoError(t, err)
	assert.NotNil(t, tables)
}

func TestCachedMetadataService_GetFields_CacheHit(t *testing.T) {
	cfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cfg)
	service := NewCachedMetadataService(c)

	ctx := context.Background()
	tenantID := "tenant-1"
	datasourceID := "ds-1"
	dsn := getTestDSNForCachedMetadata()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	fields, err := service.GetFields(ctx, tenantID, datasourceID, dsn, "tenants")

	require.NoError(t, err)
	assert.NotNil(t, fields)
}

func TestCachedMetadataService_GetTables_InvalidDSN(t *testing.T) {
	cfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cfg)
	service := NewCachedMetadataService(c)

	ctx := context.Background()
	tenantID := "tenant-1"
	datasourceID := "ds-1"
	dsn := "invalid-dsn"

	tables, err := service.GetTables(ctx, tenantID, datasourceID, dsn)

	assert.Error(t, err)
	assert.Nil(t, tables)
}

func TestCachedMetadataService_GetFields_InvalidDSN(t *testing.T) {
	cfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cfg)
	service := NewCachedMetadataService(c)

	ctx := context.Background()
	tenantID := "tenant-1"
	datasourceID := "ds-1"
	dsn := "invalid-dsn"

	fields, err := service.GetFields(ctx, tenantID, datasourceID, dsn, "nonexistent_table")

	assert.Error(t, err)
	assert.Nil(t, fields)
}

func TestCachedMetadataService_WithCache(t *testing.T) {
	dsn := getTestDSNForCachedMetadata()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	cfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cfg)
	service := NewCachedMetadataService(c)

	ctx := context.Background()
	tenantID := "tenant-cache-test"
	datasourceID := "ds-cache-test"

	tables, err := service.GetTables(ctx, tenantID, datasourceID, dsn)
	require.NoError(t, err)
	assert.NotNil(t, tables)

	fields, err := service.GetFields(ctx, tenantID, datasourceID, dsn, "tenants")
	require.NoError(t, err)
	assert.NotNil(t, fields)
}

func TestFieldInfo_JSON(t *testing.T) {
	field := FieldInfo{
		Name:     "username",
		Type:     "varchar(255)",
		Nullable: true,
		Comment:  "User name field",
	}

	data, err := json.Marshal(field)
	require.NoError(t, err)

	var unmarshaled FieldInfo
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, field.Name, unmarshaled.Name)
	assert.Equal(t, field.Type, unmarshaled.Type)
	assert.Equal(t, field.Nullable, unmarshaled.Nullable)
	assert.Equal(t, field.Comment, unmarshaled.Comment)
}

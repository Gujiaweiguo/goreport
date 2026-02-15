package render

import (
	"context"
	"os"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDSN() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func skipIfNoDB(t *testing.T) {
	if getTestDSN() == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
}

func TestEngine_Render_WithDatabase(t *testing.T) {
	skipIfNoDB(t)

	dsn := getTestDSN()
	cfg := &config.DatabaseConfig{DSN: dsn}
	db, err := database.InitWithConfig(dsn, cfg)
	require.NoError(t, err)
	require.NotNil(t, db)

	cacheCfg := config.CacheConfig{Enabled: false}
	c, _ := cache.New(cacheCfg)

	engine := NewEngine(db, c)
	ctx := context.Background()

	t.Run("render simple config", func(t *testing.T) {
		configJSON := `{
			"rows": 3,
			"cols": 3,
			"cells": [
				{"row": 0, "col": 0, "text": "Hello"},
				{"row": 1, "col": 1, "text": "World"}
			]
		}`

		result, err := engine.Render(ctx, configJSON, nil, "test-tenant")
		assert.NoError(t, err)
		assert.Contains(t, result, "Hello")
		assert.Contains(t, result, "World")
	})

	t.Run("render with datasource but no match", func(t *testing.T) {
		configJSON := `{
			"rows": 2,
			"cols": 2,
			"cells": [
				{
					"row": 0, 
					"col": 0, 
					"datasourceId": "non-existent-ds",
					"tableName": "users",
					"fieldName": "name"
				}
			]
		}`

		result, err := engine.Render(ctx, configJSON, nil, "test-tenant")
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("render with invalid JSON", func(t *testing.T) {
		configJSON := `{invalid json}`

		_, err := engine.Render(ctx, configJSON, nil, "test-tenant")
		assert.Error(t, err)
	})

	t.Run("render with params", func(t *testing.T) {
		configJSON := `{
			"rows": 1,
			"cols": 1,
			"cells": [{"row": 0, "col": 0, "text": "Test"}]
		}`

		params := map[string]interface{}{
			"page":     float64(1),
			"pageSize": float64(10),
		}

		result, err := engine.Render(ctx, configJSON, params, "test-tenant")
		assert.NoError(t, err)
		assert.Contains(t, result, "Test")
	})

	t.Run("render with empty cells", func(t *testing.T) {
		configJSON := `{
			"rows": 2,
			"cols": 2,
			"cells": []
		}`

		result, err := engine.Render(ctx, configJSON, nil, "test-tenant")
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("render with cell value field", func(t *testing.T) {
		configJSON := `{
			"rows": 1,
			"cols": 1,
			"cells": [{"row": 0, "col": 0, "value": "CellValue"}]
		}`

		result, err := engine.Render(ctx, configJSON, nil, "test-tenant")
		assert.NoError(t, err)
		assert.Contains(t, result, "CellValue")
	})
}

func TestEngine_WithNilDependencies(t *testing.T) {
	t.Run("with nil db and cache", func(t *testing.T) {
		engine := NewEngine(nil, nil)
		assert.NotNil(t, engine)
		assert.Nil(t, engine.db)
		assert.Nil(t, engine.cache)
	})
}

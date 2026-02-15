package render

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTestDSNForRender() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func skipIfNoDBForRender(t *testing.T) {
	if getTestDSNForRender() == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
}

func setupRenderIntegrationTest(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := getTestDSNForRender()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db
}

func TestEngine_FetchCellValueFromDB_Success(t *testing.T) {
	skipIfNoDBForRender(t)

	db := setupRenderIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	// Create test datasource pointing to goreport database
	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  "test-render",
		Name:      "Render Test DataSource",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  "goreport",
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	engine := NewEngine(db, nil)
	// Query the database field from data_sources (should always have data)
	tableName := "data_sources"
	fieldName := "database"

	cell := Cell{
		DatasourceID: &datasourceID,
		TableName:    &tableName,
		FieldName:    &fieldName,
	}

	value, err := engine.fetchCellValueFromDB(ctx, cell, *ds)

	assert.NoError(t, err)
	// Value should be non-empty (returns byte array representation)
	assert.NotEmpty(t, value)
}

func TestEngine_FetchCellValueFromDB_EmptyResult(t *testing.T) {
	skipIfNoDBForRender(t)

	db := setupRenderIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	// Create test datasource
	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  "test-render-empty",
		Name:      "Render Empty Test",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  "goreport",
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	engine := NewEngine(db, nil)

	// Query a non-existent table
	tableName := "non_existent_table_xyz"
	fieldName := "id"

	cell := Cell{
		DatasourceID: &datasourceID,
		TableName:    &tableName,
		FieldName:    &fieldName,
	}

	// Should return error for non-existent table
	value, err := engine.fetchCellValueFromDB(ctx, cell, *ds)

	assert.Error(t, err)
	assert.Empty(t, value)
}

func TestEngine_FetchCellValueFromDB_NullField(t *testing.T) {
	skipIfNoDBForRender(t)

	db := setupRenderIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	// Create test datasource
	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  "test-render-null",
		Name:      "Render Null Test",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  "goreport",
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	engine := NewEngine(db, nil)

	// Query deleted_at which is typically NULL for most records
	tableName := "data_sources"
	fieldName := "deleted_at"

	cell := Cell{
		DatasourceID: &datasourceID,
		TableName:    &tableName,
		FieldName:    &fieldName,
	}

	value, err := engine.fetchCellValueFromDB(ctx, cell, *ds)

	// NULL values should return empty string without error
	assert.NoError(t, err)
	assert.Empty(t, value)
}

func TestEngine_FetchCellValue_Integration(t *testing.T) {
	skipIfNoDBForRender(t)

	db := setupRenderIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)
	tenantID := "test-render-full"

	// Create test datasource
	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  tenantID,
		Name:      "Full Integration Test",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  "goreport",
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	engine := NewEngine(db, nil)

	tableName := "data_sources"
	fieldName := "type"

	cell := Cell{
		DatasourceID: &datasourceID,
		TableName:    &tableName,
		FieldName:    &fieldName,
	}

	value, err := engine.fetchCellValue(ctx, cell, tenantID)

	assert.NoError(t, err)
	// Value should be non-empty
	assert.NotEmpty(t, value)
}

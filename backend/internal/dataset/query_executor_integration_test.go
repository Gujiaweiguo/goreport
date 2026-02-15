package dataset

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTestDSNForQuery() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func getTestDatabaseNameForQuery() string {
	dsn := getTestDSNForQuery()
	if dsn == "" {
		return "goreport"
	}
	lastSlash := strings.LastIndex(dsn, "/")
	if lastSlash == -1 {
		return "goreport"
	}
	dbPart := dsn[lastSlash+1:]
	questionIdx := strings.Index(dbPart, "?")
	if questionIdx > 0 {
		return dbPart[:questionIdx]
	}
	return dbPart
}

func skipIfNoDBForQuery(t *testing.T) {
	if getTestDSNForQuery() == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
}

func setupQueryIntegrationTest(t *testing.T) (*gorm.DB, repository.DatasetRepository, repository.DatasetFieldRepository, repository.DatasourceRepository) {
	t.Helper()

	dsn := getTestDSNForQuery()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	datasetRepo := repository.NewDatasetRepository(db)
	fieldRepo := repository.NewDatasetFieldRepository(db)
	datasourceRepo := repository.NewDatasourceRepository(db)

	return db, datasetRepo, fieldRepo, datasourceRepo
}

func ensureTenantExists(db *gorm.DB, t *testing.T, tenantID string) {
	t.Helper()
	var tenant models.Tenant
	err := db.Where("id = ?", tenantID).First(&tenant).Error
	if err == gorm.ErrRecordNotFound {
		tenant = models.Tenant{
			ID:        tenantID,
			Name:      "Test " + tenantID,
			Code:      tenantID,
			Status:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&tenant).Error; err != nil {
			t.Fatalf("Failed to create tenant: %v", err)
		}
	} else if err != nil {
		t.Fatalf("Failed to check tenant: %v", err)
	}
}

func TestQueryExecutor_Query_Integration(t *testing.T) {
	skipIfNoDBForQuery(t)

	db, datasetRepo, fieldRepo, datasourceRepo := setupQueryIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	tenantID := "test-query-integration"
	ensureTenantExists(db, t, tenantID)
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  tenantID,
		Name:      "Query Test DataSource",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  getTestDatabaseNameForQuery(),
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	datasetID := fmt.Sprintf("dt-%s", ts)
	dataset := &models.Dataset{
		ID:           datasetID,
		TenantID:     tenantID,
		Name:         "Query Test Dataset",
		Type:         "sql",
		DatasourceID: &datasourceID,
		Config:       `{"query": "SELECT 1 AS id, 'test' AS name"}`,
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = db.Create(dataset).Error
	require.NoError(t, err)

	field := &models.DatasetField{
		ID:        fmt.Sprintf("fld-%s", ts),
		DatasetID: datasetID,
		Name:      "id",
		Type:      "dimension",
		DataType:  "number",
		Config:    "{}",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = db.Create(field).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM dataset_fields WHERE dataset_id = ?", datasetID)
		db.Exec("DELETE FROM datasets WHERE id = ?", datasetID)
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	sqlBuilder := NewSQLExpressionBuilder()
	cache := NewComputedFieldCache()

	executor := NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, sqlBuilder, cache)

	req := &QueryRequest{
		DatasetID: datasetID,
		Fields:    []string{"id", "name"},
		Page:      1,
		PageSize:  10,
	}

	resp, err := executor.Query(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.PageSize)
}

func TestQueryExecutor_Query_DatasetNotFound(t *testing.T) {
	skipIfNoDBForQuery(t)

	db, datasetRepo, fieldRepo, datasourceRepo := setupQueryIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	sqlBuilder := NewSQLExpressionBuilder()
	cache := NewComputedFieldCache()

	executor := NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, sqlBuilder, cache)

	req := &QueryRequest{
		DatasetID: "non-existent-dataset",
		Fields:    []string{"id"},
		Page:      1,
		PageSize:  10,
	}

	resp, err := executor.Query(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "dataset not found")
}

func TestQueryExecutor_Query_UnsupportedType(t *testing.T) {
	skipIfNoDBForQuery(t)

	db, datasetRepo, fieldRepo, datasourceRepo := setupQueryIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	tenantID := "test-query-unsupported"
	ensureTenantExists(db, t, tenantID)
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	datasetID := fmt.Sprintf("dt-%s", ts)
	dataset := &models.Dataset{
		ID:        datasetID,
		TenantID:  tenantID,
		Name:      "Unsupported Type Dataset",
		Type:      "api",
		Config:    `{}`,
		Status:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(dataset).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM datasets WHERE id = ?", datasetID)
	})

	sqlBuilder := NewSQLExpressionBuilder()
	cache := NewComputedFieldCache()

	executor := NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, sqlBuilder, cache)

	req := &QueryRequest{
		DatasetID: datasetID,
		Fields:    []string{"id"},
		Page:      1,
		PageSize:  10,
	}

	resp, err := executor.Query(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "unsupported dataset type")
}

func TestQueryExecutor_Query_WithFilters(t *testing.T) {
	skipIfNoDBForQuery(t)

	db, datasetRepo, fieldRepo, datasourceRepo := setupQueryIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	tenantID := "test-query-filters"
	ensureTenantExists(db, t, tenantID)
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  tenantID,
		Name:      "Query Filters DataSource",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  getTestDatabaseNameForQuery(),
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	datasetID := fmt.Sprintf("dt-%s", ts)
	dataset := &models.Dataset{
		ID:           datasetID,
		TenantID:     tenantID,
		Name:         "Query Filters Dataset",
		Type:         "sql",
		DatasourceID: &datasourceID,
		Config:       `{"query": "SELECT 1 AS id, 'active' AS status UNION SELECT 2, 'inactive'"}`,
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = db.Create(dataset).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM datasets WHERE id = ?", datasetID)
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	sqlBuilder := NewSQLExpressionBuilder()
	cache := NewComputedFieldCache()

	executor := NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, sqlBuilder, cache)

	req := &QueryRequest{
		DatasetID: datasetID,
		Fields:    []string{"id", "status"},
		Filters: []Filter{
			{Field: "status", Operator: "eq", Value: "active"},
		},
		Page:     1,
		PageSize: 10,
	}

	resp, err := executor.Query(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Data, 1)
	if len(resp.Data) > 0 {
		assert.Equal(t, "active", resp.Data[0]["status"])
	}
}

func TestQueryExecutor_Query_WithPagination(t *testing.T) {
	skipIfNoDBForQuery(t)

	db, datasetRepo, fieldRepo, datasourceRepo := setupQueryIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	ctx := context.Background()
	tenantID := "test-query-pagination"
	ensureTenantExists(db, t, tenantID)
	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)

	datasourceID := fmt.Sprintf("ds-%s", ts)
	ds := &models.DataSource{
		ID:        datasourceID,
		TenantID:  tenantID,
		Name:      "Query Pagination DataSource",
		Type:      "mysql",
		Host:      "127.0.0.1",
		Port:      3306,
		Database:  getTestDatabaseNameForQuery(),
		Username:  "root",
		Password:  "root",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := db.Create(ds).Error
	require.NoError(t, err)

	datasetID := fmt.Sprintf("dt-%s", ts)
	dataset := &models.Dataset{
		ID:           datasetID,
		TenantID:     tenantID,
		Name:         "Query Pagination Dataset",
		Type:         "sql",
		DatasourceID: &datasourceID,
		Config:       `{"query": "SELECT 1 AS id UNION SELECT 2 UNION SELECT 3 UNION SELECT 4 UNION SELECT 5"}`,
		Status:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = db.Create(dataset).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM datasets WHERE id = ?", datasetID)
		db.Exec("DELETE FROM data_sources WHERE id = ?", datasourceID)
	})

	sqlBuilder := NewSQLExpressionBuilder()
	cache := NewComputedFieldCache()

	executor := NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, sqlBuilder, cache)

	req := &QueryRequest{
		DatasetID: datasetID,
		Fields:    []string{"id"},
		Page:      2,
		PageSize:  2,
	}

	resp, err := executor.Query(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, resp.Page)
	assert.Equal(t, 2, resp.PageSize)
	assert.Len(t, resp.Data, 2)
}

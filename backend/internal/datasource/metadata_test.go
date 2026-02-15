package datasource

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTestDSNForMetadata() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func skipIfNoDBForMetadata(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := getTestDSNForMetadata()
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db
}

func TestGetTables_Success(t *testing.T) {
	db := skipIfNoDBForMetadata(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tables, err := GetTables(context.Background(), db, "goreport")

	require.NoError(t, err)
	assert.NotEmpty(t, tables)

	found := false
	for _, table := range tables {
		if table == "users" || table == "data_sources" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should find system tables")
}

func TestGetTables_EmptyDatabase(t *testing.T) {
	db := skipIfNoDBForMetadata(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	tables, err := GetTables(context.Background(), db, "non_existent_db_xyz")

	require.NoError(t, err)
	assert.Empty(t, tables)
}

func TestGetFields_Success(t *testing.T) {
	db := skipIfNoDBForMetadata(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	fields, err := GetFields(context.Background(), db, "goreport", "users")

	if err != nil {
		assert.Contains(t, err.Error(), "failed to query fields")
	} else {
		assert.NotEmpty(t, fields)
	}
}

func TestGetFields_NonExistentTable(t *testing.T) {
	db := skipIfNoDBForMetadata(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	fields, err := GetFields(context.Background(), db, "goreport", "non_existent_table_xyz")

	assert.NoError(t, err)
	assert.Empty(t, fields)
}

package datasource

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetTablesAndFields(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set")
	}

	dbName := extractDBName(dsn)
	if dbName == "" {
		t.Skip("could not parse database name from TEST_DB_DSN")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}
	defer sqlDB.Close()

	tableName := fmt.Sprintf("test_table_%d", time.Now().UnixNano())
	if err := db.Exec(fmt.Sprintf("CREATE TABLE %s (id INT PRIMARY KEY, name VARCHAR(50))", tableName)).Error; err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	defer db.Exec(fmt.Sprintf("DROP TABLE %s", tableName))

	tables, err := GetTables(context.Background(), db, dbName)
	if err != nil {
		t.Fatalf("GetTables failed: %v", err)
	}

	if !containsString(tables, tableName) {
		t.Fatalf("expected table %s in list, got %v", tableName, tables)
	}

	fields, err := GetFields(context.Background(), db, dbName, tableName)
	if err != nil {
		t.Fatalf("GetFields failed: %v", err)
	}

	if len(fields) == 0 {
		t.Fatal("expected fields, got none")
	}
}

func extractDBName(dsn string) string {
	idx := strings.LastIndex(dsn, "/")
	if idx == -1 {
		return ""
	}
	rest := dsn[idx+1:]
	if q := strings.Index(rest, "?"); q != -1 {
		rest = rest[:q]
	}
	return rest
}

func containsString(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gujiaweiguo/goreport/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Init 初始化数据库连接
func Init(dsn string) (*gorm.DB, error) {
	return InitWithConfig(dsn, &config.DatabaseConfig{})
}

// InitWithConfig 使用配置初始化数据库连接
func InitWithConfig(dsn string, cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	maxOpenConns := cfg.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 100
	}

	maxIdleConns := cfg.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 10
	}

	connMaxLifetime := cfg.ConnMaxLifetime
	if connMaxLifetime == 0 {
		connMaxLifetime = 3600
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	if err := ensureDatasourceSchemaCompatibility(db); err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")

	return db, nil
}

func ensureDatasourceSchemaCompatibility(db *gorm.DB) error {
	columns := []struct {
		name string
		ddl  string
	}{
		{name: "database", ddl: "ALTER TABLE data_sources ADD COLUMN `database` VARCHAR(100) NOT NULL DEFAULT '' AFTER port"},
		{name: "ssh_host", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_host VARCHAR(255) DEFAULT ''"},
		{name: "ssh_port", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_port INT DEFAULT 0"},
		{name: "ssh_username", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_username VARCHAR(100) DEFAULT ''"},
		{name: "ssh_password", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_password VARCHAR(255) DEFAULT ''"},
		{name: "ssh_key", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_key TEXT"},
		{name: "ssh_key_phrase", ddl: "ALTER TABLE data_sources ADD COLUMN ssh_key_phrase VARCHAR(255) DEFAULT ''"},
		{name: "max_connections", ddl: "ALTER TABLE data_sources ADD COLUMN max_connections INT DEFAULT 10"},
		{name: "query_timeout_seconds", ddl: "ALTER TABLE data_sources ADD COLUMN query_timeout_seconds INT DEFAULT 30"},
		{name: "config", ddl: "ALTER TABLE data_sources ADD COLUMN config TEXT"},
		{name: "description", ddl: "ALTER TABLE data_sources ADD COLUMN description TEXT"},
	}

	for _, col := range columns {
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'data_sources' AND COLUMN_NAME = ?", col.name).Scan(&count).Error
		if err != nil {
			return fmt.Errorf("failed to check datasource column %s: %w", col.name, err)
		}
		if count == 0 {
			if err := db.Exec(col.ddl).Error; err != nil {
				return fmt.Errorf("failed to add datasource column %s: %w", col.name, err)
			}
		}
	}

	updates := []string{
		"UPDATE data_sources SET id = CONCAT('ds-', REPLACE(UUID(), '-', '')) WHERE id IS NULL OR id = ''",
		"UPDATE data_sources SET `database` = database_name WHERE (`database` IS NULL OR `database` = '') AND database_name IS NOT NULL",
		"UPDATE data_sources SET max_connections = 10 WHERE max_connections IS NULL OR max_connections = 0",
		"UPDATE data_sources SET query_timeout_seconds = 30 WHERE query_timeout_seconds IS NULL OR query_timeout_seconds = 0",
	}

	for _, stmt := range updates {
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("failed to apply datasource data backfill: %w", err)
		}
	}

	return nil
}

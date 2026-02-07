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

	log.Println("Database connected successfully")

	return db, nil
}

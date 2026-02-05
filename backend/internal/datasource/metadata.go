package datasource

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type FieldInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Comment  string `json:"comment"`
}

func GetTables(ctx context.Context, db *gorm.DB, database string) ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = ? AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	var tables []string
	if err := db.WithContext(ctx).Raw(query, database).Scan(&tables).Error; err != nil {
		return nil, err
	}

	return tables, nil
}

func GetFields(ctx context.Context, db *gorm.DB, database, tableName string) ([]FieldInfo, error) {
	query := `
		SELECT column_name, data_type, is_nullable, column_comment
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ?
		ORDER BY ordinal_position
	`

	var fields []FieldInfo
	if err := db.WithContext(ctx).Raw(query, database, tableName).Scan(&fields).Error; err != nil {
		return nil, fmt.Errorf("failed to query fields: %w", err)
	}

	return fields, nil
}

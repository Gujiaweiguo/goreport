package render

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (e *Engine) fetchCellValue(ctx context.Context, cell Cell, tenantID string) (string, error) {
	var ds models.DataSource
	if err := e.db.WithContext(ctx).Where("id = ?", *cell.DatasourceID).First(&ds).Error; err != nil {
		return "", err
	}

	if e.cache == nil {
		return e.fetchCellValueFromDB(ctx, cell, ds)
	}

	domain := "report:data"
	identity := *cell.DatasourceID

	params := map[string]interface{}{
		"table": *cell.TableName,
		"field": *cell.FieldName,
		"limit": 1,
	}

	if cached, hit, err := e.cache.Get(ctx, tenantID, domain, identity, params); err == nil && hit {
		var result string
		if err := json.Unmarshal(cached, &result); err == nil {
			return result, nil
		}
	}

	value, err := e.fetchCellValueFromDB(ctx, cell, ds)
	if err != nil {
		return "", err
	}

	data, _ := json.Marshal(value)
	_ = e.cache.Set(ctx, tenantID, domain, identity, params, data)

	return value, nil
}

func (e *Engine) fetchCellValueFromDB(ctx context.Context, cell Cell, ds models.DataSource) (string, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		ds.Port,
		ds.DatabaseName,
	)

	dataDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return "", err
	}

	sqlDB, err := dataDB.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	query := fmt.Sprintf("SELECT `%s` FROM `%s` LIMIT 1", *cell.FieldName, *cell.TableName)
	rows, err := dataDB.WithContext(ctx).Raw(query).Rows()
	if err != nil {
		return "", err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}
	if len(cols) == 0 {
		return "", nil
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	if rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return "", err
		}
		if values[0] == nil {
			return "", nil
		}
		return fmt.Sprintf("%v", values[0]), nil
	}

	return "", nil
}

package render

import (
	"context"
	"fmt"

	"github.com/jeecg/jimureport-go/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (e *Engine) fetchCellValue(ctx context.Context, cell Cell) (string, error) {
	var ds models.DataSource
	if err := e.db.WithContext(ctx).Where("id = ?", *cell.DatasourceID).First(&ds).Error; err != nil {
		return "", err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		ds.Port,
		ds.Database,
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

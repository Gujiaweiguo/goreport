package render

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
)

type Engine struct {
	db *gorm.DB
}

func NewEngine(db *gorm.DB) *Engine {
	return &Engine{db: db}
}

func (e *Engine) Render(ctx context.Context, configJSON string, params map[string]interface{}) (string, error) {
	var config ReportConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return "", err
	}

	cellValues := make(map[string]string)
	for _, cell := range config.Cells {
		value := cell.Value
		if value == "" {
			value = cell.Text
		}
		if cell.DatasourceID != nil && cell.TableName != nil && cell.FieldName != nil {
			result, err := e.fetchCellValue(ctx, cell)
			if err == nil && result != "" {
				value = result
			}
		}
		cellValues[cellKey(cell.Row, cell.Col)] = value
	}

	return buildHTML(&config, cellValues), nil
}

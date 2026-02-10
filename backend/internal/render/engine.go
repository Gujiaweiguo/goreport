package render

import (
	"context"
	"encoding/json"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"gorm.io/gorm"
)

type Engine struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewEngine(db *gorm.DB, cache *cache.Cache) *Engine {
	return &Engine{
		db:    db,
		cache: cache,
	}
}

func (e *Engine) Render(ctx context.Context, configJSON string, params map[string]interface{}, tenantID string) (string, error) {
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
			result, err := e.fetchCellValue(ctx, cell, tenantID)
			if err == nil && result != "" {
				value = result
			}
		}
		cellValues[cellKey(cell.Row, cell.Col)] = value
	}

	page := 0
	pageSize := 0
	if params != nil {
		if pageVal, ok := params["page"].(float64); ok {
			page = int(pageVal)
		}
		if pageSizeVal, ok := params["pageSize"].(float64); ok {
			pageSize = int(pageSizeVal)
		}
	}

	return buildHTML(&config, cellValues, page, pageSize), nil
}

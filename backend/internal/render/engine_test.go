package render

import (
	"context"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewEngine(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	cfg := config.CacheConfig{Enabled: false}
	cacheObj, _ := cache.New(cfg)

	engine := NewEngine(db, cacheObj)
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.db)
	assert.NotNil(t, engine.cache)
}

func TestNewEngine_NilCache(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})

	engine := NewEngine(db, nil)
	assert.NotNil(t, engine)
	assert.NotNil(t, engine.db)
	assert.Nil(t, engine.cache)
}

func TestEngine_Render_EmptyConfig(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	html, err := engine.Render(context.Background(), `{"cells":[]}`, nil, "tenant-1")

	assert.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.Contains(t, html, "<table>")
	assert.Contains(t, html, "</table>")
}

func TestEngine_Render_InvalidJSON(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	html, err := engine.Render(context.Background(), `{invalid json`, nil, "tenant-1")

	assert.Error(t, err)
	assert.Empty(t, html)
}

func TestEngine_Render_StaticCells(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Header 1"},
			{"row": 0, "col": 1, "text": "Header 2"},
			{"row": 1, "col": 0, "text": "Data 1"},
			{"row": 1, "col": 1, "text": "Data 2"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Header 1")
	assert.Contains(t, html, "Header 2")
	assert.Contains(t, html, "Data 1")
	assert.Contains(t, html, "Data 2")
	assert.Contains(t, html, "<table>")
	assert.Contains(t, html, "<tr>")
	assert.Contains(t, html, "<td>")
}

func TestEngine_Render_CellWithValue(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "value": "Value1", "text": "Text1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Value1")
	assert.NotContains(t, html, "Text1")
}

func TestEngine_Render_CellWithTextNoValue(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "value": "", "text": "Text1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Text1")
}

func TestEngine_Render_WithPagination(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"},
			{"row": 1, "col": 0, "text": "Row 1"},
			{"row": 2, "col": 0, "text": "Row 2"},
			{"row": 3, "col": 0, "text": "Row 3"},
			{"row": 4, "col": 0, "text": "Row 4"},
			{"row": 5, "col": 0, "text": "Row 5"},
			{"row": 6, "col": 0, "text": "Row 6"},
			{"row": 7, "col": 0, "text": "Row 7"},
			{"row": 8, "col": 0, "text": "Row 8"},
			{"row": 9, "col": 0, "text": "Row 9"},
			{"row": 10, "col": 0, "text": "Row 10"},
			{"row": 11, "col": 0, "text": "Row 11"}
		]
	}`

	html, err := engine.Render(context.Background(), config, map[string]interface{}{
		"page":     float64(1),
		"pageSize": float64(10),
	}, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
	assert.Contains(t, html, "Row 9")
	assert.NotContains(t, html, "Row 10")
	assert.NotContains(t, html, "Row 11")
}

func TestEngine_Render_WithInvalidPagination(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"},
			{"row": 1, "col": 0, "text": "Row 1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, map[string]interface{}{
		"page":     "invalid",
		"pageSize": "invalid",
	}, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
	assert.Contains(t, html, "Row 1")
}

func TestEngine_Render_CellWithDatasource_NoDatabase(t *testing.T) {
	t.Skip("Requires actual database connection")
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	dsID := "ds-1"
	tableName := "users"
	fieldName := "name"

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Static", "datasourceId": "` + dsID + `", "tableName": "` + tableName + `", "fieldName": "` + fieldName + `"}
		]
	}`

	_, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.Error(t, err)
}

func TestEngine_Render_CellWithMissingDatasourceFields(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	dsID := "ds-1"
	tableName := "users"

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Static", "datasourceId": "` + dsID + `", "tableName": "` + tableName + `"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Static")
}

func TestEngine_Render_ParamsInt(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"},
			{"row": 1, "col": 0, "text": "Row 1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, map[string]interface{}{
		"page":     int(1),
		"pageSize": int(10),
	}, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
}

func TestEngine_Render_NilParams(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
}

func TestEngine_Render_MixedCells(t *testing.T) {
	t.Skip("Requires actual database connection")
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Static 1"},
			{"row": 0, "col": 1, "text": "Static 2", "datasourceId": "ds-1", "tableName": "users", "fieldName": "name"},
			{"row": 1, "col": 0, "value": "Value 1", "text": "Text 1"},
			{"row": 1, "col": 1, "value": "Value 2", "text": "Text 2", "datasourceId": "ds-1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Static 1")
	assert.Contains(t, html, "Static 2")
	assert.Contains(t, html, "Value 1")
	assert.Contains(t, html, "Value 2")
}

func TestEngine_Render_LargeDataset(t *testing.T) {
	t.Skip("Requires actual database connection")
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	var configJSON string
	configJSON = `{"cells": [`
	for i := 0; i < 100; i++ {
		if i > 0 {
			configJSON += ","
		}
		configJSON += `{"row":` + string(rune('0'+i%10)) + `,"col":0,"text":"Data ` + string(rune('0'+i%10)) + `"}`
	}
	configJSON += `]}`

	html, err := engine.Render(context.Background(), configJSON, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Data 0")
}

func TestEngine_Render_EmptyRowCol(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "A"},
			{"row": 0, "col": 1, "text": ""}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "A")
}

func TestEngine_Render_ContextCancellation(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	config := `{"cells": [{"row": 0, "col": 0, "text": "Test"}]}`

	html, err := engine.Render(ctx, config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Test")
}

func TestCell_Fields(t *testing.T) {
	dsID := "ds-123"
	tableName := "users"
	fieldName := "name"

	cell := Cell{
		Row:          1,
		Col:          2,
		Value:        "value",
		Text:         "text",
		DatasourceID: &dsID,
		TableName:    &tableName,
		FieldName:    &fieldName,
	}

	assert.Equal(t, 1, cell.Row)
	assert.Equal(t, 2, cell.Col)
	assert.Equal(t, "value", cell.Value)
	assert.Equal(t, "text", cell.Text)
	assert.Equal(t, &dsID, cell.DatasourceID)
	assert.Equal(t, &tableName, cell.TableName)
	assert.Equal(t, &fieldName, cell.FieldName)
}

func TestReportConfig_Fields(t *testing.T) {
	config := ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "A"},
			{Row: 1, Col: 1, Text: "B"},
		},
	}

	assert.Len(t, config.Cells, 2)
	assert.Equal(t, "A", config.Cells[0].Text)
	assert.Equal(t, "B", config.Cells[1].Text)
}

func TestEngine_Render_WithDatasourceCell(t *testing.T) {
	t.Skip("Requires actual database connection")
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	dsID := "ds-1"
	tableName := "users"
	fieldName := "name"

	// Cell with datasource binding but DB lookup will fail
	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Fallback", "datasourceId": "` + dsID + `", "tableName": "` + tableName + `", "fieldName": "` + fieldName + `"}
		]
	}`

	// Should use fallback text when datasource lookup fails
	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Fallback")
}

func TestEngine_Render_WithAllOptionalFields(t *testing.T) {
	t.Skip("Requires actual database connection")
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	dsID := "ds-1"
	tableName := "users"
	fieldName := "name"

	// Cell with all optional fields
	config := `{
		"cells": [
			{
				"row": 0,
				"col": 0,
				"value": "CellValue",
				"text": "CellText",
				"datasourceId": "` + dsID + `",
				"tableName": "` + tableName + `",
				"fieldName": "` + fieldName + `"
			}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "CellValue")
}

func TestEngine_Render_WithNegativePage(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"},
			{"row": 1, "col": 0, "text": "Row 1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, map[string]interface{}{
		"page":     float64(-1),
		"pageSize": float64(10),
	}, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
}

func TestEngine_Render_WithZeroPageSize(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Row 0"},
			{"row": 1, "col": 0, "text": "Row 1"}
		]
	}`

	html, err := engine.Render(context.Background(), config, map[string]interface{}{
		"page":     float64(0),
		"pageSize": float64(0),
	}, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Row 0")
	assert.Contains(t, html, "Row 1")
}

func TestEngine_Render_WithCacheNil(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	// Verify cache is nil
	assert.Nil(t, engine.cache)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Test"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Test")
}

func TestEngine_Render_WithCacheEnabled(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	cfg := config.CacheConfig{Enabled: false}
	cacheObj, _ := cache.New(cfg)
	engine := NewEngine(db, cacheObj)

	// Verify cache is not nil (but degraded since disabled)
	assert.NotNil(t, engine.cache)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "Test"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Test")
}

func TestEngine_Render_LargeRowCol(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 1000, "col": 500, "text": "Far Cell"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "Far Cell")
}

func TestEngine_Render_OnlyTextNoValue(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "text": "OnlyText"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "OnlyText")
}

func TestEngine_Render_OnlyValueNoText(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "value": "OnlyValue"}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "OnlyValue")
}

func TestEngine_Render_EmptyValueAndText(t *testing.T) {
	db, _ := gorm.Open(mysql.Open("user:pass@tcp(localhost:3306)/db"), &gorm.Config{})
	engine := NewEngine(db, nil)

	config := `{
		"cells": [
			{"row": 0, "col": 0, "value": "", "text": ""}
		]
	}`

	html, err := engine.Render(context.Background(), config, nil, "tenant-1")

	assert.NoError(t, err)
	assert.Contains(t, html, "<td></td>")
}

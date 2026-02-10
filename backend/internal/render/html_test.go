package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellKey(t *testing.T) {
	assert.Equal(t, "0:0", cellKey(0, 0))
	assert.Equal(t, "5:10", cellKey(5, 10))
	assert.Equal(t, "100:200", cellKey(100, 200))
}

func TestBuildHTML_NoPagination(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 99, Col: 0, Text: "Row 99"},
		},
	}

	cellValues := map[string]string{
		"0:0":  "Row 0",
		"99:0": "Row 99",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "Row 0")
	assert.Contains(t, result, "Row 99")
}

func TestBuildHTML_PaginationPage1(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 49, Col: 0, Text: "Row 49"},
			{Row: 50, Col: 0, Text: "Row 50"},
			{Row: 99, Col: 0, Text: "Row 99"},
		},
	}

	cellValues := map[string]string{
		"0:0":  "Row 0",
		"49:0": "Row 49",
		"50:0": "Row 50",
		"99:0": "Row 99",
	}

	result := buildHTML(config, cellValues, 1, 50)
	assert.Contains(t, result, "Row 0")
	assert.Contains(t, result, "Row 49")
	assert.NotContains(t, result, "Row 50")
	assert.NotContains(t, result, "Row 99")
}

func TestBuildHTML_PaginationPage2(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 49, Col: 0, Text: "Row 49"},
			{Row: 50, Col: 0, Text: "Row 50"},
			{Row: 99, Col: 0, Text: "Row 99"},
		},
	}

	cellValues := map[string]string{
		"0:0":  "Row 0",
		"49:0": "Row 49",
		"50:0": "Row 50",
		"99:0": "Row 99",
	}

	result := buildHTML(config, cellValues, 2, 50)
	assert.NotContains(t, result, "Row 0")
	assert.NotContains(t, result, "Row 49")
	assert.Contains(t, result, "Row 50")
	assert.Contains(t, result, "Row 99")
}

func TestGetTotalRows(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "A1"},
			{Row: 99, Col: 0, Text: "A100"},
		},
	}

	totalRows := GetTotalRows(config)
	assert.Equal(t, 100, totalRows)
}

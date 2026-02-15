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

func TestBuildHTMLEmptyConfig(t *testing.T) {
	config := &ReportConfig{Cells: []Cell{}}
	cellValues := map[string]string{}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "<table>")
	assert.Contains(t, result, "</table>")
}

func TestBuildHTML_EscapeHTML(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Test"},
		},
	}

	cellValues := map[string]string{
		"0:0": "<script>alert('xss')</script>",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "&lt;script&gt;")
	assert.NotContains(t, result, "<script>")
}

func TestBuildHTML_SpecialCharacters(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Test"},
		},
	}

	cellValues := map[string]string{
		"0:0": "A & B < C > D \"E\" 'F'",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "&amp;")
	assert.Contains(t, result, "&lt;")
	assert.Contains(t, result, "&gt;")
}

func TestBuildHTML_MultipleColumns(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "A"},
			{Row: 0, Col: 1, Text: "B"},
			{Row: 0, Col: 2, Text: "C"},
		},
	}

	cellValues := map[string]string{
		"0:0": "A",
		"0:1": "B",
		"0:2": "C",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "A")
	assert.Contains(t, result, "B")
	assert.Contains(t, result, "C")
}

func TestBuildHTML_PaginationBeyondData(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 1, Col: 0, Text: "Row 1"},
		},
	}

	cellValues := map[string]string{
		"0:0": "Row 0",
		"1:0": "Row 1",
	}

	result := buildHTML(config, cellValues, 10, 10)
	assert.NotContains(t, result, "Row 0")
	assert.NotContains(t, result, "Row 1")
}

func TestBuildHTML_PageZero(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 100, Col: 0, Text: "Row 100"},
		},
	}

	cellValues := map[string]string{
		"0:0":   "Row 0",
		"100:0": "Row 100",
	}

	result := buildHTML(config, cellValues, 0, 10)
	assert.Contains(t, result, "Row 0")
	assert.Contains(t, result, "Row 100")
}

func TestBuildHTML_PageSizeZero(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Row 0"},
			{Row: 100, Col: 0, Text: "Row 100"},
		},
	}

	cellValues := map[string]string{
		"0:0":   "Row 0",
		"100:0": "Row 100",
	}

	result := buildHTML(config, cellValues, 1, 0)
	assert.Contains(t, result, "Row 0")
	assert.Contains(t, result, "Row 100")
}

func TestGetTotalRows_EmptyConfig(t *testing.T) {
	config := &ReportConfig{Cells: []Cell{}}

	totalRows := GetTotalRows(config)
	assert.Equal(t, 1, totalRows)
}

func TestGetTotalRows_SingleCell(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 5, Col: 3, Text: "Test"},
		},
	}

	totalRows := GetTotalRows(config)
	assert.Equal(t, 6, totalRows)
}

func TestBuildHTML_EmptyCellValue(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "A"},
		},
	}

	cellValues := map[string]string{
		"0:0": "",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "<td></td>")
}

func TestBuildHTML_UnicodeContent(t *testing.T) {
	config := &ReportConfig{
		Cells: []Cell{
			{Row: 0, Col: 0, Text: "Test"},
		},
	}

	cellValues := map[string]string{
		"0:0": "你好世界 🌍",
	}

	result := buildHTML(config, cellValues, 0, 0)
	assert.Contains(t, result, "你好世界")
	assert.Contains(t, result, "🌍")
}

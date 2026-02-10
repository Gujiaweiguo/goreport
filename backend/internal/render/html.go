package render

import (
	"fmt"
	"html"
	"strings"
)

func buildHTML(config *ReportConfig, cellValues map[string]string, page, pageSize int) string {
	maxRow := 0
	maxCol := 0
	for _, cell := range config.Cells {
		if cell.Row > maxRow {
			maxRow = cell.Row
		}
		if cell.Col > maxCol {
			maxCol = cell.Col
		}
	}

	var b strings.Builder
	b.WriteString("<table>")

	startRow := 0
	endRow := maxRow + 1

	if page > 0 && pageSize > 0 {
		startRow = (page - 1) * pageSize
		endRow = page * pageSize
		if endRow > maxRow+1 {
			endRow = maxRow + 1
		}
	}

	for r := startRow; r < endRow; r++ {
		b.WriteString("<tr>")
		for c := 0; c <= maxCol; c++ {
			value := cellValues[cellKey(r, c)]
			b.WriteString("<td>")
			b.WriteString(html.EscapeString(value))
			b.WriteString("</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")
	return b.String()
}

func cellKey(row, col int) string {
	return fmt.Sprintf("%d:%d", row, col)
}

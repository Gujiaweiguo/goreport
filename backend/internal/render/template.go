package render

type ReportConfig struct {
	Cells []Cell `json:"cells"`
}

type Cell struct {
	Row          int     `json:"row"`
	Col          int     `json:"col"`
	Value        string  `json:"value"`
	Text         string  `json:"text"`
	DatasourceID *string `json:"datasourceId"`
	TableName    *string `json:"tableName"`
	FieldName    *string `json:"fieldName"`
}

func GetTotalRows(config *ReportConfig) int {
	maxRow := 0
	for _, cell := range config.Cells {
		if cell.Row > maxRow {
			maxRow = cell.Row
		}
	}
	return maxRow + 1
}

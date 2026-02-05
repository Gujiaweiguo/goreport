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

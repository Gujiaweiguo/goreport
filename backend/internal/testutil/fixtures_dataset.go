package testutil

import (
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
)

type DatasetFixtures struct {
	Datasets     []*models.Dataset
	Fields       []*models.DatasetField
	Sources      []*models.DatasetSource
	DatasourceID string
	TenantID     string
	UserID       string
}

func NewDatasetFixtures() *DatasetFixtures {
	now := time.Now()
	datasourceID := "ds-test-001"
	tenantID := "tenant-test-001"
	userID := "user-test-001"

	displayNameOrderDate := "订单日期"
	displayNameRegion := "地区"
	displayNameCategory := "类别"
	displayNameAmount := "金额"
	displayNameQuantity := "数量"
	displayNameTotalPrice := "总价"
	displayNamePricePerItem := "单价"
	displayNameID := "ID"

	exprTotalPrice := "[amount] * [quantity]"
	exprPricePerItem := "[amount] / [quantity]"

	return &DatasetFixtures{
		Datasets: []*models.Dataset{
			{
				ID:        "dataset-001",
				Name:      "Sales Data",
				Type:      "sql",
				TenantID:  tenantID,
				Status:    1,
				CreatedBy: userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "dataset-002",
				Name:      "API Orders",
				Type:      "api",
				TenantID:  tenantID,
				Status:    0,
				CreatedBy: userID,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Fields: []*models.DatasetField{
			{
				ID:          "field-001",
				DatasetID:   "dataset-001",
				Name:        "order_date",
				DisplayName: &displayNameOrderDate,
				Type:        "dimension",
				DataType:    "date",
				SortIndex:   0,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-002",
				DatasetID:   "dataset-001",
				Name:        "region",
				DisplayName: &displayNameRegion,
				Type:        "dimension",
				DataType:    "string",
				SortIndex:   1,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-003",
				DatasetID:   "dataset-001",
				Name:        "category",
				DisplayName: &displayNameCategory,
				Type:        "dimension",
				DataType:    "string",
				SortIndex:   2,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-004",
				DatasetID:   "dataset-001",
				Name:        "amount",
				DisplayName: &displayNameAmount,
				Type:        "measure",
				DataType:    "number",
				SortIndex:   3,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-005",
				DatasetID:   "dataset-001",
				Name:        "quantity",
				DisplayName: &displayNameQuantity,
				Type:        "measure",
				DataType:    "number",
				SortIndex:   4,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-006",
				DatasetID:   "dataset-001",
				Name:        "total_price",
				DisplayName: &displayNameTotalPrice,
				Type:        "measure",
				DataType:    "number",
				Expression:  &exprTotalPrice,
				SortIndex:   5,
				IsComputed:  true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-007",
				DatasetID:   "dataset-001",
				Name:        "price_per_item",
				DisplayName: &displayNamePricePerItem,
				Type:        "measure",
				DataType:    "number",
				Expression:  &exprPricePerItem,
				SortIndex:   6,
				IsComputed:  true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "field-008",
				DatasetID:   "dataset-002",
				Name:        "id",
				DisplayName: &displayNameID,
				Type:        "dimension",
				DataType:    "string",
				SortIndex:   0,
				IsComputed:  false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		Sources: []*models.DatasetSource{
			{
				ID:           "source-001",
				DatasetID:    "dataset-001",
				SourceType:   "datasource",
				SourceID:     &datasourceID,
				SourceConfig: `{"sql": "SELECT * FROM sales"}`,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			{
				ID:           "source-002",
				DatasetID:    "dataset-002",
				SourceType:   "api",
				SourceConfig: `{"url": "https://api.example.com/orders"}`,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
		},
		DatasourceID: datasourceID,
		TenantID:     tenantID,
		UserID:       userID,
	}
}

func (f *DatasetFixtures) GetDatasetByID(id string) *models.Dataset {
	for _, ds := range f.Datasets {
		if ds.ID == id {
			return ds
		}
	}
	return nil
}

func (f *DatasetFixtures) GetFieldsByDatasetID(datasetID string) []*models.DatasetField {
	var fields []*models.DatasetField
	for _, field := range f.Fields {
		if field.DatasetID == datasetID {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f *DatasetFixtures) GetDimensions(datasetID string) []*models.DatasetField {
	var fields []*models.DatasetField
	for _, field := range f.Fields {
		if field.DatasetID == datasetID && field.Type == "dimension" {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f *DatasetFixtures) GetMeasures(datasetID string) []*models.DatasetField {
	var fields []*models.DatasetField
	for _, field := range f.Fields {
		if field.DatasetID == datasetID && field.Type == "measure" {
			fields = append(fields, field)
		}
	}
	return fields
}

func (f *DatasetFixtures) GetComputedFields(datasetID string) []*models.DatasetField {
	var fields []*models.DatasetField
	for _, field := range f.Fields {
		if field.DatasetID == datasetID && field.IsComputed {
			fields = append(fields, field)
		}
	}
	return fields
}

package models

import (
	"testing"
	"time"
)

func TestDatasetField_TableName(t *testing.T) {
	f := &DatasetField{}
	if f.TableName() != "dataset_fields" {
		t.Errorf("Expected table name 'dataset_fields', got '%s'", f.TableName())
	}
}

func TestDatasetField_BasicProperties(t *testing.T) {
	displayName := "Region"
	f := &DatasetField{
		ID:          "field-1",
		DatasetID:   "dataset-1",
		Name:        "region",
		DisplayName: &displayName,
		Type:        "dimension",
		DataType:    "string",
		IsComputed:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if f.ID != "field-1" {
		t.Errorf("Expected ID 'field-1', got '%s'", f.ID)
	}
	if f.Name != "region" {
		t.Errorf("Expected name 'region', got '%s'", f.Name)
	}
	if f.Type != "dimension" {
		t.Errorf("Expected type 'dimension', got '%s'", f.Type)
	}
	if f.DataType != "string" {
		t.Errorf("Expected dataType 'string', got '%s'", f.DataType)
	}
	if f.IsComputed != false {
		t.Errorf("Expected isComputed false, got %v", f.IsComputed)
	}
}

func TestDatasetField_ComputedField(t *testing.T) {
	displayName := "Growth Rate"
	expr := "sales / last_year_sales * 100"
	f := &DatasetField{
		ID:          "field-3",
		DatasetID:   "dataset-1",
		Name:        "growth_rate",
		DisplayName: &displayName,
		Type:        "measure",
		DataType:    "number",
		IsComputed:  true,
		Expression:  &expr,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if f.IsComputed != true {
		t.Errorf("Expected isComputed true, got %v", f.IsComputed)
	}
	if f.Expression == nil {
		t.Error("Expected expression to be set")
	}
	if *f.Expression != expr {
		t.Errorf("Expected expression '%s', got '%s'", expr, *f.Expression)
	}
}

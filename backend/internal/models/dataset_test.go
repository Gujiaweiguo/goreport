package models

import (
	"testing"
	"time"
)

func TestDataset_TableName(t *testing.T) {
	d := &Dataset{}
	if d.TableName() != "datasets" {
		t.Errorf("Expected table name 'datasets', got '%s'", d.TableName())
	}
}

func TestDataset_BasicProperties(t *testing.T) {
	d := &Dataset{
		ID:        "test-id",
		TenantID:  "tenant-1",
		Name:      "Test Dataset",
		Type:      "sql",
		Config:    `{"sql": "SELECT * FROM users"}`,
		Status:    1,
		CreatedBy: "user-1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if d.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", d.ID)
	}
	if d.Name != "Test Dataset" {
		t.Errorf("Expected name 'Test Dataset', got '%s'", d.Name)
	}
	if d.Type != "sql" {
		t.Errorf("Expected type 'sql', got '%s'", d.Type)
	}
}

func TestDatasetSource_BasicProperties(t *testing.T) {
	sourceID := "datasource-1"
	s := &DatasetSource{
		ID:           "source-1",
		DatasetID:    "dataset-1",
		SourceType:   "datasource",
		SourceID:     &sourceID,
		SourceConfig: `{"host": "localhost", "port": 3306}`,
		JoinType:     "inner",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if s.ID != "source-1" {
		t.Errorf("Expected ID 'source-1', got '%s'", s.ID)
	}
	if s.SourceType != "datasource" {
		t.Errorf("Expected sourceType 'datasource', got '%s'", s.SourceType)
	}
	if s.JoinType != "inner" {
		t.Errorf("Expected joinType 'inner', got '%s'", s.JoinType)
	}
}

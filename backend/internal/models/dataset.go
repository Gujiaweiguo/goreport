package models

import (
	"time"

	"gorm.io/gorm"
)

type Dataset struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID     string         `gorm:"index;type:varchar(36);not null" json:"tenantId"`
	DatasourceID *string        `gorm:"index;type:varchar(36)" json:"datasourceId"`
	Name         string         `gorm:"type:varchar(200);not null" json:"name"`
	Type         string         `gorm:"type:enum('sql','api','file');not null" json:"type"`
	Config       string         `gorm:"type:json" json:"config"`
	Action       string         `gorm:"type:varchar(20)" json:"action,omitempty"`
	Status       int            `gorm:"type:tinyint;default:1" json:"status"`
	CreatedBy    string         `gorm:"type:varchar(36)" json:"createdBy"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Fields  []DatasetField  `gorm:"foreignKey:DatasetID" json:"fields,omitempty"`
	Sources []DatasetSource `gorm:"foreignKey:DatasetID" json:"sources,omitempty"`
}

func (Dataset) TableName() string {
	return "datasets"
}

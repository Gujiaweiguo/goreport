package models

import (
	"time"
)

type DatasetSource struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DatasetID     string    `gorm:"index;type:varchar(36);not null" json:"datasetId"`
	SourceType    string    `gorm:"type:enum('datasource','api','file');not null" json:"sourceType"`
	SourceID      *string   `gorm:"type:varchar(36)" json:"sourceId"`
	SourceConfig  string    `gorm:"type:json" json:"sourceConfig"`
	JoinType      string    `gorm:"type:enum('inner','left','right','full');default:inner" json:"joinType"`
	JoinCondition *string   `gorm:"type:text" json:"joinCondition,omitempty"`
	SortIndex     int       `gorm:"type:int;default:0" json:"sortIndex"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`

	Dataset Dataset `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
}

func (DatasetSource) TableName() string {
	return "dataset_sources"
}

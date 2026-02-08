package models

import (
	"time"
)

type DatasetField struct {
	ID               string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	DatasetID        string    `gorm:"index;type:varchar(36);not null" json:"datasetId"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	DisplayName      *string   `gorm:"type:varchar(100)" json:"displayName"`
	Type             string    `gorm:"type:enum('dimension','measure');not null" json:"type"`
	DataType         string    `gorm:"type:enum('string','number','date','boolean');not null" json:"dataType"`
	IsComputed       bool      `gorm:"type:boolean;default:false" json:"isComputed"`
	Expression       *string   `gorm:"type:text" json:"expression,omitempty"`
	IsSortable       bool      `gorm:"type:boolean;default:true" json:"isSortable"`
	IsGroupable      bool      `gorm:"type:boolean;default:true" json:"isGroupable"`
	DefaultSortOrder string    `gorm:"type:enum('asc','desc','none');default:none" json:"defaultSortOrder"`
	SortIndex        int       `gorm:"type:int;default:0" json:"sortIndex"`
	Config           string    `gorm:"type:json" json:"config,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	Dataset Dataset `gorm:"foreignKey:DatasetID" json:"dataset,omitempty"`
}

func (DatasetField) TableName() string {
	return "dataset_fields"
}

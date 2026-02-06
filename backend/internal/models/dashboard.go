package models

import (
	"time"

	"gorm.io/gorm"
)

type Dashboard struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID  string         `gorm:"index;type:varchar(36)" json:"tenantId"`
	Name      string         `gorm:"type:varchar(200)" json:"name"`
	Code      string         `gorm:"type:varchar(100)" json:"code"`
	Config    string         `gorm:"type:json" json:"config"`
	Thumbnail string         `gorm:"type:varchar(500)" json:"thumbnail"`
	Status    int            `gorm:"type:tinyint" json:"status"`
	ViewCount int            `gorm:"type:int" json:"viewCount"`
	CreatedBy string         `gorm:"type:varchar(36)" json:"createdBy"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

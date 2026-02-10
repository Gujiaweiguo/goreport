package models

import (
	"time"

	"gorm.io/gorm"
)

type Chart struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID  string         `gorm:"index;type:varchar(36)" json:"tenantId"`
	Name      string         `gorm:"type:varchar(200)" json:"name"`
	Code      string         `gorm:"type:varchar(100)" json:"code"`
	Type      string         `gorm:"type:varchar(50)" json:"type"`
	Config    string         `gorm:"type:json" json:"config"`
	Status    int            `gorm:"type:tinyint" json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

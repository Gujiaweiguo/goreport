package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Code      string         `gorm:"type:varchar(50)" json:"code"`
	Status    int            `gorm:"type:tinyint" json:"status"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

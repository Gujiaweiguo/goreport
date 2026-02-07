package models

import (
	"gorm.io/gorm"
	"time"
)

type DataSource struct {
	ID           string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name         string         `gorm:"type:varchar(100)" json:"name"`
	Type         string         `gorm:"type:varchar(20)" json:"type"`
	Host         string         `gorm:"type:varchar(255)" json:"host"`
	Port         int            `gorm:"type:int" json:"port"`
	DatabaseName string         `gorm:"column:database_name;type:varchar(100)" json:"database"`
	Username     string         `gorm:"type:varchar(100)" json:"username"`
	Password     string         `gorm:"type:varchar(255)" json:"-"`
	TenantID     string         `gorm:"column:tenant_id;index;type:varchar(36)" json:"tenantId"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

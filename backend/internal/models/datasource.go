package models

import (
	"gorm.io/gorm"
	"time"
)

type DataSource struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Type      string         `gorm:"type:varchar(20)" json:"type"`
	Host      string         `gorm:"type:varchar(255)" json:"host"`
	Port      int            `gorm:"type:int" json:"port"`
	Database  string         `gorm:"type:varchar(100)" json:"database"`
	Username  string         `gorm:"type:varchar(100)" json:"username"`
	Password  string         `gorm:"type:varchar(255)" json:"-"`
	TenantID  string         `gorm:"index;type:varchar(36)" json:"tenantId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

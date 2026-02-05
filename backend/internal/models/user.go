package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username  string         `gorm:"uniqueIndex;type:varchar(50)" json:"username"`
	Password  string         `gorm:"type:varchar(255)" json:"-"`
	Role      string         `gorm:"type:varchar(50)" json:"role"`
	TenantID  string         `gorm:"index;type:varchar(36)" json:"tenantId"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

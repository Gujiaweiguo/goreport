package models

import (
	"gorm.io/gorm"
	"time"
)

type DataSource struct {
	ID        string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string `gorm:"type:varchar(100)" json:"name"`
	Type      string `gorm:"type:varchar(20)" json:"type"`
	Host      string `gorm:"type:varchar(255)" json:"host"`
	Port      int    `gorm:"type:int" json:"port"`
	Database  string `gorm:"column:database;type:varchar(100)" json:"database"`
	Username  string `gorm:"type:varchar(100)" json:"username"`
	Password  string `gorm:"type:varchar(255)" json:"-"`
	TenantID  string `gorm:"column:tenant_id;index;type:varchar(36)" json:"tenantId"`
	CreatedBy string `gorm:"column:created_by;type:varchar(36)" json:"createdBy"`

	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	// SSH Tunnel Settings
	SSHHost      string `gorm:"column:ssh_host;type:varchar(255)" json:"sshHost,omitempty"`
	SSHPort      int    `gorm:"column:ssh_port;type:int" json:"sshPort,omitempty"`
	SSHUsername  string `gorm:"column:ssh_username;type:varchar(100)" json:"sshUsername,omitempty"`
	SSHPassword  string `gorm:"column:ssh_password;type:varchar(255)" json:"-"`
	SSHKey       string `gorm:"column:ssh_key;type:text" json:"sshKey,omitempty"`
	SSHKeyPhrase string `gorm:"column:ssh_key_phrase;type:varchar(255)" json:"sshKeyPhrase,omitempty"`

	// Runtime Controls
	MaxConnections      int `gorm:"column:max_connections;type:int;default:10" json:"maxConnections,omitempty"`
	QueryTimeoutSeconds int `gorm:"column:query_timeout_seconds;type:int;default:30" json:"queryTimeoutSeconds,omitempty"`

	// Connector Config (JSON for type-specific settings)
	Config string `gorm:"column:config;type:text" json:"config,omitempty"`

	// Metadata Description (for UI display)
	Description string `gorm:"column:description;type:text" json:"description,omitempty"`
}

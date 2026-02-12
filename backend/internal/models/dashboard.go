package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func (d *Dashboard) BeforeCreate(tx *gorm.DB) error {
	return d.serialize()
}

func (d *Dashboard) BeforeUpdate(tx *gorm.DB) error {
	return d.serialize()
}

func (d *Dashboard) AfterFind(tx *gorm.DB) error {
	return d.deserialize()
}

func (d *Dashboard) serialize() error {
	if d.Config.Width != 0 || d.Config.Height != 0 {
		configBytes, err := json.Marshal(d.Config)
		if err != nil {
			return err
		}
		d.ConfigJSON = string(configBytes)
	}
	if len(d.Components) > 0 {
		componentsBytes, err := json.Marshal(d.Components)
		if err != nil {
			return err
		}
		d.ComponentsJSON = string(componentsBytes)
	}
	return nil
}

func (d *Dashboard) deserialize() error {
	if d.ConfigJSON != "" {
		if err := json.Unmarshal([]byte(d.ConfigJSON), &d.Config); err != nil {
			d.Config = DashboardConfig{Width: 1920, Height: 1080, BackgroundColor: "#0a0e27"}
		}
	}
	if d.ComponentsJSON != "" {
		if err := json.Unmarshal([]byte(d.ComponentsJSON), &d.Components); err != nil {
			d.Components = []DashboardComponent{}
		}
	}
	return nil
}

type DashboardConfig struct {
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	BackgroundColor string `json:"backgroundColor"`
}

type DashboardComponent struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
	Width       int                    `json:"width"`
	Height      int                    `json:"height"`
	X           int                    `json:"x"`
	Y           int                    `json:"y"`
	Visible     bool                   `json:"visible"`
	Locked      bool                   `json:"locked"`
	Style       map[string]interface{} `json:"style"`
	Data        map[string]interface{} `json:"data"`
	Interaction map[string]interface{} `json:"interaction"`
}

type Dashboard struct {
	ID             string               `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID       string               `gorm:"index;type:varchar(36)" json:"tenantId"`
	Name           string               `gorm:"type:varchar(200)" json:"name"`
	Code           string               `gorm:"type:varchar(100)" json:"code"`
	Config         DashboardConfig      `gorm:"-" json:"config"`
	ConfigJSON     string               `gorm:"type:json;column:config" json:"-"`
	Components     []DashboardComponent `gorm:"-" json:"components"`
	ComponentsJSON string               `gorm:"type:json;column:components" json:"-"`
	Thumbnail      string               `gorm:"type:varchar(500)" json:"thumbnail"`
	Status         int                  `gorm:"type:tinyint" json:"status"`
	ViewCount      int                  `gorm:"type:int" json:"viewCount"`
	CreatedBy      string               `gorm:"type:varchar(36)" json:"createdBy"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt       `gorm:"index" json:"-"`
}

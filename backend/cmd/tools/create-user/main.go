package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/jeecg/jimureport-go/internal/auth"
	"github.com/jeecg/jimureport-go/internal/config"
	"github.com/jeecg/jimureport-go/internal/database"
	"github.com/jeecg/jimureport-go/internal/models"
	"gorm.io/gorm"
)

func main() {
	username := flag.String("username", "admin", "")
	password := flag.String("password", "admin123", "")
	role := flag.String("role", "admin", "")
	tenant := flag.String("tenant", "default-tenant", "")
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	db, err := database.Init(cfg.Database.DSN)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}

	if err := db.AutoMigrate(&models.User{}, &models.DataSource{}); err != nil {
		fmt.Printf("Failed to migrate users table: %v\n", err)
		return
	}

	hash, err := auth.HashPassword(*password)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		return
	}

	var user models.User
	lookup := db.Unscoped().Where("username = ?", *username).First(&user)
	if lookup.Error != nil && !errors.Is(lookup.Error, gorm.ErrRecordNotFound) {
		fmt.Printf("Failed to lookup user: %v\n", lookup.Error)
		return
	}

	if errors.Is(lookup.Error, gorm.ErrRecordNotFound) {
		user = models.User{
			ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
			Username: *username,
			Password: hash,
			Role:     *role,
			TenantID: *tenant,
		}

		if err := db.Create(&user).Error; err != nil {
			fmt.Printf("Failed to create user: %v\n", err)
			return
		}

		fmt.Printf("Created user: %s (%s)\n", user.Username, user.ID)
		return
	}

	updates := map[string]interface{}{
		"password":   hash,
		"role":       *role,
		"tenant_id":  *tenant,
		"deleted_at": nil,
	}

	if err := db.Unscoped().Model(&user).Updates(updates).Error; err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		return
	}

	fmt.Printf("Updated user: %s (%s)\n", user.Username, user.ID)
}

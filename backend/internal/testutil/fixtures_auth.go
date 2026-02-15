package testutil

import (
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthFixtures struct {
	Users   []*models.User
	Tenants []*models.Tenant
}

func NewAuthFixtures() *AuthFixtures {
	now := time.Now()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	return &AuthFixtures{
		Tenants: []*models.Tenant{
			{
				ID:        "tenant-auth-001",
				Name:      "Auth Test Tenant",
				Code:      "auth-test",
				Status:    1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Users: []*models.User{
			{
				ID:        "user-auth-001",
				Username:  "testuser",
				Password:  string(hashedPassword),
				Role:      "user",
				TenantID:  "tenant-auth-001",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "user-admin-001",
				Username:  "admin",
				Password:  string(hashedPassword),
				Role:      "admin",
				TenantID:  "tenant-auth-001",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}
}

func (f *AuthFixtures) Setup(db *gorm.DB) error {
	for _, tenant := range f.Tenants {
		if err := db.Create(tenant).Error; err != nil {
			return err
		}
	}
	for _, user := range f.Users {
		if err := db.Create(user).Error; err != nil {
			return err
		}
	}
	return nil
}

func (f *AuthFixtures) Cleanup(db *gorm.DB) error {
	tenantIDs := make([]string, len(f.Tenants))
	for i, t := range f.Tenants {
		tenantIDs[i] = t.ID
	}
	CleanupTenantData(db, tenantIDs)
	return nil
}

func (f *AuthFixtures) GetUserByUsername(username string) *models.User {
	for _, u := range f.Users {
		if u.Username == username {
			return u
		}
	}
	return nil
}

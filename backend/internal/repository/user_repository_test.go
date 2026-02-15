package repository

import (
	"context"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository(nil)
	assert.NotNil(t, repo)
}

func TestNewUserRepository_WithDB(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:invalid@tcp(localhost:3306)/test"), &gorm.Config{})
	if err != nil {
		t.Skip("Cannot connect to database")
	}
	repo := NewUserRepository(db)
	assert.NotNil(t, repo)
}

func TestUserRepository_Interface(t *testing.T) {
	var _ UserRepository = (*userRepository)(nil)
}

func TestUserRepository_GetByID(t *testing.T) {
	db, _ := setupDataSourceRepo(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user, err := repo.GetByID(ctx, "test-user-id")
	require.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepository_GetByID_ContextCancellation(t *testing.T) {
	db, _ := setupDataSourceRepo(t)
	repo := NewUserRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.GetByID(ctx, "any-id")
	assert.Error(t, err)
}

func TestUserModel_Fields(t *testing.T) {
	user := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Password: "hashedpassword",
		TenantID: "tenant-123",
	}

	assert.Equal(t, "user-123", user.ID)
	assert.Equal(t, "testuser", user.Username)
}

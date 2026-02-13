package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	repo := NewUserRepository(nil)
	assert.NotNil(t, repo)
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

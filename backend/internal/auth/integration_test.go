package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getTestDSNForAuth() string {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	return dsn
}

func skipIfNoDBForAuth(t *testing.T) {
	if getTestDSNForAuth() == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}
}

func setupAuthIntegrationTest(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := getTestDSNForAuth()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return db
}

func TestGetUserByCredentials_Success(t *testing.T) {
	skipIfNoDBForAuth(t)

	db := setupAuthIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)
	username := fmt.Sprintf("testuser-%s", ts)

	hashedPassword, err := HashPassword("password123")
	require.NoError(t, err)

	user := &models.User{
		ID:        fmt.Sprintf("u-%s", ts),
		Username:  username,
		Password:  hashedPassword,
		TenantID:  "test-auth-tenant",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	result, err := GetUserByCredentials(db, username, "password123")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, username, result.Username)
	assert.Equal(t, user.ID, result.ID)
}

func TestGetUserByCredentials_WrongPassword(t *testing.T) {
	skipIfNoDBForAuth(t)

	db := setupAuthIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	now := time.Now()
	ts := fmt.Sprintf("%08x", now.UnixNano()&0xFFFFFFFF)
	username := fmt.Sprintf("testuser-wp-%s", ts)

	hashedPassword, err := HashPassword("correct-password")
	require.NoError(t, err)

	user := &models.User{
		ID:        fmt.Sprintf("u-wp-%s", ts),
		Username:  username,
		Password:  hashedPassword,
		TenantID:  "test-auth-tenant",
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	result, err := GetUserByCredentials(db, username, "wrong-password")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid password")
}

func TestGetUserByCredentials_UserNotFound(t *testing.T) {
	skipIfNoDBForAuth(t)

	db := setupAuthIntegrationTest(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	result, err := GetUserByCredentials(db, "non-existent-user-xyz", "password")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("with roles set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Set(string(RolesKey), []string{"admin", "user"})

		roles := GetRoles(c)

		assert.Len(t, roles, 2)
		assert.Contains(t, roles, "admin")
		assert.Contains(t, roles, "user")
	})

	t.Run("without roles set", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)

		roles := GetRoles(c)

		assert.Nil(t, roles)
	})
}

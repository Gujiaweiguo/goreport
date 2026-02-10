package auth

import (
	"testing"

	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateToken(t *testing.T) {
	InitJWT(&config.JWTConfig{
		Secret:   "test-secret",
		Issuer:   "goreport-test",
		Audience: "goreport-test",
	})

	user := &models.User{
		ID:       "u-1",
		Username: "alice",
		Role:     "admin",
		TenantID: "tenant-1",
	}

	token, err := GenerateToken(user)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.TenantID, claims.TenantID)
	require.Len(t, claims.Roles, 1)
	assert.Equal(t, user.Role, claims.Roles[0])
}

func TestValidateToken_InvalidToken(t *testing.T) {
	InitJWT(&config.JWTConfig{
		Secret:   "test-secret",
		Issuer:   "goreport-test",
		Audience: "goreport-test",
	})

	_, err := ValidateToken("invalid-token")
	assert.Error(t, err)
}

func TestHashAndCheckPassword(t *testing.T) {
	plain := "admin123"

	hashed, err := HashPassword(plain)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	assert.True(t, CheckPassword(plain, hashed))
	assert.False(t, CheckPassword("wrong-password", hashed))
}

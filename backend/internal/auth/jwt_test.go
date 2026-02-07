package auth

import (
	"testing"

	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/models"
)

func TestHashPasswordAndCheckPassword(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Fatal("CheckPassword should return true for correct password")
	}

	if CheckPassword("wrong", hash) {
		t.Fatal("CheckPassword should return false for incorrect password")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	cfg := &config.JWTConfig{
		Secret:   "test-secret",
		Issuer:   "test-issuer",
		Audience: "test-audience",
	}

	InitJWT(cfg)

	user := &models.User{
		ID:       "user-1",
		Username: "alice",
		Role:     "admin",
		TenantID: "tenant-1",
	}

	token, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != user.ID {
		t.Fatalf("expected userId %s, got %s", user.ID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Fatalf("expected username %s, got %s", user.Username, claims.Username)
	}
	if claims.TenantID != user.TenantID {
		t.Fatalf("expected tenantId %s, got %s", user.TenantID, claims.TenantID)
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != user.Role {
		t.Fatalf("expected roles [%s], got %v", user.Role, claims.Roles)
	}
}

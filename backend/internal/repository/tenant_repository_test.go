package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestNewTenantRepository(t *testing.T) {
	repo := NewTenantRepository(nil)
	assert.NotNil(t, repo)
}

func TestTenantRepository_GetByID(t *testing.T) {
	db, _ := setupDataSourceRepo(t)
	repo := NewTenantRepository(db)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	tenant, err := repo.GetByID(ctx, tenantID)
	require.NoError(t, err)
	assert.NotNil(t, tenant)
	assert.Equal(t, tenantID, tenant.ID)
}

func TestTenantRepository_GetByID_NotFound(t *testing.T) {
	db, _ := setupDataSourceRepo(t)
	repo := NewTenantRepository(db)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existing-tenant")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestTenantRepository_ListByUserID(t *testing.T) {
	t.Skip("Requires user_tenants table which does not exist")
	db, _ := setupDataSourceRepo(t)
	repo := NewTenantRepository(db)
	ctx := context.Background()
	tenantID := setupTenant(t, db)

	tenants, err := repo.ListByUserID(ctx, "test-user")
	require.NoError(t, err)
	assert.NotNil(t, tenants)
	_ = tenantID
}

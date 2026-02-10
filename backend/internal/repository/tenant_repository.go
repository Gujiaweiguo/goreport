package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type TenantRepository interface {
	GetByID(ctx context.Context, id string) (*models.Tenant, error)
	ListByUserID(ctx context.Context, userID string) ([]*models.Tenant, error)
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) ListByUserID(ctx context.Context, userID string) ([]*models.Tenant, error) {
	var tenants []*models.Tenant
	err := r.db.WithContext(ctx).
		Table("tenants t").
		Select("t.*").
		Joins("JOIN user_tenants ut ON ut.tenant_id = t.id").
		Where("ut.user_id = ?", userID).
		Order("ut.is_default DESC, t.created_at DESC").
		Find(&tenants).Error
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

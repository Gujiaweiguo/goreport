package chart

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, chart *models.Chart) error
	Update(ctx context.Context, chart *models.Chart) error
	Delete(ctx context.Context, id, tenantID string) error
	Get(ctx context.Context, id, tenantID string) (*models.Chart, error)
	List(ctx context.Context, tenantID string) ([]*models.Chart, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, chart *models.Chart) error {
	return r.db.WithContext(ctx).Create(chart).Error
}

func (r *repository) Update(ctx context.Context, chart *models.Chart) error {
	return r.db.WithContext(ctx).Model(&models.Chart{}).
		Where("id = ? AND tenant_id = ?", chart.ID, chart.TenantID).
		Updates(chart).Error
}

func (r *repository) Delete(ctx context.Context, id, tenantID string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&models.Chart{}).Error
}

func (r *repository) Get(ctx context.Context, id, tenantID string) (*models.Chart, error) {
	var chart models.Chart
	if err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&chart).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

func (r *repository) List(ctx context.Context, tenantID string) ([]*models.Chart, error) {
	var charts []*models.Chart
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("updated_at desc").
		Find(&charts).Error; err != nil {
		return nil, err
	}
	return charts, nil
}

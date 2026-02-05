package report

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, report *Report) error
	Update(ctx context.Context, report *Report) error
	Delete(ctx context.Context, id, tenantID string) error
	Get(ctx context.Context, id, tenantID string) (*Report, error)
	List(ctx context.Context, tenantID string) ([]*Report, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, report *Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *repository) Update(ctx context.Context, report *Report) error {
	return r.db.WithContext(ctx).Model(&Report{}).
		Where("id = ? AND tenant_id = ?", report.ID, report.TenantID).
		Updates(report).Error
}

func (r *repository) Delete(ctx context.Context, id, tenantID string) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&Report{}).Error
}

func (r *repository) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	var report Report
	if err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *repository) List(ctx context.Context, tenantID string) ([]*Report, error) {
	var reports []*Report
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("updated_at desc").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

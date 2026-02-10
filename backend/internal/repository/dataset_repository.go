package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DatasetRepository interface {
	Create(ctx context.Context, dataset *models.Dataset) error
	GetByID(ctx context.Context, id string) (*models.Dataset, error)
	GetByIDWithFields(ctx context.Context, id string) (*models.Dataset, error)
	List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error)
	Update(ctx context.Context, dataset *models.Dataset) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}

type datasetRepository struct {
	db *gorm.DB
}

func NewDatasetRepository(db *gorm.DB) DatasetRepository {
	return &datasetRepository{db: db}
}

func (r *datasetRepository) Create(ctx context.Context, dataset *models.Dataset) error {
	return r.db.WithContext(ctx).Create(dataset).Error
}

func (r *datasetRepository) GetByID(ctx context.Context, id string) (*models.Dataset, error) {
	var dataset models.Dataset
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

func (r *datasetRepository) GetByIDWithFields(ctx context.Context, id string) (*models.Dataset, error) {
	var dataset models.Dataset
	err := r.db.WithContext(ctx).Preload("Fields").Preload("Sources").Where("id = ?", id).First(&dataset).Error
	if err != nil {
		return nil, err
	}
	return &dataset, nil
}

func (r *datasetRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error) {
	var datasets []*models.Dataset
	var total int64

	offset := (page - 1) * pageSize

	err := r.db.WithContext(ctx).Model(&models.Dataset{}).Where("tenant_id = ?", tenantID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).
		Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&datasets).Error

	if err != nil {
		return nil, 0, err
	}

	return datasets, total, nil
}

func (r *datasetRepository) Update(ctx context.Context, dataset *models.Dataset) error {
	return r.db.WithContext(ctx).Model(dataset).Updates(dataset).Error
}

func (r *datasetRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.Dataset{}).Error
}

func (r *datasetRepository) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Dataset{}).Error
}

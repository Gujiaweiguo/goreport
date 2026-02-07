package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DataSourceRepository interface {
	Create(ctx context.Context, ds *models.DataSource) error
	GetByID(ctx context.Context, id string) (*models.DataSource, error)
	List(ctx context.Context, tenantID string) ([]*models.DataSource, error)
	Update(ctx context.Context, ds *models.DataSource) error
	Delete(ctx context.Context, id string) error
}

type dataSourceRepository struct {
	db *gorm.DB
}

func NewDataSourceRepository(db *gorm.DB) DataSourceRepository {
	return &dataSourceRepository{db: db}
}

func (r *dataSourceRepository) Create(ctx context.Context, ds *models.DataSource) error {
	return r.db.WithContext(ctx).Create(ds).Error
}

func (r *dataSourceRepository) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	var ds models.DataSource
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ds).Error
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

func (r *dataSourceRepository) List(ctx context.Context, tenantID string) ([]*models.DataSource, error) {
	var datasources []*models.DataSource
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&datasources).Error
	if err != nil {
		return nil, err
	}
	return datasources, nil
}

func (r *dataSourceRepository) Update(ctx context.Context, ds *models.DataSource) error {
	return r.db.WithContext(ctx).Model(ds).Updates(ds).Error
}

func (r *dataSourceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DataSource{}).Error
}

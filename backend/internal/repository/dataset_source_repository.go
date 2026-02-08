package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DatasetSourceRepository interface {
	Create(ctx context.Context, source *models.DatasetSource) error
	GetByID(ctx context.Context, id string) (*models.DatasetSource, error)
	List(ctx context.Context, datasetID string) ([]*models.DatasetSource, error)
	Update(ctx context.Context, source *models.DatasetSource) error
	Delete(ctx context.Context, id string) error
}

type datasetSourceRepository struct {
	db *gorm.DB
}

func NewDatasetSourceRepository(db *gorm.DB) DatasetSourceRepository {
	return &datasetSourceRepository{db: db}
}

func (r *datasetSourceRepository) Create(ctx context.Context, source *models.DatasetSource) error {
	return r.db.WithContext(ctx).Create(source).Error
}

func (r *datasetSourceRepository) GetByID(ctx context.Context, id string) (*models.DatasetSource, error) {
	var source models.DatasetSource
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&source).Error
	if err != nil {
		return nil, err
	}
	return &source, nil
}

func (r *datasetSourceRepository) List(ctx context.Context, datasetID string) ([]*models.DatasetSource, error) {
	var sources []*models.DatasetSource
	err := r.db.WithContext(ctx).Where("dataset_id = ?", datasetID).
		Order("sort_index ASC").
		Find(&sources).Error
	if err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *datasetSourceRepository) Update(ctx context.Context, source *models.DatasetSource) error {
	return r.db.WithContext(ctx).Model(source).Updates(source).Error
}

func (r *datasetSourceRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DatasetSource{}).Error
}

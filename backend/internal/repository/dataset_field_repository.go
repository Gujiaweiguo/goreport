package repository

import (
	"context"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DatasetFieldRepository interface {
	Create(ctx context.Context, field *models.DatasetField) error
	GetByID(ctx context.Context, id string) (*models.DatasetField, error)
	List(ctx context.Context, datasetID string) ([]*models.DatasetField, error)
	ListByType(ctx context.Context, datasetID string, fieldType string) ([]*models.DatasetField, error)
	Update(ctx context.Context, field *models.DatasetField) error
	Delete(ctx context.Context, id string) error
	DeleteComputedFields(ctx context.Context, datasetID string) error
}

type datasetFieldRepository struct {
	db *gorm.DB
}

func NewDatasetFieldRepository(db *gorm.DB) DatasetFieldRepository {
	return &datasetFieldRepository{db: db}
}

func (r *datasetFieldRepository) Create(ctx context.Context, field *models.DatasetField) error {
	return r.db.WithContext(ctx).Create(field).Error
}

func (r *datasetFieldRepository) GetByID(ctx context.Context, id string) (*models.DatasetField, error) {
	var field models.DatasetField
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&field).Error
	if err != nil {
		return nil, err
	}
	return &field, nil
}

func (r *datasetFieldRepository) List(ctx context.Context, datasetID string) ([]*models.DatasetField, error) {
	var fields []*models.DatasetField
	err := r.db.WithContext(ctx).Where("dataset_id = ?", datasetID).
		Order("sort_index ASC").
		Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (r *datasetFieldRepository) ListByType(ctx context.Context, datasetID string, fieldType string) ([]*models.DatasetField, error) {
	var fields []*models.DatasetField
	err := r.db.WithContext(ctx).Where("dataset_id = ? AND type = ?", datasetID, fieldType).
		Order("sort_index ASC").
		Find(&fields).Error
	if err != nil {
		return nil, err
	}
	return fields, nil
}

func (r *datasetFieldRepository) Update(ctx context.Context, field *models.DatasetField) error {
	return r.db.WithContext(ctx).Model(field).Updates(field).Error
}

func (r *datasetFieldRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.DatasetField{}).Error
}

func (r *datasetFieldRepository) DeleteComputedFields(ctx context.Context, datasetID string) error {
	return r.db.WithContext(ctx).Where("dataset_id = ? AND is_computed = ?", datasetID, true).Delete(&models.DatasetField{}).Error
}

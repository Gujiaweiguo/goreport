package dashboard

import (
	"errors"

	"github.com/jeecg/jimureport-go/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(dashboard *models.Dashboard) error
	Update(dashboard *models.Dashboard) error
	Delete(id, tenantID string) error
	Get(id, tenantID string) (*models.Dashboard, error)
	List(tenantID string) ([]*models.Dashboard, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(dashboard *models.Dashboard) error {
	return r.db.Create(dashboard).Error
}

func (r *repository) Update(dashboard *models.Dashboard) error {
	return r.db.Save(dashboard).Error
}

func (r *repository) Delete(id, tenantID string) error {
	result := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Dashboard{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("dashboard not found")
	}
	return nil
}

func (r *repository) Get(id, tenantID string) (*models.Dashboard, error) {
	var dashboard models.Dashboard
	err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&dashboard).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("dashboard not found")
		}
		return nil, err
	}
	return &dashboard, nil
}

func (r *repository) List(tenantID string) ([]*models.Dashboard, error) {
	var dashboards []*models.Dashboard
	err := r.db.Where("tenant_id = ?", tenantID).Find(&dashboards).Error
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

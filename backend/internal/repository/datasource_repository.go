package repository

import (
	"context"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DatasourceRepository interface {
	Create(ctx context.Context, ds *models.DataSource) error
	GetByID(ctx context.Context, id string) (*models.DataSource, error)
	List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error)
	Update(ctx context.Context, ds *models.DataSource) error
	Delete(ctx context.Context, id, tenantID string) error
	Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error)
	Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error)
	Move(ctx context.Context, id, tenantID string) error
	Rename(ctx context.Context, id, tenantID string, newName string) error
}

type datasourceRepository struct {
	db *gorm.DB
}

func NewDatasourceRepository(db *gorm.DB) DatasourceRepository {
	return &datasourceRepository{db: db}
}

func (r *datasourceRepository) Create(ctx context.Context, ds *models.DataSource) error {
	ds.CreatedAt = time.Now()
	ds.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Create(ds).Error
}

func (r *datasourceRepository) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	var ds models.DataSource
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ds).Error
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

func (r *datasourceRepository) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	var datasources []*models.DataSource
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DataSource{}).Where("tenant_id = ?", tenantID)

	query.Count(&total)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&datasources).Error
	if err != nil {
		return nil, 0, err
	}

	return datasources, total, nil
}

func (r *datasourceRepository) Update(ctx context.Context, ds *models.DataSource) error {
	ds.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(ds).Updates(ds).Error
}

func (r *datasourceRepository) Delete(ctx context.Context, id, tenantID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.DataSource{}).Error
}

func (r *datasourceRepository) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	var datasources []*models.DataSource
	var total int64

	query := r.db.WithContext(ctx).Model(&models.DataSource{}).Where("tenant_id = ?", tenantID)

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&datasources).Error
	if err != nil {
		return nil, 0, err
	}

	return datasources, total, nil
}

func (r *datasourceRepository) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	original, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if original.TenantID != tenantID {
		return nil, gorm.ErrRecordNotFound
	}

	copy := &models.DataSource{
		Name:                original.Name + " (副本)",
		Type:                original.Type,
		Host:                original.Host,
		Port:                original.Port,
		Username:            original.Username,
		Password:            original.Password,
		Database:            original.Database,
		SSHHost:             original.SSHHost,
		SSHPort:             original.SSHPort,
		SSHUsername:         original.SSHUsername,
		SSHPassword:         original.SSHPassword,
		SSHKey:              original.SSHKey,
		SSHKeyPhrase:        original.SSHKeyPhrase,
		MaxConnections:      original.MaxConnections,
		QueryTimeoutSeconds: original.QueryTimeoutSeconds,
		Config:              original.Config,
		TenantID:            tenantID,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := r.Create(ctx, copy); err != nil {
		return nil, err
	}
	return copy, nil
}

func (r *datasourceRepository) Move(ctx context.Context, id, tenantID string) error {
	ds, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ds.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}

	ds.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", ds.ID).Updates(ds).Error
}

func (r *datasourceRepository) Rename(ctx context.Context, id, tenantID string, newName string) error {
	ds, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if ds.TenantID != tenantID {
		return gorm.ErrRecordNotFound
	}

	if newName == "" || len(newName) > 255 {
		return gorm.ErrInvalidData
	}

	existing, _, err := r.List(ctx, tenantID, 1, 1000)
	if err != nil {
		return err
	}

	for _, other := range existing {
		if other.ID != id && other.Name == newName {
			return gorm.ErrDuplicatedKey
		}
	}

	ds.Name = newName
	ds.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(ds).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(ds).Error
}

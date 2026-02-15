package testutil

import (
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/gorm"
)

type DatasourceFixtures struct {
	Datasources []*models.DataSource
	TenantID    string
	UserID      string
}

func NewDatasourceFixtures() *DatasourceFixtures {
	now := time.Now()
	tenantID := "tenant-ds-001"
	userID := "user-ds-001"

	return &DatasourceFixtures{
		TenantID: tenantID,
		UserID:   userID,
		Datasources: []*models.DataSource{
			{
				ID:                  "ds-test-001",
				Name:                "Test MySQL",
				Type:                "mysql",
				Host:                "localhost",
				Port:                3306,
				Database:            "test_db",
				Username:            "root",
				Password:            "root",
				TenantID:            tenantID,
				CreatedBy:           userID,
				MaxConnections:      10,
				QueryTimeoutSeconds: 30,
				CreatedAt:           now,
				UpdatedAt:           now,
			},
			{
				ID:                  "ds-test-002",
				Name:                "Test PostgreSQL",
				Type:                "postgresql",
				Host:                "localhost",
				Port:                5432,
				Database:            "test_db",
				Username:            "postgres",
				Password:            "postgres",
				TenantID:            tenantID,
				CreatedBy:           userID,
				MaxConnections:      10,
				QueryTimeoutSeconds: 30,
				CreatedAt:           now,
				UpdatedAt:           now,
			},
		},
	}
}

func (f *DatasourceFixtures) Setup(db *gorm.DB) error {
	for _, ds := range f.Datasources {
		if err := db.Create(ds).Error; err != nil {
			return err
		}
	}
	return nil
}

func (f *DatasourceFixtures) Cleanup(db *gorm.DB) error {
	ids := make([]string, len(f.Datasources))
	for i, ds := range f.Datasources {
		ids[i] = ds.ID
	}
	return db.Where("id IN ?", ids).Delete(&models.DataSource{}).Error
}

func (f *DatasourceFixtures) GetDatasourceByName(name string) *models.DataSource {
	for _, ds := range f.Datasources {
		if ds.Name == name {
			return ds
		}
	}
	return nil
}

func (f *DatasourceFixtures) GetDatasourceByID(id string) *models.DataSource {
	for _, ds := range f.Datasources {
		if ds.ID == id {
			return ds
		}
	}
	return nil
}

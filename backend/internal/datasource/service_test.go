package datasource

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDatasourceRepo struct {
	datasource  *models.DataSource
	datasources []*models.DataSource
	createErr   error
	updateErr   error
	deleteErr   error
	getErr      error
	listErr     error
	searchErr   error
	copyErr     error
	moveErr     error
	renameErr   error
}

type mockProfileValidator struct {
	mock.Mock
}

func (m *mockProfileValidator) Validate(datasourceType string, config map[string]interface{}) error {
	args := m.Called(datasourceType, config)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *mockDatasourceRepo) Create(ctx context.Context, datasource *models.DataSource) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.datasource = datasource
	return nil
}

func (m *mockDatasourceRepo) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.datasource != nil && m.datasource.ID == id {
		return m.datasource, nil
	}
	return nil, errors.New("datasource not found")
}

func (m *mockDatasourceRepo) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if m.listErr != nil {
		return nil, 0, m.listErr
	}
	return m.datasources, int64(len(m.datasources)), nil
}

func (m *mockDatasourceRepo) Update(ctx context.Context, datasource *models.DataSource) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.datasource = datasource
	return nil
}

func (m *mockDatasourceRepo) Delete(ctx context.Context, id, tenantID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.datasource = nil
	return nil
}

func (m *mockDatasourceRepo) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if m.searchErr != nil {
		return nil, 0, m.searchErr
	}
	return m.datasources, int64(len(m.datasources)), nil
}

func (m *mockDatasourceRepo) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	if m.copyErr != nil {
		return nil, m.copyErr
	}
	if m.datasource != nil {
		copyDs := *m.datasource
		return &copyDs, nil
	}
	return nil, errors.New("datasource not found")
}

func (m *mockDatasourceRepo) Move(ctx context.Context, id, tenantID string) error {
	if m.moveErr != nil {
		return m.moveErr
	}
	return nil
}

func (m *mockDatasourceRepo) Rename(ctx context.Context, id, tenantID, newName string) error {
	if m.renameErr != nil {
		return m.renameErr
	}
	if m.datasource != nil {
		m.datasource.Name = newName
		return nil
	}
	return errors.New("datasource not found")
}

func TestService_Create(t *testing.T) {
	repo := &mockDatasourceRepo{}
	service := NewService(repo)

	t.Run("成功创建MySQL数据源", func(t *testing.T) {
		req := &CreateRequest{
			Name:      "MySQL Test",
			Type:      "mysql",
			Host:      "localhost",
			Port:      3306,
			Database:  "testdb",
			Username:  "root",
			Password:  "password",
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "MySQL Test", datasource.Name)
		assert.Equal(t, "mysql", datasource.Type)
	})

	t.Run("成功创建PostgreSQL数据源", func(t *testing.T) {
		req := &CreateRequest{
			Name:      "PostgreSQL Test",
			Type:      "postgres",
			Host:      "localhost",
			Port:      5432,
			Database:  "testdb",
			Username:  "postgres",
			Password:  "password",
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, "postgres", datasource.Type)
	})

	t.Run("创建失败", func(t *testing.T) {
		repo.createErr = errors.New("database error")

		req := &CreateRequest{
			Name:     "Test",
			Type:     "mysql",
			Host:     "localhost",
			Port:     3306,
			TenantID: "tenant-1",
		}

		_, err := service.Create(context.Background(), req)

		assert.Error(t, err)
	})
}

func TestService_GetByID(t *testing.T) {
	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		Name:     "Test Datasource",
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
		TenantID: "tenant-1",
	}

	repo := &mockDatasourceRepo{
		datasource: existingDatasource,
	}
	service := NewService(repo)

	t.Run("成功获取数据源", func(t *testing.T) {
		datasource, err := service.GetByID(context.Background(), "ds-1")

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "ds-1", datasource.ID)
		assert.Equal(t, "Test Datasource", datasource.Name)
	})

	t.Run("数据源不存在", func(t *testing.T) {
		repo.getErr = errors.New("not found")

		_, err := service.GetByID(context.Background(), "not-exist")

		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	datasources := []*models.DataSource{
		{
			ID:       "ds-1",
			Name:     "MySQL Datasource",
			Type:     "mysql",
			TenantID: "tenant-1",
		},
		{
			ID:       "ds-2",
			Name:     "PostgreSQL Datasource",
			Type:     "postgres",
			TenantID: "tenant-1",
		},
	}

	repo := &mockDatasourceRepo{
		datasources: datasources,
	}
	service := NewService(repo)

	t.Run("成功获取数据源列表", func(t *testing.T) {
		list, total, err := service.List(context.Background(), "tenant-1", 1, 10)

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, int64(2), total)
		assert.Equal(t, "MySQL Datasource", list[0].Name)
	})

	t.Run("空列表", func(t *testing.T) {
		repo.datasources = []*models.DataSource{}

		list, total, err := service.List(context.Background(), "tenant-1", 1, 10)

		assert.NoError(t, err)
		assert.Len(t, list, 0)
		assert.Equal(t, int64(0), total)
	})

	t.Run("列表获取失败", func(t *testing.T) {
		repo.listErr = errors.New("database error")

		_, _, err := service.List(context.Background(), "tenant-1", 1, 10)

		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		TenantID: "tenant-1",
		Name:     "Old Name",
		Type:     "mysql",
	}

	repo := &mockDatasourceRepo{
		datasource: existingDatasource,
	}
	service := NewService(repo)

	t.Run("成功更新数据源", func(t *testing.T) {
		req := &UpdateRequest{
			ID:       "ds-1",
			Name:     "Updated Name",
			Host:     "newhost.com",
			Port:     3307,
			TenantID: "tenant-1",
		}

		datasource, err := service.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "Updated Name", datasource.Name)
	})

	t.Run("更新失败", func(t *testing.T) {
		repo.updateErr = errors.New("update failed")

		req := &UpdateRequest{
			ID:       "ds-1",
			Name:     "Test",
			TenantID: "tenant-1",
		}

		_, err := service.Update(context.Background(), req)

		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		Name:     "Test Datasource",
		Type:     "mysql",
		TenantID: "tenant-1",
	}

	repo := &mockDatasourceRepo{
		datasource: existingDatasource,
	}
	service := NewService(repo)

	t.Run("成功删除数据源", func(t *testing.T) {
		err := service.Delete(context.Background(), "ds-1", "tenant-1")

		assert.NoError(t, err)
	})

	t.Run("删除失败", func(t *testing.T) {
		repo.deleteErr = errors.New("delete failed")

		err := service.Delete(context.Background(), "ds-1", "tenant-1")

		assert.Error(t, err)
	})
}

func TestService_Search(t *testing.T) {
	datasources := []*models.DataSource{
		{
			ID:       "ds-1",
			Name:     "Production MySQL",
			Type:     "mysql",
			TenantID: "tenant-1",
		},
		{
			ID:       "ds-2",
			Name:     "Staging MySQL",
			Type:     "mysql",
			TenantID: "tenant-1",
		},
	}

	repo := &mockDatasourceRepo{
		datasources: datasources,
	}
	service := NewService(repo)

	t.Run("成功搜索数据源", func(t *testing.T) {
		list, total, err := service.Search(context.Background(), "tenant-1", "MySQL", 1, 10)

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("搜索失败", func(t *testing.T) {
		repo.searchErr = errors.New("search failed")

		_, _, err := service.Search(context.Background(), "tenant-1", "test", 1, 10)

		assert.Error(t, err)
	})
}

func TestService_Copy(t *testing.T) {
	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		Name:     "Original",
		Type:     "mysql",
		TenantID: "tenant-1",
	}

	repo := &mockDatasourceRepo{
		datasource: existingDatasource,
	}
	service := NewService(repo)

	t.Run("成功复制数据源", func(t *testing.T) {
		datasource, err := service.Copy(context.Background(), "ds-1", "tenant-1")

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "Original", datasource.Name)
	})

	t.Run("复制失败", func(t *testing.T) {
		repo.copyErr = errors.New("copy failed")

		_, err := service.Copy(context.Background(), "not-exist", "tenant-1")

		assert.Error(t, err)
	})
}

func TestService_Move(t *testing.T) {
	existingDatasource := &models.DataSource{
		ID:       "ds-1",
		Name:     "Test Datasource",
		Type:     "mysql",
		TenantID: "tenant-1",
	}

	repo := &mockDatasourceRepo{
		datasource: existingDatasource,
	}
	service := NewService(repo)

	t.Run("成功移动数据源", func(t *testing.T) {
		err := service.Move(context.Background(), "ds-1", "tenant-1")

		assert.NoError(t, err)
	})

	t.Run("移动失败", func(t *testing.T) {
		repo.moveErr = errors.New("move failed")

		err := service.Move(context.Background(), "ds-1", "tenant-1")

		assert.Error(t, err)
	})
}

func TestService_Rename(t *testing.T) {
	t.Run("成功重命名数据源", func(t *testing.T) {
		datasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: datasource,
		}
		service := NewService(repo)

		result, err := service.Rename(context.Background(), "ds-1", "tenant-1", "New Name")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("重命名失败-数据源不存在", func(t *testing.T) {
		datasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: datasource,
			getErr:     errors.New("not found"),
		}
		service := NewService(repo)

		_, err := service.Rename(context.Background(), "not-exist", "tenant-1", "Test")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("重命名失败-空名字", func(t *testing.T) {
		datasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: datasource,
		}
		service := NewService(repo)

		_, err := service.Rename(context.Background(), "ds-1", "tenant-1", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("重命名失败-名字太长", func(t *testing.T) {
		datasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: datasource,
		}
		service := NewService(repo)

		longName := string(make([]byte, 256))
		_, err := service.Rename(context.Background(), "ds-1", "tenant-1", longName)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "too long")
	})

	t.Run("重命名失败-更新错误", func(t *testing.T) {
		datasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: datasource,
			updateErr:  errors.New("update failed"),
		}
		service := NewService(repo)

		_, err := service.Rename(context.Background(), "ds-1", "tenant-1", "New Name")

		assert.Error(t, err)
	})
}

func TestService_Create_WithAdvanced(t *testing.T) {
	t.Run("成功创建带SSH配置的数据源", func(t *testing.T) {
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		sshHost := "ssh.example.com"
		sshPort := 22
		req := &CreateRequest{
			Name:      "MySQL with SSH",
			Type:      "mysql",
			Host:      "localhost",
			Port:      3306,
			Database:  "testdb",
			Username:  "root",
			Password:  "password",
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
			Advanced: &AdvancedConfig{
				SSHHost:             sshHost,
				SSHPort:             sshPort,
				SSHUsername:         "sshuser",
				SSHPassword:         "sshpass",
				MaxConnections:      10,
				QueryTimeoutSeconds: 30,
			},
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "MySQL with SSH", datasource.Name)
		assert.Equal(t, sshHost, datasource.SSHHost)
		assert.Equal(t, sshPort, datasource.SSHPort)
		assert.Equal(t, 10, datasource.MaxConnections)
		assert.Equal(t, 30, datasource.QueryTimeoutSeconds)
	})

	t.Run("成功创建MongoDB数据源", func(t *testing.T) {
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		req := &CreateRequest{
			Name:      "MongoDB Test",
			Type:      "mongodb",
			Host:      "localhost",
			Port:      27017,
			Database:  "testdb",
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "MongoDB Test", datasource.Name)
		assert.Equal(t, "mongodb", datasource.Type)
	})
}

func TestService_Update_WithAdvanced(t *testing.T) {
	t.Run("成功更新带SSH配置", func(t *testing.T) {
		existingDatasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Old Name",
			Type:     "mysql",
			TenantID: "tenant-1",
		}

		repo := &mockDatasourceRepo{
			datasource: existingDatasource,
		}
		service := NewService(repo)

		req := &UpdateRequest{
			ID:       "ds-1",
			Name:     "Updated Name",
			Host:     "newhost.com",
			Port:     3307,
			TenantID: "tenant-1",
			Advanced: &AdvancedConfig{
				SSHHost:             "new.ssh.com",
				SSHPort:             2222,
				SSHUsername:         "newsshuser",
				MaxConnections:      20,
				QueryTimeoutSeconds: 60,
			},
		}

		datasource, err := service.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "Updated Name", datasource.Name)
		assert.Equal(t, "new.ssh.com", datasource.SSHHost)
		assert.Equal(t, 2222, datasource.SSHPort)
		assert.Equal(t, 20, datasource.MaxConnections)
	})

	t.Run("更新失败-数据源不属于租户", func(t *testing.T) {
		existingDatasource := &models.DataSource{
			ID:       "ds-1",
			Name:     "Test",
			Type:     "mysql",
			TenantID: "tenant-2",
		}

		repo := &mockDatasourceRepo{
			datasource: existingDatasource,
		}
		service := NewService(repo)

		req := &UpdateRequest{
			ID:       "ds-1",
			Name:     "Updated Name",
			TenantID: "tenant-1",
		}

		_, err := service.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestService_List_Pagination(t *testing.T) {
	t.Run("列表分页-第一页", func(t *testing.T) {
		datasources := []*models.DataSource{
			{ID: "ds-1", Name: "DS1", TenantID: "tenant-1"},
			{ID: "ds-2", Name: "DS2", TenantID: "tenant-1"},
		}

		repo := &mockDatasourceRepo{
			datasources: datasources,
		}
		service := NewService(repo)

		list, total, err := service.List(context.Background(), "tenant-1", 1, 2)

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("列表分页-page小于等于0", func(t *testing.T) {
		repo := &mockDatasourceRepo{
			datasources: []*models.DataSource{{ID: "ds-1", Name: "DS1", TenantID: "tenant-1"}},
		}
		service := NewService(repo)

		list, total, err := service.List(context.Background(), "tenant-1", 0, 10)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(list))
		assert.Equal(t, int64(1), total)
	})

	t.Run("列表分页-pageSize小于等于0", func(t *testing.T) {
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		list, total, err := service.List(context.Background(), "tenant-1", 1, 0)

		assert.NoError(t, err)
		assert.Len(t, list, 0)
		assert.Equal(t, int64(0), total)
	})

	t.Run("列表分页-pageSize大于100", func(t *testing.T) {
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		list, total, err := service.List(context.Background(), "tenant-1", 1, 200)

		assert.NoError(t, err)
		assert.Len(t, list, 0)
		assert.Equal(t, int64(0), total)
	})
}

func TestService_Search_EmptyKeyword(t *testing.T) {
	t.Run("搜索-空关键字", func(t *testing.T) {
		datasources := []*models.DataSource{
			{ID: "ds-1", Name: "DS1", TenantID: "tenant-1"},
		}

		repo := &mockDatasourceRepo{
			datasources: datasources,
		}
		service := NewService(repo)

		list, total, err := service.Search(context.Background(), "tenant-1", "", 1, 10)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})

	t.Run("搜索分页", func(t *testing.T) {
		datasources := []*models.DataSource{
			{ID: "ds-1", Name: "MySQL", TenantID: "tenant-1"},
		}

		repo := &mockDatasourceRepo{
			datasources: datasources,
		}
		service := NewService(repo)

		list, total, err := service.Search(context.Background(), "tenant-1", "MySQL", 1, 10)

		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, int64(1), total)
	})
}

func TestService_Create_SpecialTypes(t *testing.T) {
	t.Run("成功创建API数据源", func(t *testing.T) {
		t.Skip("Requires validator mock implementation")
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		req := &CreateRequest{
			Name:      "API Test",
			Type:      "api",
			Host:      "api.example.com",
			Port:      443,
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "API Test", datasource.Name)
		assert.Equal(t, "api", datasource.Type)
	})

	t.Run("成功创建CSV数据源", func(t *testing.T) {
		t.Skip("Requires validator mock implementation")
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		req := &CreateRequest{
			Name:      "CSV Test",
			Type:      "csv",
			Host:      "localhost",
			Port:      8080,
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "CSV Test", datasource.Name)
		assert.Equal(t, "csv", datasource.Type)
	})

	t.Run("成功创建Excel数据源", func(t *testing.T) {
		t.Skip("Requires validator mock implementation")
		repo := &mockDatasourceRepo{}
		service := NewService(repo)

		req := &CreateRequest{
			Name:      "Excel Test",
			Type:      "excel",
			Host:      "localhost",
			Port:      8080,
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		datasource, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, datasource)
		assert.Equal(t, "Excel Test", datasource.Name)
		assert.Equal(t, "excel", datasource.Type)
	})
}

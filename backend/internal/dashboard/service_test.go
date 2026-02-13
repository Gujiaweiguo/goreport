package dashboard

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
)

type mockDashboardRepo struct {
	dashboard  *models.Dashboard
	dashboards []*models.Dashboard
	createErr  error
	updateErr  error
	deleteErr  error
	getErr     error
	listErr    error
}

func (m *mockDashboardRepo) Create(dashboard *models.Dashboard) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.dashboard = dashboard
	return nil
}

func (m *mockDashboardRepo) Update(dashboard *models.Dashboard) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.dashboard = dashboard
	return nil
}

func (m *mockDashboardRepo) Delete(id, tenantID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.dashboard = nil
	return nil
}

func (m *mockDashboardRepo) Get(id, tenantID string) (*models.Dashboard, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.dashboard != nil && m.dashboard.ID == id {
		return m.dashboard, nil
	}
	return nil, errors.New("dashboard not found")
}

func (m *mockDashboardRepo) List(tenantID string) ([]*models.Dashboard, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.dashboards, nil
}

func TestNewService(t *testing.T) {
	service := NewService(nil)
	assert.NotNil(t, service)
}

func TestService_Create(t *testing.T) {
	repo := &mockDashboardRepo{}
	service := NewService(repo)

	t.Run("成功创建仪表盘", func(t *testing.T) {
		req := &CreateRequest{
			Name:        "Test Dashboard",
			Code:        "test-dashboard",
			Description: "A test dashboard",
			Config: models.DashboardConfig{
				Width:           1280,
				Height:          720,
				BackgroundColor: "#ffffff",
			},
			Components: []models.DashboardComponent{
				{ID: "comp-1", Type: "chart", Title: "Sales Chart"},
			},
			TenantID:  "tenant-1",
			CreatedBy: "user-1",
		}

		dashboard, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, dashboard)
		assert.Equal(t, "Test Dashboard", dashboard.Name)
		assert.Equal(t, "test-dashboard", dashboard.Code)
		assert.Equal(t, 1280, dashboard.Config.Width)
		assert.Equal(t, 720, dashboard.Config.Height)
		assert.Equal(t, "#ffffff", dashboard.Config.BackgroundColor)
		assert.Len(t, dashboard.Components, 1)
	})

	t.Run("名称不能为空", func(t *testing.T) {
		req := &CreateRequest{
			Name:     "",
			Code:     "test",
			TenantID: "tenant-1",
		}

		_, err := service.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("使用默认配置", func(t *testing.T) {
		req := &CreateRequest{
			Name:     "Default Config Dashboard",
			TenantID: "tenant-1",
		}

		dashboard, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 1920, dashboard.Config.Width)
		assert.Equal(t, 1080, dashboard.Config.Height)
		assert.Equal(t, "#0a0e27", dashboard.Config.BackgroundColor)
	})
}

func TestService_Update(t *testing.T) {
	existingDashboard := &models.Dashboard{
		ID:       "dashboard-1",
		TenantID: "tenant-1",
		Name:     "Old Name",
		Code:     "old-code",
	}

	repo := &mockDashboardRepo{
		dashboard: existingDashboard,
	}
	service := NewService(repo)

	t.Run("成功更新仪表盘", func(t *testing.T) {
		newName := "Updated Name"
		req := &UpdateRequest{
			ID:       "dashboard-1",
			Name:     newName,
			TenantID: "tenant-1",
			Status:   1,
		}

		dashboard, err := service.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, dashboard)
		assert.Equal(t, newName, dashboard.Name)
		assert.Equal(t, 1, dashboard.Status)
	})

	t.Run("更新配置", func(t *testing.T) {
		req := &UpdateRequest{
			ID: "dashboard-1",
			Config: models.DashboardConfig{
				Width:           2560,
				Height:          1440,
				BackgroundColor: "#000000",
			},
			TenantID: "tenant-1",
		}

		dashboard, err := service.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 2560, dashboard.Config.Width)
		assert.Equal(t, 1440, dashboard.Config.Height)
		assert.Equal(t, "#000000", dashboard.Config.BackgroundColor)
	})

	t.Run("更新组件", func(t *testing.T) {
		newComponents := []models.DashboardComponent{
			{ID: "comp-2", Type: "metric", Title: "Total Users"},
		}
		req := &UpdateRequest{
			ID:         "dashboard-1",
			Components: newComponents,
			TenantID:   "tenant-1",
		}

		dashboard, err := service.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.Len(t, dashboard.Components, 1)
		assert.Equal(t, "comp-2", dashboard.Components[0].ID)
	})

	t.Run("仪表盘ID不能为空", func(t *testing.T) {
		req := &UpdateRequest{
			Name:     "Test",
			TenantID: "tenant-1",
		}

		_, err := service.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "id is required")
	})

	t.Run("仪表盘不存在", func(t *testing.T) {
		repo.getErr = errors.New("not found")
		req := &UpdateRequest{
			ID:       "not-exist",
			Name:     "Test",
			TenantID: "tenant-1",
		}

		_, err := service.Update(context.Background(), req)

		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	existingDashboard := &models.Dashboard{
		ID:       "dashboard-1",
		TenantID: "tenant-1",
		Name:     "Test Dashboard",
	}

	repo := &mockDashboardRepo{
		dashboard: existingDashboard,
	}
	service := NewService(repo)

	t.Run("成功删除仪表盘", func(t *testing.T) {
		err := service.Delete(context.Background(), "dashboard-1", "tenant-1")

		assert.NoError(t, err)
	})

	t.Run("删除失败", func(t *testing.T) {
		repo.deleteErr = errors.New("delete failed")

		err := service.Delete(context.Background(), "dashboard-1", "tenant-1")

		assert.Error(t, err)
	})
}

func TestService_Get(t *testing.T) {
	existingDashboard := &models.Dashboard{
		ID:       "dashboard-1",
		TenantID: "tenant-1",
		Name:     "Test Dashboard",
		Code:     "test-code",
	}

	repo := &mockDashboardRepo{
		dashboard: existingDashboard,
	}
	service := NewService(repo)

	t.Run("成功获取仪表盘", func(t *testing.T) {
		dashboard, err := service.Get(context.Background(), "dashboard-1", "tenant-1")

		assert.NoError(t, err)
		assert.NotNil(t, dashboard)
		assert.Equal(t, "dashboard-1", dashboard.ID)
		assert.Equal(t, "Test Dashboard", dashboard.Name)
		assert.Equal(t, "test-code", dashboard.Code)
	})

	t.Run("仪表盘不存在", func(t *testing.T) {
		repo.getErr = errors.New("not found")

		_, err := service.Get(context.Background(), "not-exist", "tenant-1")

		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	dashboards := []*models.Dashboard{
		{
			ID:       "dashboard-1",
			TenantID: "tenant-1",
			Name:     "Dashboard 1",
		},
		{
			ID:       "dashboard-2",
			TenantID: "tenant-1",
			Name:     "Dashboard 2",
		},
	}

	repo := &mockDashboardRepo{
		dashboards: dashboards,
	}
	service := NewService(repo)

	t.Run("成功获取仪表盘列表", func(t *testing.T) {
		list, err := service.List(context.Background(), "tenant-1")

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, "Dashboard 1", list[0].Name)
		assert.Equal(t, "Dashboard 2", list[1].Name)
	})

	t.Run("空列表", func(t *testing.T) {
		repo.dashboards = []*models.Dashboard{}

		list, err := service.List(context.Background(), "tenant-1")

		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("列表获取失败", func(t *testing.T) {
		repo.listErr = errors.New("database error")

		_, err := service.List(context.Background(), "tenant-1")

		assert.Error(t, err)
	})
}

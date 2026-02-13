package chart

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/dataset"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockQueryExecutor struct {
	mock.Mock
}

func (m *mockQueryExecutor) Query(ctx context.Context, req *dataset.QueryRequest) (*dataset.QueryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dataset.QueryResponse), args.Error(1)
}

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, chart *models.Chart) error {
	args := m.Called(ctx, chart)
	return args.Error(0)
}

func (m *mockRepository) Update(ctx context.Context, chart *models.Chart) error {
	args := m.Called(ctx, chart)
	return args.Error(0)
}

func (m *mockRepository) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockRepository) Get(ctx context.Context, id, tenantID string) (*models.Chart, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Chart), args.Error(1)
}

func (m *mockRepository) List(ctx context.Context, tenantID string) ([]*models.Chart, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Chart), args.Error(1)
}

func TestNewService(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)
	assert.NotNil(t, svc)
}

func TestService_Create(t *testing.T) {
	t.Run("成功创建图表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Create", mock.Anything, mock.MatchedBy(func(c *models.Chart) bool {
			return c.Name == "Test Chart"
		})).Return(nil).Once()

		chart, err := svc.Create(context.Background(), &CreateRequest{
			TenantID: "tenant-1",
			Name:     "Test Chart",
			Type:     "bar",
			Config: ChartConfig{
				Title: "Test",
				Series: []SeriesConfig{
					{
						Name: "Test Series",
						Type: "bar",
						Data: []any{1, 2, 3},
					},
				},
			},
		})

		assert.NoError(t, err)
		assert.NotNil(t, chart)
		assert.Equal(t, "Test Chart", chart.Name)
		assert.Equal(t, "tenant-1", chart.TenantID)
		repo.AssertExpectations(t)
	})

	t.Run("创建失败-Repository错误", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()

		_, err := svc.Create(context.Background(), &CreateRequest{
			TenantID: "tenant-1",
			Name:     "Test Chart",
			Type:     "bar",
			Config: ChartConfig{
				Series: []SeriesConfig{{Name: "S1", Type: "bar", Data: []any{1}}},
			},
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
		repo.AssertExpectations(t)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("成功获取图表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		existingChart := &models.Chart{ID: "c1", Name: "Test Chart", TenantID: "tenant-1"}
		repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()

		chart, err := svc.Get(context.Background(), "c1", "tenant-1")
		assert.NoError(t, err)
		assert.Equal(t, "c1", chart.ID)
		assert.Equal(t, "Test Chart", chart.Name)
		repo.AssertExpectations(t)
	})

	t.Run("图表不存在", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Get", mock.Anything, "not-found", "tenant-1").Return(nil, assert.AnError).Once()

		chart, err := svc.Get(context.Background(), "not-found", "tenant-1")
		assert.Nil(t, chart)
		assert.ErrorIs(t, err, ErrNotFound)
		repo.AssertExpectations(t)
	})
}

func TestService_Update(t *testing.T) {
	t.Run("成功更新图表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		existingChart := &models.Chart{ID: "c1", Name: "Old Name", TenantID: "tenant-1"}
		repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()
		repo.On("Update", mock.Anything, mock.MatchedBy(func(c *models.Chart) bool {
			return c.Name == "New Name"
		})).Return(nil).Once()

		chart, err := svc.Update(context.Background(), &UpdateRequest{
			TenantID: "tenant-1",
			ID:       "c1",
			Name:     "New Name",
			Config: ChartConfig{
				Series: []SeriesConfig{{Name: "S1", Type: "bar", Data: []any{1}}},
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, "New Name", chart.Name)
		repo.AssertExpectations(t)
	})

	t.Run("更新失败-图表不存在", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Get", mock.Anything, "not-found", "tenant-1").Return(nil, assert.AnError).Once()

		_, err := svc.Update(context.Background(), &UpdateRequest{
			TenantID: "tenant-1",
			ID:       "not-found",
			Name:     "New Name",
		})

		assert.ErrorIs(t, err, ErrNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("更新失败-Repository错误", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		existingChart := &models.Chart{ID: "c1", Name: "Old Name", TenantID: "tenant-1"}
		repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()
		repo.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed")).Once()

		_, err := svc.Update(context.Background(), &UpdateRequest{
			TenantID: "tenant-1",
			ID:       "c1",
			Name:     "New Name",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
		repo.AssertExpectations(t)
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("成功删除图表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Delete", mock.Anything, "c1", "tenant-1").Return(nil).Once()

		err := svc.Delete(context.Background(), "c1", "tenant-1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("删除失败", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Delete", mock.Anything, "c1", "tenant-1").Return(errors.New("delete failed")).Once()

		err := svc.Delete(context.Background(), "c1", "tenant-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
		repo.AssertExpectations(t)
	})
}

func TestService_List(t *testing.T) {
	t.Run("成功获取图表列表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		expected := []*models.Chart{{ID: "c1", Name: "Chart 1"}, {ID: "c2", Name: "Chart 2"}}
		repo.On("List", mock.Anything, "tenant-1").Return(expected, nil).Once()

		charts, err := svc.List(context.Background(), "tenant-1")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(charts))
		assert.Equal(t, "c1", charts[0].ID)
		assert.Equal(t, "Chart 1", charts[0].Name)
		repo.AssertExpectations(t)
	})

	t.Run("空列表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("List", mock.Anything, "tenant-1").Return([]*models.Chart{}, nil).Once()

		charts, err := svc.List(context.Background(), "tenant-1")
		assert.NoError(t, err)
		assert.Empty(t, charts)
		repo.AssertExpectations(t)
	})

	t.Run("列表获取失败", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("List", mock.Anything, "tenant-1").Return(nil, errors.New("list error")).Once()

		_, err := svc.List(context.Background(), "tenant-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "list error")
		repo.AssertExpectations(t)
	})
}

func TestService_Render(t *testing.T) {
	t.Run("成功渲染图表", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		configJSON, _ := json.Marshal(ChartConfig{
			Series: []SeriesConfig{
				{
					Name:      "Series 1",
					Type:      "bar",
					DatasetID: "ds-1",
					Query:     dataset.QueryRequest{DatasetID: "ds-1"},
				},
			},
		})

		existingChart := &models.Chart{
			ID:       "c1",
			Name:     "Test Chart",
			TenantID: "tenant-1",
			Config:   string(configJSON),
		}

		repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()

		queryExec.On("Query", mock.Anything, mock.MatchedBy(func(req *dataset.QueryRequest) bool {
			return req.DatasetID == "ds-1"
		})).Return(&dataset.QueryResponse{
			Data: []map[string]interface{}{
				{"Series 1": "value1"},
				{"Series 1": "value2"},
			},
		}, nil).Once()

		result, err := svc.Render(context.Background(), "c1", "tenant-1")
		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		repo.AssertExpectations(t)
		queryExec.AssertExpectations(t)
	})

	t.Run("渲染失败-图表不存在", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		repo.On("Get", mock.Anything, "not-found", "tenant-1").Return(nil, assert.AnError).Once()

		_, err := svc.Render(context.Background(), "not-found", "tenant-1")
		assert.ErrorIs(t, err, ErrNotFound)
		repo.AssertExpectations(t)
	})

	t.Run("渲染失败-查询错误", func(t *testing.T) {
		queryExec := &mockQueryExecutor{}
		repo := &mockRepository{}
		svc := NewService(repo, queryExec)

		configJSON, _ := json.Marshal(ChartConfig{
			Series: []SeriesConfig{
				{
					Name:      "Series 1",
					Type:      "bar",
					DatasetID: "ds-1",
				},
			},
		})

		existingChart := &models.Chart{
			ID:       "c1",
			Name:     "Test Chart",
			TenantID: "tenant-1",
			Config:   string(configJSON),
		}

		repo.On("Get", mock.Anything, "c1", "tenant-1").Return(existingChart, nil).Once()
		queryExec.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("query failed")).Once()

		_, err := svc.Render(context.Background(), "c1", "tenant-1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "query failed")

		repo.AssertExpectations(t)
		queryExec.AssertExpectations(t)
	})
}

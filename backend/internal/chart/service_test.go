package chart

import (
	"context"
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
	repo.AssertExpectations(t)
}

func TestService_Get_NotFound(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)

	repo.On("Get", mock.Anything, "not-found", "tenant-1").Return(nil, assert.AnError).Once()

	chart, err := svc.Get(context.Background(), "not-found", "tenant-1")
	assert.Nil(t, chart)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	queryExec := &mockQueryExecutor{}
	repo := &mockRepository{}
	svc := NewService(repo, queryExec)

	expected := []*models.Chart{{ID: "c1", Name: "Chart 1"}}
	repo.On("List", mock.Anything, "tenant-1").Return(expected, nil).Once()

	charts, err := svc.List(context.Background(), "tenant-1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(charts))
	assert.Equal(t, "c1", charts[0].ID)
	repo.AssertExpectations(t)
}

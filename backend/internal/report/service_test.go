package report

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, report *Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *mockRepository) Update(ctx context.Context, report *Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *mockRepository) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockRepository) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Report), args.Error(1)
}

func (m *mockRepository) List(ctx context.Context, tenantID string) ([]*Report, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Report), args.Error(1)
}

func TestNewService(t *testing.T) {
	svc := NewService(nil, nil, nil)
	assert.NotNil(t, svc)
}

func TestService_Create_DefaultType(t *testing.T) {
	repo := &mockRepository{}
	svc := &service{repo: repo}

	req := &CreateRequest{
		TenantID: "test-tenant",
		Name:     "Sales Report",
		Code:     "RPT-001",
		Config:   json.RawMessage(`{"layout":"a4"}`),
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(r *Report) bool {
		return r.TenantID == "test-tenant" && r.Name == "Sales Report" && r.Type == "report" && r.Status == 1
	})).Return(nil).Once()

	created, err := svc.Create(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "report", created.Type)
	assert.Equal(t, 1, created.Status)
	repo.AssertExpectations(t)
}

func TestService_Update_NotFound(t *testing.T) {
	repo := &mockRepository{}
	svc := &service{repo: repo}

	req := &UpdateRequest{TenantID: "test-tenant", ID: "not-exists", Name: "Updated"}

	repo.On("Get", mock.Anything, "not-exists", "test-tenant").Return(nil, errors.New("record not found")).Once()

	updated, err := svc.Update(context.Background(), req)
	assert.Nil(t, updated)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepository{}
	svc := &service{repo: repo}

	repo.On("Get", mock.Anything, "not-exists", "test-tenant").Return(nil, errors.New("record not found")).Once()

	got, err := svc.Get(context.Background(), "not-exists", "test-tenant")
	assert.Nil(t, got)
	assert.ErrorIs(t, err, ErrNotFound)
	repo.AssertExpectations(t)
}

func TestService_List(t *testing.T) {
	repo := &mockRepository{}
	svc := &service{repo: repo}

	expected := []*Report{{ID: "r1", TenantID: "test-tenant", Name: "A"}, {ID: "r2", TenantID: "test-tenant", Name: "B"}}
	repo.On("List", mock.Anything, "test-tenant").Return(expected, nil).Once()

	list, err := svc.List(context.Background(), "test-tenant")
	require.NoError(t, err)
	assert.Len(t, list, 2)
	assert.Equal(t, "r1", list[0].ID)
	repo.AssertExpectations(t)
}

func TestDefaultReportType(t *testing.T) {
	assert.Equal(t, "report", defaultReportType(""))
	assert.Equal(t, "custom", defaultReportType("custom"))
}

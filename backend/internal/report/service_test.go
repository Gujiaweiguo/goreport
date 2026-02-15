package report

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/render"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockReportRepository struct {
	mock.Mock
}

func (m *mockReportRepository) Create(ctx context.Context, report *Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *mockReportRepository) Update(ctx context.Context, report *Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *mockReportRepository) Delete(ctx context.Context, id, tenantID string) error {
	args := m.Called(ctx, id, tenantID)
	return args.Error(0)
}

func (m *mockReportRepository) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	args := m.Called(ctx, id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Report), args.Error(1)
}

func (m *mockReportRepository) List(ctx context.Context, tenantID string) ([]*Report, error) {
	args := m.Called(ctx, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Report), args.Error(1)
}

func TestReportService_Create_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	req := &CreateRequest{
		TenantID: "tenant-1",
		Name:     "Test Report",
		Code:     "RPT001",
		Type:     "report",
		Config:   []byte(`{"cells":[]}`),
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	report, err := svc.Create(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Test Report", report.Name)
	assert.Equal(t, "tenant-1", report.TenantID)
	assert.Equal(t, 1, report.Status)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Create_RepoError(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	req := &CreateRequest{
		TenantID: "tenant-1",
		Name:     "Test Report",
		Config:   []byte(`{}`),
	}

	mockRepo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))

	report, err := svc.Create(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, report)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Update_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	req := &UpdateRequest{
		TenantID: "tenant-1",
		ID:       "r-1",
		Name:     "Updated Report",
		Code:     "RPT002",
		Type:     "dashboard",
		Config:   []byte(`{"cells":[]}`),
	}

	existingReport := &Report{
		ID:       "r-1",
		TenantID: "tenant-1",
		Name:     "Old Report",
		Code:     "RPT001",
		Type:     "report",
	}

	mockRepo.On("Get", mock.Anything, "r-1", "tenant-1").Return(existingReport, nil)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	report, err := svc.Update(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "Updated Report", report.Name)
	assert.Equal(t, "RPT002", report.Code)
	assert.Equal(t, "dashboard", report.Type)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Update_NotFound(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	req := &UpdateRequest{
		TenantID: "tenant-1",
		ID:       "not-exist",
		Name:     "Updated Report",
	}

	mockRepo.On("Get", mock.Anything, "not-exist", "tenant-1").Return(nil, errors.New("not found"))

	report, err := svc.Update(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, report)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Delete_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "r-1", "tenant-1").Return(nil)

	err := svc.Delete(context.Background(), "r-1", "tenant-1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Delete_Error(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	mockRepo.On("Delete", mock.Anything, "r-1", "tenant-1").Return(errors.New("db error"))

	err := svc.Delete(context.Background(), "r-1", "tenant-1")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Get_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	expectedReport := &Report{
		ID:       "r-1",
		TenantID: "tenant-1",
		Name:     "Test Report",
	}

	mockRepo.On("Get", mock.Anything, "r-1", "tenant-1").Return(expectedReport, nil)

	report, err := svc.Get(context.Background(), "r-1", "tenant-1")

	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, "r-1", report.ID)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Get_NotFound(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	mockRepo.On("Get", mock.Anything, "not-exist", "tenant-1").Return(nil, errors.New("not found"))

	report, err := svc.Get(context.Background(), "not-exist", "tenant-1")

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, report)
	mockRepo.AssertExpectations(t)
}

func TestReportService_List_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	expectedReports := []*Report{
		{ID: "r-1", TenantID: "tenant-1", Name: "Report 1"},
		{ID: "r-2", TenantID: "tenant-1", Name: "Report 2"},
	}

	mockRepo.On("List", mock.Anything, "tenant-1").Return(expectedReports, nil)

	reports, err := svc.List(context.Background(), "tenant-1")

	assert.NoError(t, err)
	assert.Len(t, reports, 2)
	mockRepo.AssertExpectations(t)
}

func TestReportService_List_Error(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	mockRepo.On("List", mock.Anything, "tenant-1").Return(nil, errors.New("db error"))

	reports, err := svc.List(context.Background(), "tenant-1")

	assert.Error(t, err)
	assert.Nil(t, reports)
	mockRepo.AssertExpectations(t)
}

func TestDefaultReportType(t *testing.T) {
	assert.Equal(t, "report", defaultReportType(""))
	assert.Equal(t, "dashboard", defaultReportType("dashboard"))
	assert.Equal(t, "custom", defaultReportType("custom"))
}

func TestReportService_Preview_Success(t *testing.T) {
	mockRepo := &mockReportRepository{}
	renderEngine := render.NewEngine(nil, nil)
	svc := NewService(mockRepo, renderEngine, nil)

	existingReport := &Report{
		ID:       "r-1",
		TenantID: "tenant-1",
		Name:     "Test Report",
		Config:   `{"cells":[{"row":0,"col":0,"text":"Hello"}]}`,
	}

	req := &PreviewRequest{
		ID:       "r-1",
		TenantID: "tenant-1",
		Params:   nil,
	}

	mockRepo.On("Get", mock.Anything, "r-1", "tenant-1").Return(existingReport, nil)

	resp, err := svc.Preview(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.HTML, "Hello")
	mockRepo.AssertExpectations(t)
}

func TestReportService_Preview_ReportNotFound(t *testing.T) {
	mockRepo := &mockReportRepository{}
	svc := NewService(mockRepo, nil, nil)

	req := &PreviewRequest{
		ID:       "not-exist",
		TenantID: "tenant-1",
		Params:   nil,
	}

	mockRepo.On("Get", mock.Anything, "not-exist", "tenant-1").Return(nil, errors.New("not found"))

	resp, err := svc.Preview(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

func TestReportService_Preview_InvalidConfig(t *testing.T) {
	mockRepo := &mockReportRepository{}
	renderEngine := render.NewEngine(nil, nil)
	svc := NewService(mockRepo, renderEngine, nil)

	existingReport := &Report{
		ID:       "r-1",
		TenantID: "tenant-1",
		Name:     "Test Report",
		Config:   `{invalid json}`,
	}

	req := &PreviewRequest{
		ID:       "r-1",
		TenantID: "tenant-1",
		Params:   nil,
	}

	mockRepo.On("Get", mock.Anything, "r-1", "tenant-1").Return(existingReport, nil)

	resp, err := svc.Preview(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockRepo.AssertExpectations(t)
}

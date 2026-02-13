package report

import (
	"context"
	"errors"
	"testing"

	"github.com/gujiaweiguo/goreport/internal/render"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	svc := NewService(nil, nil, nil)
	assert.NotNil(t, svc)
}

type mockReportRepo struct {
	report    *Report
	reports   []*Report
	createErr error
	updateErr error
	deleteErr error
	getErr    error
	listErr   error
}

func (m *mockReportRepo) Create(ctx context.Context, report *Report) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.report = report
	return nil
}

func (m *mockReportRepo) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.report != nil && m.report.ID == id {
		return m.report, nil
	}
	return nil, ErrNotFound
}

func (m *mockReportRepo) List(ctx context.Context, tenantID string) ([]*Report, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.reports, nil
}

func (m *mockReportRepo) Update(ctx context.Context, report *Report) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.report = report
	return nil
}

func (m *mockReportRepo) Delete(ctx context.Context, id, tenantID string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.report = nil
	return nil
}

type mockRenderer struct {
	html       string
	previewErr error
}

func (m *mockRenderer) Render(ctx context.Context, configJSON string, params map[string]interface{}, tenantID string) (string, error) {
	if m.previewErr != nil {
		return "", m.previewErr
	}
	return m.html, nil
}

func TestService_Create(t *testing.T) {
	t.Run("成功创建报表", func(t *testing.T) {
		config := `{"layout":"a4"}`

		req := &CreateRequest{
			TenantID: "tenant-1",
			Name:     "Test Report",
			Code:     "TEST-001",
			Type:     "report",
			Config:   []byte(config),
		}

		repo := &mockReportRepo{}
		svc := NewService(repo, nil, nil)

		report, err := svc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, "Test Report", report.Name)
		assert.Equal(t, "TEST-001", report.Code)
		assert.Equal(t, "tenant-1", report.TenantID)
	})

	t.Run("创建失败-Repository错误", func(t *testing.T) {
		config := `{"layout":"a4"}`

		req := &CreateRequest{
			TenantID: "tenant-1",
			Name:     "Test Report",
			Config:   []byte(config),
		}

		repo := &mockReportRepo{createErr: errors.New("db error")}
		svc := NewService(repo, nil, nil)

		_, err := svc.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("创建成功-空配置", func(t *testing.T) {
		req := &CreateRequest{
			TenantID: "tenant-1",
			Name:     "Test",
			Config:   []byte(`{}`),
		}

		repo := &mockReportRepo{}
		svc := NewService(repo, nil, nil)

		report, err := svc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, report)
	})
}

func TestService_Get(t *testing.T) {
	t.Run("成功获取报表", func(t *testing.T) {
		existingReport := &Report{
			ID:       "report-1",
			Name:     "Test Report",
			Code:     "TEST-001",
			Type:     "report",
			Status:   1,
			TenantID: "tenant-1",
		}

		repo := &mockReportRepo{report: existingReport}
		svc := NewService(repo, nil, nil)

		report, err := svc.Get(context.Background(), "report-1", "tenant-1")

		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, "report-1", report.ID)
		assert.Equal(t, "Test Report", report.Name)
	})

	t.Run("报表不存在", func(t *testing.T) {
		repo := &mockReportRepo{getErr: ErrNotFound}
		svc := NewService(repo, nil, nil)

		_, err := svc.Get(context.Background(), "not-exist", "tenant-1")

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestService_List(t *testing.T) {
	t.Run("成功获取报表列表", func(t *testing.T) {
		reports := []*Report{
			{
				ID:       "report-1",
				Name:     "Report 1",
				Type:     "report",
				Status:   1,
				TenantID: "tenant-1",
			},
			{
				ID:       "report-2",
				Name:     "Report 2",
				Type:     "chart",
				Status:   1,
				TenantID: "tenant-1",
			},
		}

		repo := &mockReportRepo{reports: reports}
		svc := NewService(repo, nil, nil)

		list, err := svc.List(context.Background(), "tenant-1")

		assert.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, "Report 1", list[0].Name)
		assert.Equal(t, "Report 2", list[1].Name)
	})

	t.Run("空列表", func(t *testing.T) {
		repo := &mockReportRepo{reports: []*Report{}}
		svc := NewService(repo, nil, nil)

		list, err := svc.List(context.Background(), "tenant-1")

		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})

	t.Run("列表获取失败", func(t *testing.T) {
		repo := &mockReportRepo{listErr: errors.New("database error")}
		svc := NewService(repo, nil, nil)

		_, err := svc.List(context.Background(), "tenant-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
	})
}

func TestService_Update(t *testing.T) {
	t.Run("成功更新报表", func(t *testing.T) {
		existingReport := &Report{
			ID:       "report-1",
			Name:     "Old Name",
			Code:     "OLD-001",
			Type:     "report",
			Status:   0,
			TenantID: "tenant-1",
		}

		newConfig := `{"layout":"b3"}`

		req := &UpdateRequest{
			ID:       "report-1",
			TenantID: "tenant-1",
			Name:     "Updated Name",
			Code:     "NEW-001",
			Type:     "report",
			Config:   []byte(newConfig),
		}

		repo := &mockReportRepo{report: existingReport}
		svc := NewService(repo, nil, nil)

		report, err := svc.Update(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, report)
		assert.Equal(t, "Updated Name", report.Name)
		assert.Equal(t, "NEW-001", report.Code)
	})

	t.Run("更新失败-报表不存在", func(t *testing.T) {
		repo := &mockReportRepo{getErr: ErrNotFound}
		svc := NewService(repo, nil, nil)

		req := &UpdateRequest{
			ID:       "report-1",
			TenantID: "tenant-1",
			Name:     "Test",
		}

		_, err := svc.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("更新失败-Repository错误", func(t *testing.T) {
		existingReport := &Report{
			ID:       "report-1",
			Name:     "Test",
			Type:     "report",
			Status:   0,
			TenantID: "tenant-1",
		}

		repo := &mockReportRepo{
			report:    existingReport,
			updateErr: errors.New("update failed"),
		}
		svc := NewService(repo, nil, nil)

		req := &UpdateRequest{
			ID:       "report-1",
			TenantID: "tenant-1",
			Name:     "Test",
		}

		_, err := svc.Update(context.Background(), req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestService_Delete(t *testing.T) {
	t.Run("成功删除报表", func(t *testing.T) {
		existingReport := &Report{
			ID:       "report-1",
			Name:     "Test Report",
			Type:     "report",
			TenantID: "tenant-1",
		}

		repo := &mockReportRepo{report: existingReport}
		svc := NewService(repo, nil, nil)

		err := svc.Delete(context.Background(), "report-1", "tenant-1")

		assert.NoError(t, err)
	})

	t.Run("删除失败", func(t *testing.T) {
		repo := &mockReportRepo{deleteErr: errors.New("delete failed")}
		svc := NewService(repo, nil, nil)

		err := svc.Delete(context.Background(), "not-exist", "tenant-1")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}

func TestService_Preview(t *testing.T) {
	t.Skip("Preview requires proper render.Engine setup with database")

	t.Run("成功预览报表", func(t *testing.T) {
		existingReport := &Report{
			ID:       "report-1",
			Name:     "Test Report",
			Code:     "TEST-001",
			Type:     "report",
			Status:   1,
			TenantID: "tenant-1",
			Config:   `{"layout":"a4"}`,
		}

		params := map[string]interface{}{
			"param1": "value1",
		}

		engine := &render.Engine{}

		repo := &mockReportRepo{report: existingReport}

		req := &PreviewRequest{
			TenantID: "tenant-1",
			ID:       "report-1",
			Params:   params,
		}

		svc := NewService(repo, engine, nil)

		result, err := svc.Preview(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.HTML)
	})

	t.Run("报表不存在", func(t *testing.T) {
		repo := &mockReportRepo{getErr: ErrNotFound}
		svc := NewService(repo, nil, nil)

		req := &PreviewRequest{
			TenantID: "tenant-1",
			ID:       "not-exist",
			Params:   map[string]interface{}{},
		}

		_, err := svc.Preview(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestDefaultReportType(t *testing.T) {
	assert.Equal(t, "report", defaultReportType(""))
	assert.Equal(t, "custom", defaultReportType("custom"))
}

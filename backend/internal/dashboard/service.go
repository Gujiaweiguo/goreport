package dashboard

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
)

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*models.Dashboard, error)
	Update(ctx context.Context, req *UpdateRequest) (*models.Dashboard, error)
	Delete(ctx context.Context, id, tenantID string) error
	Get(ctx context.Context, id, tenantID string) (*models.Dashboard, error)
	List(ctx context.Context, tenantID string) ([]*models.Dashboard, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type CreateRequest struct {
	Name        string                      `json:"name"`
	Code        string                      `json:"code"`
	Description string                      `json:"description"`
	Config      models.DashboardConfig      `json:"config"`
	Components  []models.DashboardComponent `json:"components"`
	TenantID    string                      `json:"-"`
	CreatedBy   string                      `json:"-"`
}

type UpdateRequest struct {
	ID          string                      `json:"id"`
	Name        string                      `json:"name"`
	Code        string                      `json:"code"`
	Description string                      `json:"description"`
	Config      models.DashboardConfig      `json:"config"`
	Components  []models.DashboardComponent `json:"components"`
	Status      int                         `json:"status"`
	TenantID    string                      `json:"-"`
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*models.Dashboard, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	config := req.Config
	if config.Width == 0 {
		config.Width = 1920
	}
	if config.Height == 0 {
		config.Height = 1080
	}
	if config.BackgroundColor == "" {
		config.BackgroundColor = "#0a0e27"
	}

	dashboard := &models.Dashboard{
		ID:         fmt.Sprintf("dashboard-%d", time.Now().UnixNano()),
		TenantID:   req.TenantID,
		Name:       req.Name,
		Code:       req.Code,
		Config:     config,
		Components: req.Components,
		Status:     1,
		CreatedBy:  req.CreatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (s *service) Update(ctx context.Context, req *UpdateRequest) (*models.Dashboard, error) {
	if req.ID == "" {
		return nil, errors.New("id is required")
	}

	dashboard, err := s.repo.Get(req.ID, req.TenantID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		dashboard.Name = req.Name
	}
	if req.Code != "" {
		dashboard.Code = req.Code
	}
	if req.Config.Width != 0 || req.Config.Height != 0 {
		dashboard.Config = req.Config
	}
	if len(req.Components) > 0 {
		dashboard.Components = req.Components
	}
	if req.Status != 0 {
		dashboard.Status = req.Status
	}
	dashboard.UpdatedAt = time.Now()

	if err := s.repo.Update(dashboard); err != nil {
		return nil, err
	}

	return dashboard, nil
}

func (s *service) Delete(ctx context.Context, id, tenantID string) error {
	return s.repo.Delete(id, tenantID)
}

func (s *service) Get(ctx context.Context, id, tenantID string) (*models.Dashboard, error) {
	return s.repo.Get(id, tenantID)
}

func (s *service) List(ctx context.Context, tenantID string) ([]*models.Dashboard, error) {
	return s.repo.List(tenantID)
}

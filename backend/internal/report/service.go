package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jeecg/jimureport-go/internal/cache"
	"github.com/jeecg/jimureport-go/internal/render"
)

var ErrNotFound = errors.New("report not found")

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*Report, error)
	Update(ctx context.Context, req *UpdateRequest) (*Report, error)
	Delete(ctx context.Context, id, tenantID string) error
	Get(ctx context.Context, id, tenantID string) (*Report, error)
	List(ctx context.Context, tenantID string) ([]*Report, error)
	Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error)
}

type service struct {
	repo   Repository
	render *render.Engine
	cache  *cache.Cache
}

func NewService(repo Repository, engine *render.Engine, cache *cache.Cache) Service {
	return &service{repo: repo, render: engine, cache: cache}
}

type CreateRequest struct {
	TenantID string          `json:"-"`
	Name     string          `json:"name" binding:"required"`
	Code     string          `json:"code"`
	Type     string          `json:"type"`
	Config   json.RawMessage `json:"config" binding:"required"`
}

type UpdateRequest struct {
	TenantID string          `json:"-"`
	ID       string          `json:"id" binding:"required"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Type     string          `json:"type"`
	Config   json.RawMessage `json:"config"`
}

type PreviewRequest struct {
	TenantID string                 `json:"-"`
	ID       string                 `json:"id" binding:"required"`
	Params   map[string]interface{} `json:"params"`
}

type PreviewResponse struct {
	HTML string `json:"html"`
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*Report, error) {
	report := &Report{
		ID:       fmt.Sprintf("report-%d", time.Now().UnixNano()),
		TenantID: req.TenantID,
		Name:     req.Name,
		Code:     req.Code,
		Type:     defaultReportType(req.Type),
		Config:   string(req.Config),
		Status:   1,
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.Invalidate(ctx, req.TenantID, "report:data")
	}

	return report, nil
}

func (s *service) Update(ctx context.Context, req *UpdateRequest) (*Report, error) {
	report, err := s.repo.Get(ctx, req.ID, req.TenantID)
	if err != nil {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		report.Name = req.Name
	}
	if req.Code != "" {
		report.Code = req.Code
	}
	if req.Type != "" {
		report.Type = req.Type
	}
	if len(req.Config) > 0 {
		report.Config = string(req.Config)
	}

	if err := s.repo.Update(ctx, report); err != nil {
		return nil, err
	}

	if s.cache != nil {
		_ = s.cache.Invalidate(ctx, req.TenantID, "report:data")
	}

	return report, nil
}

func (s *service) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, id, tenantID); err != nil {
		return err
	}
	return nil
}

func (s *service) Get(ctx context.Context, id, tenantID string) (*Report, error) {
	report, err := s.repo.Get(ctx, id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	return report, nil
}

func (s *service) List(ctx context.Context, tenantID string) ([]*Report, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *service) Preview(ctx context.Context, req *PreviewRequest) (*PreviewResponse, error) {
	report, err := s.repo.Get(ctx, req.ID, req.TenantID)
	if err != nil {
		return nil, ErrNotFound
	}

	html, err := s.render.Render(ctx, report.Config, req.Params, req.TenantID)
	if err != nil {
		return nil, err
	}

	return &PreviewResponse{HTML: html}, nil
}

func defaultReportType(value string) string {
	if value == "" {
		return "report"
	}
	return value
}

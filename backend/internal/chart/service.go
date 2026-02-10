package chart

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gujiaweiguo/goreport/internal/dataset"
	"github.com/gujiaweiguo/goreport/internal/models"
)

var ErrNotFound = errors.New("chart not found")

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*models.Chart, error)
	Update(ctx context.Context, req *UpdateRequest) (*models.Chart, error)
	Delete(ctx context.Context, id, tenantID string) error
	Get(ctx context.Context, id, tenantID string) (*models.Chart, error)
	List(ctx context.Context, tenantID string) ([]*models.Chart, error)
	Render(ctx context.Context, id, tenantID string) (string, error)
}

type service struct {
	repo          Repository
	queryExecutor dataset.QueryExecutor
}

func NewService(repo Repository, queryExecutor dataset.QueryExecutor) Service {
	return &service{repo: repo, queryExecutor: queryExecutor}
}

type CreateRequest struct {
	TenantID string      `json:"-"`
	Name     string      `json:"name" binding:"required"`
	Code     string      `json:"code"`
	Type     string      `json:"type" binding:"required"`
	Config   ChartConfig `json:"config" binding:"required"`
}

type UpdateRequest struct {
	TenantID string      `json:"-"`
	ID       string      `json:"id" binding:"required"`
	Name     string      `json:"name"`
	Code     string      `json:"code"`
	Type     string      `json:"type"`
	Config   ChartConfig `json:"config"`
}

type ChartConfig struct {
	Title  string                 `json:"title"`
	XAxis  *AxisConfig            `json:"xAxis,omitempty"`
	YAxis  *AxisConfig            `json:"yAxis,omitempty"`
	Series []SeriesConfig         `json:"series"`
	Params map[string]interface{} `json:"params,omitempty"`
}

type AxisConfig struct {
	Type string   `json:"type"`
	Data []string `json:"data"`
	Name string   `json:"name"`
}

type SeriesConfig struct {
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Data      []any                `json:"data"`
	DatasetID string               `json:"datasetId,omitempty"`
	Query     dataset.QueryRequest `json:"query,omitempty"`
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*models.Chart, error) {
	configJSON, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}

	chart := &models.Chart{
		ID:       fmt.Sprintf("chart-%d", time.Now().UnixNano()),
		TenantID: req.TenantID,
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		Config:   string(configJSON),
		Status:   1,
	}

	if err := s.repo.Create(ctx, chart); err != nil {
		return nil, err
	}

	return chart, nil
}

func (s *service) Update(ctx context.Context, req *UpdateRequest) (*models.Chart, error) {
	chart, err := s.repo.Get(ctx, req.ID, req.TenantID)
	if err != nil {
		return nil, ErrNotFound
	}

	if req.Name != "" {
		chart.Name = req.Name
	}
	if req.Code != "" {
		chart.Code = req.Code
	}
	if req.Type != "" {
		chart.Type = req.Type
	}
	if req.Config.Series != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			return nil, err
		}
		chart.Config = string(configJSON)
	}

	if err := s.repo.Update(ctx, chart); err != nil {
		return nil, err
	}

	return chart, nil
}

func (s *service) Delete(ctx context.Context, id, tenantID string) error {
	return s.repo.Delete(ctx, id, tenantID)
}

func (s *service) Get(ctx context.Context, id, tenantID string) (*models.Chart, error) {
	chart, err := s.repo.Get(ctx, id, tenantID)
	if err != nil {
		return nil, ErrNotFound
	}
	return chart, nil
}

func (s *service) List(ctx context.Context, tenantID string) ([]*models.Chart, error) {
	return s.repo.List(ctx, tenantID)
}

func (s *service) Render(ctx context.Context, id, tenantID string) (string, error) {
	chart, err := s.repo.Get(ctx, id, tenantID)
	if err != nil {
		return "", ErrNotFound
	}

	var config ChartConfig
	if err := json.Unmarshal([]byte(chart.Config), &config); err != nil {
		return "", err
	}

	for i, series := range config.Series {
		if series.DatasetID != "" {
			query := series.Query
			query.DatasetID = series.DatasetID
			resp, err := s.queryExecutor.Query(ctx, &query)
			if err != nil {
				return "", err
			}

			if resp.Data != nil {
				data := make([]any, 0, len(resp.Data))
				for _, row := range resp.Data {
					if val, exists := row[series.Name]; exists {
						data = append(data, val)
					}
				}
				config.Series[i].Data = data
			}
		}
	}

	resultJSON, err := json.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}

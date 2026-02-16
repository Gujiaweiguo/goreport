package datasource

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/repository"
)

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*models.DataSource, error)
	GetByID(ctx context.Context, id string) (*models.DataSource, error)
	List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error)
	Update(ctx context.Context, req *UpdateRequest) (*models.DataSource, error)
	Delete(ctx context.Context, id, tenantID string) error
	Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error)
	Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error)
	Move(ctx context.Context, id, tenantID string) error
	Rename(ctx context.Context, id, tenantID string, newName string) (*models.DataSource, error)
}

type CreateRequest struct {
	Name      string `json:"name" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=mysql postgres mongodb excel csv api"`
	Host      string `json:"host" binding:"required"`
	Port      int    `json:"port" binding:"required,min=1,max=65535"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Database  string `json:"database"`
	TenantID  string `json:"tenantId"`
	CreatedBy string `json:"createdBy"`

	Advanced *AdvancedConfig `json:"advanced,omitempty"`
}

type UpdateRequest struct {
	ID       string `json:"id"` // Set by handler from URL param
	Name     string `json:"name" binding:"max=255"`
	Type     string `json:"type" binding:"omitempty,oneof=mysql postgres mongodb excel csv api"`
	Host     string `json:"host" binding:"omitempty,max=255"`
	Port     int    `json:"port" binding:"omitempty,min=0,max=65535"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	TenantID string `json:"tenantId"` // Set by handler from auth token

	Advanced *AdvancedConfig `json:"advanced,omitempty"`
}

type AdvancedConfig struct {
	SSHHost             string `json:"sshHost,omitempty"`
	SSHPort             int    `json:"sshPort,omitempty"`
	SSHUsername         string `json:"sshUsername,omitempty"`
	SSHPassword         string `json:"sshPassword,omitempty"`
	SSHKey              string `json:"sshKey,omitempty"`
	SSHKeyPhrase        string `json:"sshKeyPhrase,omitempty"`
	MaxConnections      int    `json:"maxConnections,omitempty"`
	QueryTimeoutSeconds int    `json:"queryTimeoutSeconds,omitempty"`
}

type service struct {
	dsRepo           repository.DatasourceRepository
	profileValidator *ProfileValidator
}

func NewService(dsRepo repository.DatasourceRepository) Service {
	return &service{
		dsRepo:           dsRepo,
		profileValidator: NewProfileValidator(),
	}
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*models.DataSource, error) {
	config := map[string]interface{}{
		"host":     req.Host,
		"port":     req.Port,
		"database": req.Database,
		"username": req.Username,
		"password": req.Password,
	}

	if req.Advanced != nil {
		config["ssh_host"] = req.Advanced.SSHHost
		config["ssh_port"] = req.Advanced.SSHPort
		config["ssh_username"] = req.Advanced.SSHUsername
		config["ssh_password"] = req.Advanced.SSHPassword
		config["ssh_key"] = req.Advanced.SSHKey
		config["ssh_key_phrase"] = req.Advanced.SSHKeyPhrase
		config["max_connections"] = req.Advanced.MaxConnections
		config["query_timeout_seconds"] = req.Advanced.QueryTimeoutSeconds
	}

	if err := s.profileValidator.Validate(req.Type, config); err != nil {
		return nil, err
	}

	ds := &models.DataSource{
		ID:        fmt.Sprintf("ds-%d", time.Now().UnixNano()),
		Name:      req.Name,
		Type:      req.Type,
		Host:      req.Host,
		Port:      req.Port,
		Username:  req.Username,
		Password:  req.Password,
		Database:  req.Database,
		TenantID:  req.TenantID,
		CreatedBy: req.CreatedBy,
	}

	if req.Advanced != nil {
		ds.SSHHost = req.Advanced.SSHHost
		ds.SSHPort = req.Advanced.SSHPort
		ds.SSHUsername = req.Advanced.SSHUsername
		ds.SSHPassword = req.Advanced.SSHPassword
		ds.SSHKey = req.Advanced.SSHKey
		ds.SSHKeyPhrase = req.Advanced.SSHKeyPhrase
		ds.MaxConnections = req.Advanced.MaxConnections
		ds.QueryTimeoutSeconds = req.Advanced.QueryTimeoutSeconds
	}

	if err := s.dsRepo.Create(ctx, ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*models.DataSource, error) {
	return s.dsRepo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	return s.dsRepo.List(ctx, tenantID, page, pageSize)
}

func (s *service) Update(ctx context.Context, req *UpdateRequest) (*models.DataSource, error) {
	ds, err := s.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if ds.TenantID != req.TenantID {
		return nil, errors.New("datasource not found")
	}

	config := map[string]interface{}{
		"host":     req.Host,
		"port":     req.Port,
		"database": req.Database,
		"username": req.Username,
		"password": req.Password,
	}

	if req.Advanced != nil {
		config["ssh_host"] = req.Advanced.SSHHost
		config["ssh_port"] = req.Advanced.SSHPort
		config["ssh_username"] = req.Advanced.SSHUsername
		config["ssh_password"] = req.Advanced.SSHPassword
		config["ssh_key"] = req.Advanced.SSHKey
		config["ssh_key_phrase"] = req.Advanced.SSHKeyPhrase
		config["max_connections"] = req.Advanced.MaxConnections
		config["query_timeout_seconds"] = req.Advanced.QueryTimeoutSeconds
	}

	if err := s.profileValidator.Validate(ds.Type, config); err != nil {
		return nil, err
	}

	if req.Name != "" {
		ds.Name = req.Name
	}
	if req.Type != "" {
		ds.Type = req.Type
	}
	if req.Host != "" {
		ds.Host = req.Host
	}
	if req.Port > 0 {
		ds.Port = req.Port
	}
	if req.Username != "" {
		ds.Username = req.Username
	}
	if req.Password != "" {
		ds.Password = req.Password
	}
	if req.Database != "" {
		ds.Database = req.Database
	}

	if req.Advanced != nil {
		if req.Advanced.SSHHost != "" {
			ds.SSHHost = req.Advanced.SSHHost
		}
		if req.Advanced.SSHPort > 0 {
			ds.SSHPort = req.Advanced.SSHPort
		}
		if req.Advanced.SSHUsername != "" {
			ds.SSHUsername = req.Advanced.SSHUsername
		}
		if req.Advanced.SSHPassword != "" {
			ds.SSHPassword = req.Advanced.SSHPassword
		}
		if req.Advanced.SSHKey != "" {
			ds.SSHKey = req.Advanced.SSHKey
		}
		if req.Advanced.SSHKeyPhrase != "" {
			ds.SSHKeyPhrase = req.Advanced.SSHKeyPhrase
		}
		if req.Advanced.MaxConnections > 0 {
			ds.MaxConnections = req.Advanced.MaxConnections
		}
		if req.Advanced.QueryTimeoutSeconds > 0 {
			ds.QueryTimeoutSeconds = req.Advanced.QueryTimeoutSeconds
		}
	}

	if err := s.dsRepo.Update(ctx, ds); err != nil {
		return nil, err
	}

	return ds, nil
}

func (s *service) Delete(ctx context.Context, id, tenantID string) error {
	return s.dsRepo.Delete(ctx, id, tenantID)
}

func (s *service) Search(ctx context.Context, tenantID, keyword string, page, pageSize int) ([]*models.DataSource, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	if keyword == "" {
		return s.List(ctx, tenantID, page, pageSize)
	}

	return s.dsRepo.Search(ctx, tenantID, keyword, page, pageSize)
}

func (s *service) Copy(ctx context.Context, id, tenantID string) (*models.DataSource, error) {
	return s.dsRepo.Copy(ctx, id, tenantID)
}

func (s *service) Move(ctx context.Context, id, tenantID string) error {
	return s.dsRepo.Move(ctx, id, tenantID)
}

func (s *service) Rename(ctx context.Context, id, tenantID string, newName string) (*models.DataSource, error) {
	if newName == "" {
		return nil, errors.New("name cannot be empty")
	}
	if len(newName) > 255 {
		return nil, fmt.Errorf("name too long: %d > 255", len(newName))
	}

	ds, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ds.TenantID != tenantID {
		return nil, errors.New("datasource not found")
	}

	ds.Name = newName

	if err := s.dsRepo.Update(ctx, ds); err != nil {
		return nil, err
	}

	return ds, nil
}

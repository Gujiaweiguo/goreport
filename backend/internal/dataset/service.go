package dataset

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/repository"
)

type Service interface {
	Create(ctx context.Context, req *CreateRequest) (*models.Dataset, error)
	Get(ctx context.Context, id, tenantID string) (*models.Dataset, error)
	GetWithFields(ctx context.Context, id, tenantID string) (*models.Dataset, error)
	List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error)
	Update(ctx context.Context, req *UpdateRequest) (*models.Dataset, error)
	Delete(ctx context.Context, id, tenantID string) error
	Preview(ctx context.Context, id, tenantID string) ([]map[string]interface{}, error)
	GetSchema(ctx context.Context, id, tenantID string) (*SchemaResponse, error)

	CreateComputedField(ctx context.Context, req *CreateFieldRequest) (*models.DatasetField, error)
	UpdateField(ctx context.Context, req *UpdateFieldRequest) (*models.DatasetField, error)
	DeleteField(ctx context.Context, fieldID, tenantID string) error
	ListDimensions(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error)
	ListMeasures(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error)
	ListFields(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error)
}

type service struct {
	datasetRepo    repository.DatasetRepository
	fieldRepo      repository.DatasetFieldRepository
	sourceRepo     repository.DatasetSourceRepository
	datasourceRepo repository.DataSourceRepository
	sqlBuilder     SQLExpressionBuilder
	apiBuilder     APIExpressionBuilder
	cache          *ComputedFieldCache
}

func NewService(
	datasetRepo repository.DatasetRepository,
	fieldRepo repository.DatasetFieldRepository,
	sourceRepo repository.DatasetSourceRepository,
	datasourceRepo repository.DataSourceRepository,
) Service {
	return &service{
		datasetRepo:    datasetRepo,
		fieldRepo:      fieldRepo,
		sourceRepo:     sourceRepo,
		datasourceRepo: datasourceRepo,
		sqlBuilder:     NewSQLExpressionBuilder(),
		apiBuilder:     NewAPIExpressionBuilder(),
		cache:          NewComputedFieldCache(),
	}
}

type CreateRequest struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	DatasourceID *string         `json:"datasourceId"`
	Config       json.RawMessage `json:"config"`
	TenantID     string          `json:"-"`
	CreatedBy    string          `json:"-"`
}

type UpdateRequest struct {
	ID       string          `json:"id"`
	Name     *string         `json:"name"`
	Config   json.RawMessage `json:"config"`
	Status   *int            `json:"status"`
	Action   *string         `json:"action"`
	TenantID string          `json:"-"`
}

type SchemaResponse struct {
	Dimensions []*models.DatasetField `json:"dimensions"`
	Measures   []*models.DatasetField `json:"measures"`
	Computed   []*models.DatasetField `json:"computed"`
}

func (s *service) Create(ctx context.Context, req *CreateRequest) (*models.Dataset, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if req.Type == "" {
		return nil, errors.New("type is required")
	}

	if req.Type == "sql" && req.DatasourceID == nil {
		return nil, errors.New("datasourceId is required for SQL datasets")
	}

	if req.DatasourceID != nil {
		datasource, err := s.datasourceRepo.GetByID(ctx, *req.DatasourceID)
		if err != nil {
			return nil, fmt.Errorf("datasource not found: %w", err)
		}
		if datasource.TenantID != req.TenantID {
			return nil, errors.New("datasource does not belong to this tenant")
		}
	}

	dataset := &models.Dataset{
		ID:           fmt.Sprintf("dataset-%d", time.Now().UnixNano()),
		TenantID:     req.TenantID,
		DatasourceID: req.DatasourceID,
		Name:         req.Name,
		Type:         req.Type,
		Config:       string(req.Config),
		Status:       1,
		CreatedBy:    req.CreatedBy,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.datasetRepo.Create(ctx, dataset); err != nil {
		return nil, err
	}

	if err := s.extractFields(ctx, dataset); err != nil {
		// Keep create flow consistent: if field extraction fails, remove the dataset created in this request.
		if rollbackErr := s.datasetRepo.Delete(ctx, dataset.ID); rollbackErr != nil {
			return nil, fmt.Errorf("failed to extract fields: %w (rollback failed: %v)", err, rollbackErr)
		}
		return nil, fmt.Errorf("failed to extract fields: %w", err)
	}

	return s.datasetRepo.GetByIDWithFields(ctx, dataset.ID)
}

func (s *service) Get(ctx context.Context, id, tenantID string) (*models.Dataset, error) {
	dataset, err := s.datasetRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != tenantID {
		return nil, errors.New("dataset not found")
	}

	return dataset, nil
}

func (s *service) GetWithFields(ctx context.Context, id, tenantID string) (*models.Dataset, error) {
	dataset, err := s.datasetRepo.GetByIDWithFields(ctx, id)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != tenantID {
		return nil, errors.New("dataset not found")
	}

	return dataset, nil
}

func (s *service) List(ctx context.Context, tenantID string, page, pageSize int) ([]*models.Dataset, int64, error) {
	return s.datasetRepo.List(ctx, tenantID, page, pageSize)
}

func (s *service) Update(ctx context.Context, req *UpdateRequest) (*models.Dataset, error) {
	if req.ID == "" {
		return nil, errors.New("id is required")
	}

	dataset, err := s.datasetRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != req.TenantID {
		return nil, errors.New("dataset not found")
	}

	if req.Name != nil {
		dataset.Name = *req.Name
	}
	if req.Config != nil {
		oldConfig := dataset.Config
		dataset.Config = string(req.Config)
		if oldConfig != dataset.Config {
			if err := s.extractFields(ctx, dataset); err != nil {
				return nil, fmt.Errorf("failed to re-extract fields: %w", err)
			}
		}
	}
	if req.Status != nil {
		dataset.Status = *req.Status
	}
	dataset.UpdatedAt = time.Now()

	if err := s.datasetRepo.Update(ctx, dataset); err != nil {
		return nil, err
	}

	return s.datasetRepo.GetByIDWithFields(ctx, dataset.ID)
}

func (s *service) Delete(ctx context.Context, id, tenantID string) error {
	dataset, err := s.datasetRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if dataset.TenantID != tenantID {
		return errors.New("dataset not found")
	}

	if err := s.datasetRepo.SoftDelete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *service) Preview(ctx context.Context, id, tenantID string) ([]map[string]interface{}, error) {
	dataset, err := s.GetWithFields(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	if dataset.Type == "sql" && dataset.DatasourceID != nil {
		return s.executeSQLPreview(ctx, dataset)
	}

	return nil, errors.New("preview not implemented for this dataset type")
}

func (s *service) GetSchema(ctx context.Context, id, tenantID string) (*SchemaResponse, error) {
	dataset, err := s.GetWithFields(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}

	var dimensions []*models.DatasetField
	var measures []*models.DatasetField
	var computed []*models.DatasetField

	for _, field := range dataset.Fields {
		fieldCopy := field
		if fieldCopy.IsComputed {
			computed = append(computed, &fieldCopy)
		} else if fieldCopy.Type == "dimension" {
			dimensions = append(dimensions, &fieldCopy)
		} else {
			measures = append(measures, &fieldCopy)
		}
	}

	return &SchemaResponse{
		Dimensions: dimensions,
		Measures:   measures,
		Computed:   computed,
	}, nil
}

func (s *service) extractFields(ctx context.Context, dataset *models.Dataset) error {
	if dataset.Type == "sql" && dataset.DatasourceID != nil {
		return s.extractSQLFields(ctx, dataset)
	}
	return nil
}

func (s *service) extractSQLFields(ctx context.Context, dataset *models.Dataset) error {
	datasource, err := s.datasourceRepo.GetByID(ctx, *dataset.DatasourceID)
	if err != nil {
		return err
	}

	var config struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(dataset.Config), &config); err != nil {
		return err
	}

	db, err := s.getDBConnection(datasource)
	if err != nil {
		return err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT * FROM (%s) AS tmp LIMIT 0", config.Query)
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	if err := s.fieldRepo.DeleteComputedFields(ctx, dataset.ID); err != nil {
		return err
	}

	for i, colName := range columns {
		dataType := mapSQLTypeToDataType(columnTypes[i].DatabaseTypeName())
		field := &models.DatasetField{
			ID:          fmt.Sprintf("field-%d", time.Now().UnixNano()),
			DatasetID:   dataset.ID,
			Name:        colName,
			DisplayName: &colName,
			Type:        inferFieldType(columnTypes[i].DatabaseTypeName()),
			DataType:    dataType,
			IsComputed:  false,
			Config:      "{}",
			SortIndex:   i,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.fieldRepo.Create(ctx, field); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) executeSQLPreview(ctx context.Context, dataset *models.Dataset) ([]map[string]interface{}, error) {
	datasource, err := s.datasourceRepo.GetByID(ctx, *dataset.DatasourceID)
	if err != nil {
		return nil, err
	}

	var config struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(dataset.Config), &config); err != nil {
		return nil, err
	}

	db, err := s.getDBConnection(datasource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf("%s LIMIT 100", config.Query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (s *service) getDBConnection(datasource *models.DataSource) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		datasource.Username,
		datasource.Password,
		datasource.Host,
		datasource.Port,
		datasource.DatabaseName,
	)

	return sql.Open("mysql", dsn)
}

func mapSQLTypeToDataType(sqlType string) string {
	switch sqlType {
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE", "DECIMAL":
		return "number"
	case "DATE", "DATETIME", "TIMESTAMP", "TIME", "YEAR":
		return "date"
	case "BOOLEAN", "TINYINT(1)":
		return "boolean"
	default:
		return "string"
	}
}

func inferFieldType(sqlType string) string {
	switch sqlType {
	case "VARCHAR", "CHAR", "TEXT", "DATE", "DATETIME", "TIMESTAMP":
		return "dimension"
	case "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "FLOAT", "DOUBLE", "DECIMAL":
		return "measure"
	case "BOOLEAN", "TINYINT(1)":
		return "dimension"
	default:
		return "dimension"
	}
}

type CreateFieldRequest struct {
	DatasetID       string  `json:"datasetId"`
	Name            string  `json:"name"`
	DisplayName     *string `json:"displayName"`
	Type            string  `json:"type"`
	DataType        string  `json:"dataType"`
	Expression      *string `json:"expression"`
	IsGroupingField bool    `json:"isGroupingField"`
	GroupingRule    *string `json:"groupingRule"`
	GroupingEnabled *bool   `json:"groupingEnabled"`
	TenantID        string  `json:"-"`
}

type UpdateFieldRequest struct {
	FieldID         string  `json:"fieldId"`
	DisplayName     *string `json:"displayName"`
	Type            *string `json:"type"`
	DataType        *string `json:"dataType"`
	IsSortable      *bool   `json:"isSortable"`
	IsGroupable     *bool   `json:"isGroupable"`
	SortOrder       *string `json:"sortOrder"`
	IsGroupingField *bool   `json:"isGroupingField"`
	GroupingRule    *string `json:"groupingRule"`
	GroupingEnabled *bool   `json:"groupingEnabled"`
	TenantID        string  `json:"-"`
}

func (s *service) CreateComputedField(ctx context.Context, req *CreateFieldRequest) (*models.DatasetField, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}

	if req.Type != "dimension" && req.Type != "measure" {
		return nil, errors.New("type must be 'dimension' or 'measure'")
	}

	if req.IsGroupingField {
		if req.GroupingRule == nil || *req.GroupingRule == "" {
			return nil, errors.New("groupingRule is required for grouping fields")
		}
	} else {
		if req.Expression == nil || *req.Expression == "" {
			return nil, errors.New("expression is required for computed fields")
		}
	}

	dataset, err := s.datasetRepo.GetByID(ctx, req.DatasetID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != req.TenantID {
		return nil, errors.New("dataset not found")
	}

	if !req.IsGroupingField {
		fields, err := s.fieldRepo.List(ctx, req.DatasetID)
		if err != nil {
			return nil, err
		}

		var fieldNames []string
		for _, f := range fields {
			fieldNames = append(fieldNames, f.Name)
		}

		if err := s.sqlBuilder.Validate(*req.Expression, fieldNames); err != nil {
			return nil, fmt.Errorf("invalid expression: %w", err)
		}

		if strings.Contains(*req.Expression, fmt.Sprintf("[%s]", req.Name)) {
			return nil, errors.New("expression cannot reference itself")
		}
	}

	field := &models.DatasetField{
		ID:               fmt.Sprintf("field-%d", time.Now().UnixNano()),
		DatasetID:        req.DatasetID,
		Name:             req.Name,
		DisplayName:      req.DisplayName,
		Type:             req.Type,
		DataType:         req.DataType,
		IsComputed:       !req.IsGroupingField,
		Expression:       req.Expression,
		IsGroupingField:  req.IsGroupingField,
		GroupingRule:     req.GroupingRule,
		GroupingEnabled:  req.GroupingEnabled,
		Config:           "{}",
		IsSortable:       true,
		IsGroupable:      req.Type == "dimension",
		DefaultSortOrder: "none",
		SortIndex:        0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.fieldRepo.Create(ctx, field); err != nil {
		return nil, err
	}

	s.cache.InvalidateField(field.ID)

	return field, nil
}

func (s *service) UpdateField(ctx context.Context, req *UpdateFieldRequest) (*models.DatasetField, error) {
	field, err := s.fieldRepo.GetByID(ctx, req.FieldID)
	if err != nil {
		return nil, err
	}

	dataset, err := s.datasetRepo.GetByID(ctx, field.DatasetID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != req.TenantID {
		return nil, errors.New("field not found")
	}

	if req.DisplayName != nil {
		field.DisplayName = req.DisplayName
	}
	if req.Type != nil {
		field.Type = *req.Type
	}
	if req.DataType != nil {
		field.DataType = *req.DataType
	}
	if req.IsSortable != nil {
		field.IsSortable = *req.IsSortable
	}
	if req.IsGroupable != nil {
		field.IsGroupable = *req.IsGroupable
	}
	if req.SortOrder != nil {
		field.DefaultSortOrder = *req.SortOrder
	}
	if req.IsGroupingField != nil {
		field.IsGroupingField = *req.IsGroupingField
	}
	if req.GroupingRule != nil {
		field.GroupingRule = req.GroupingRule
	}
	if req.GroupingEnabled != nil {
		field.GroupingEnabled = req.GroupingEnabled
	}
	field.UpdatedAt = time.Now()

	if err := s.fieldRepo.Update(ctx, field); err != nil {
		return nil, err
	}

	return field, nil
}

func (s *service) DeleteField(ctx context.Context, fieldID, tenantID string) error {
	field, err := s.fieldRepo.GetByID(ctx, fieldID)
	if err != nil {
		return err
	}

	if !field.IsComputed {
		return errors.New("cannot delete non-computed fields")
	}

	dataset, err := s.datasetRepo.GetByID(ctx, field.DatasetID)
	if err != nil {
		return err
	}

	if dataset.TenantID != tenantID {
		return errors.New("field not found")
	}

	return s.fieldRepo.Delete(ctx, fieldID)
}

func (s *service) ListDimensions(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	dataset, err := s.datasetRepo.GetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != tenantID {
		return nil, errors.New("dataset not found")
	}

	return s.fieldRepo.ListByType(ctx, datasetID, "dimension")
}

func (s *service) ListMeasures(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	dataset, err := s.datasetRepo.GetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != tenantID {
		return nil, errors.New("dataset not found")
	}

	return s.fieldRepo.ListByType(ctx, datasetID, "measure")
}

func (s *service) ListFields(ctx context.Context, datasetID, tenantID string) ([]*models.DatasetField, error) {
	dataset, err := s.datasetRepo.GetByID(ctx, datasetID)
	if err != nil {
		return nil, err
	}

	if dataset.TenantID != tenantID {
		return nil, errors.New("dataset not found")
	}

	return s.fieldRepo.List(ctx, datasetID)
}

func (s *service) ResolveComputedFieldDependencies(ctx context.Context, datasetID, fieldID string) ([]*models.DatasetField, error) {
	var dependencies []*models.DatasetField
	visited := make(map[string]bool)

	if err := s.resolveFieldDependencies(ctx, datasetID, fieldID, &dependencies, visited); err != nil {
		return nil, err
	}

	return dependencies, nil
}

func (s *service) resolveFieldDependencies(ctx context.Context, datasetID, fieldID string, dependencies *[]*models.DatasetField, visited map[string]bool) error {
	if visited[fieldID] {
		return fmt.Errorf("circular dependency detected for field: %s", fieldID)
	}
	visited[fieldID] = true

	field, err := s.fieldRepo.GetByID(ctx, fieldID)
	if err != nil {
		return err
	}

	if field.DatasetID != datasetID {
		return fmt.Errorf("field %s does not belong to dataset %s", fieldID, datasetID)
	}

	if !field.IsComputed {
		return nil
	}

	if field.Expression != nil {
		expression := *field.Expression
		fieldPattern := regexp.MustCompile(`\[(\w+)\]`)
		matches := fieldPattern.FindAllStringSubmatch(expression, -1)

		for _, match := range matches {
			if len(match) > 1 {
				refFieldName := match[1]
				fields, err := s.fieldRepo.List(ctx, datasetID)
				if err != nil {
					return err
				}

				for _, f := range fields {
					if f.Name == refFieldName {
						if f.IsComputed {
							if err := s.resolveFieldDependencies(ctx, datasetID, f.ID, dependencies, visited); err != nil {
								return err
							}
						}
					}
				}
			}
		}
	}

	*dependencies = append(*dependencies, field)
	return nil
}

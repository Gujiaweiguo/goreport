package dataset

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/gujiaweiguo/goreport/internal/repository"
)

type QueryExecutor interface {
	Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error)
}

type QueryRequest struct {
	DatasetID    string                 `json:"datasetId"`
	Fields       []string               `json:"fields"`
	Filters      []Filter               `json:"filters"`
	SortBy       string                 `json:"sortBy"`
	SortOrder    string                 `json:"sortOrder"`
	Page         int                    `json:"page"`
	PageSize     int                    `json:"pageSize"`
	GroupBy      []string               `json:"groupBy"`
	Aggregations map[string]Aggregation `json:"aggregations"`
}

type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type Aggregation struct {
	Function string `json:"function"`
	Field    string `json:"field"`
}

type QueryResponse struct {
	Data          []map[string]interface{} `json:"data"`
	Total         int64                    `json:"total"`
	Page          int                      `json:"page"`
	PageSize      int                      `json:"pageSize"`
	ExecutionTime int64                    `json:"executionTime"`
}

type queryExecutor struct {
	datasetRepo    repository.DatasetRepository
	fieldRepo      repository.DatasetFieldRepository
	datasourceRepo repository.DataSourceRepository
	sqlBuilder     SQLExpressionBuilder
	cache          *ComputedFieldCache
}

func NewQueryExecutor(
	datasetRepo repository.DatasetRepository,
	fieldRepo repository.DatasetFieldRepository,
	datasourceRepo repository.DataSourceRepository,
	sqlBuilder SQLExpressionBuilder,
	cache *ComputedFieldCache,
) QueryExecutor {
	return &queryExecutor{
		datasetRepo:    datasetRepo,
		fieldRepo:      fieldRepo,
		datasourceRepo: datasourceRepo,
		sqlBuilder:     sqlBuilder,
		cache:          cache,
	}
}

func (q *queryExecutor) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	dataset, err := q.datasetRepo.GetByIDWithFields(ctx, req.DatasetID)
	if err != nil {
		return nil, fmt.Errorf("dataset not found: %w", err)
	}

	if dataset.Type == "sql" && dataset.DatasourceID != nil {
		return q.querySQLDataset(ctx, dataset, req)
	}

	return nil, fmt.Errorf("unsupported dataset type: %s", dataset.Type)
}

func (q *queryExecutor) querySQLDataset(ctx context.Context, dataset *models.Dataset, req *QueryRequest) (*QueryResponse, error) {
	datasource, err := q.datasourceRepo.GetByID(ctx, *dataset.DatasourceID)
	if err != nil {
		return nil, fmt.Errorf("datasource not found: %w", err)
	}

	db, err := q.getDBConnection(datasource)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var config struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal([]byte(dataset.Config), &config); err != nil {
		return nil, fmt.Errorf("invalid dataset config: %w", err)
	}

	selectClause := q.buildSelectClause(dataset, req.Fields)
	whereClause := q.buildWhereClause(req.Filters)
	groupByClause := q.buildGroupByClause(req.GroupBy)
	orderByClause := q.buildOrderByClause(req.SortBy, req.SortOrder)
	limitClause := q.buildLimitClause(req.Page, req.PageSize)

	query := fmt.Sprintf("SELECT %s FROM (%s) AS dataset_query %s %s %s %s",
		selectClause, config.Query, whereClause, groupByClause, orderByClause, limitClause)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
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
		data = append(data, row)
	}

	total, err := q.countQueryResults(db, config.Query, whereClause)
	if err != nil {
		return nil, err
	}

	return &QueryResponse{
		Data:          data,
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
		ExecutionTime: 0,
	}, nil
}

func (q *queryExecutor) buildSelectClause(dataset *models.Dataset, requestedFields []string) string {
	if len(requestedFields) == 0 {
		return "*"
	}

	var selectParts []string
	fieldMap := make(map[string]*models.DatasetField)

	for i := range dataset.Fields {
		field := &dataset.Fields[i]
		fieldMap[field.Name] = field
	}

	for _, fieldName := range requestedFields {
		field, exists := fieldMap[fieldName]
		if !exists {
			selectParts = append(selectParts, fmt.Sprintf("`%s`", fieldName))
			continue
		}

		if field.IsComputed && field.Expression != nil {
			if cachedSQL, ok := q.cache.GetSQL(field.ID); ok {
				selectParts = append(selectParts, fmt.Sprintf("%s AS `%s`", cachedSQL, field.Name))
			} else {
				computedSQL, err := q.sqlBuilder.Build(*field.Expression, []string{field.Name})
				if err != nil {
					selectParts = append(selectParts, fmt.Sprintf("NULL AS `%s`", field.Name))
				} else {
					q.cache.SetSQL(field.ID, computedSQL, 3600000000000)
					selectParts = append(selectParts, fmt.Sprintf("%s AS `%s`", computedSQL, field.Name))
				}
			}
		} else {
			selectParts = append(selectParts, fmt.Sprintf("`%s`", field.Name))
		}
	}

	return strings.Join(selectParts, ", ")
}

func (q *queryExecutor) buildWhereClause(filters []Filter) string {
	if len(filters) == 0 {
		return ""
	}

	var conditions []string
	for _, filter := range filters {
		condition := fmt.Sprintf("`%s` %s ?", filter.Field, q.translateOperator(filter.Operator))
		conditions = append(conditions, condition)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
}

func (q *queryExecutor) buildGroupByClause(groupBy []string) string {
	if len(groupBy) == 0 {
		return ""
	}

	var groupParts []string
	for _, field := range groupBy {
		groupParts = append(groupParts, fmt.Sprintf("`%s`", field))
	}

	return fmt.Sprintf("GROUP BY %s", strings.Join(groupParts, ", "))
}

func (q *queryExecutor) buildOrderByClause(sortBy, sortOrder string) string {
	if sortBy == "" {
		return ""
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	return fmt.Sprintf("ORDER BY `%s` %s", sortBy, strings.ToUpper(sortOrder))
}

func (q *queryExecutor) buildLimitClause(page, pageSize int) string {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	offset := (page - 1) * pageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset)
}

func (q *queryExecutor) countQueryResults(db *sql.DB, baseQuery, whereClause string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s) AS dataset_query %s", baseQuery, whereClause)
	var total int64
	err := db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (q *queryExecutor) getDBConnection(datasource *models.DataSource) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		datasource.Username,
		datasource.Password,
		datasource.Host,
		datasource.Port,
		datasource.DatabaseName,
	)

	return sql.Open("mysql", dsn)
}

func (q *queryExecutor) translateOperator(operator string) string {
	switch operator {
	case "eq":
		return "="
	case "neq":
		return "<>"
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lt":
		return "<"
	case "lte":
		return "<="
	case "like":
		return "LIKE"
	case "in":
		return "IN"
	default:
		return "="
	}
}

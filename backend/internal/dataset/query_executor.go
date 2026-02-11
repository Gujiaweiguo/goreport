package dataset

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
	Aggregations  map[string]interface{}   `json:"aggregations,omitempty"`
}

type queryExecutor struct {
	datasetRepo    repository.DatasetRepository
	fieldRepo      repository.DatasetFieldRepository
	datasourceRepo repository.DatasourceRepository
	sqlBuilder     SQLExpressionBuilder
	cache          *ComputedFieldCache
}

func NewQueryExecutor(
	datasetRepo repository.DatasetRepository,
	fieldRepo repository.DatasetFieldRepository,
	datasourceRepo repository.DatasourceRepository,
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
	if err := validateSQLSafety(config.Query); err != nil {
		return nil, fmt.Errorf("query validation failed: %w", err)
	}

	selectedFields := req.Fields
	if len(selectedFields) == 0 && len(req.GroupBy) > 0 {
		selectedFields = req.GroupBy
	}

	selectClause := q.buildSelectClause(dataset, selectedFields)
	aggregationSelectClause, aggregationAliases := q.buildAggregationSelects(req.Aggregations)
	if aggregationSelectClause != "" {
		if len(req.GroupBy) == 0 {
			// Ungrouped aggregation should only project aggregate columns.
			selectClause = aggregationSelectClause
		} else {
			selectClause = fmt.Sprintf("%s, %s", selectClause, aggregationSelectClause)
		}
	}

	whereClause, whereArgs, err := q.buildWhereClause(req.Filters)
	if err != nil {
		return nil, fmt.Errorf("invalid filter condition: %w", err)
	}
	groupByClause := q.buildGroupByClause(req.GroupBy)
	orderByClause := q.buildOrderByClause(req.SortBy, req.SortOrder)
	limitClause, page, pageSize := q.buildLimitClause(req.Page, req.PageSize)

	query := fmt.Sprintf("SELECT %s FROM (%s) AS dataset_query %s %s %s %s",
		selectClause, config.Query, whereClause, groupByClause, orderByClause, limitClause)

	queryCtx, cancel := withDatasetQueryTimeout(ctx)
	defer cancel()

	rows, err := db.QueryContext(queryCtx, query, whereArgs...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("query execution timeout")
		}
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	var aggregations map[string]interface{}

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

		if len(aggregationAliases) > 0 && len(data) == 0 {
			aggregations = make(map[string]interface{})
			for _, alias := range aggregationAliases {
				if val, ok := row[alias]; ok {
					aggregations[alias] = val
				}
			}
		}

		data = append(data, row)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("query execution timeout")
		}
		return nil, err
	}

	total, err := q.countQueryResults(queryCtx, db, config.Query, whereClause, whereArgs)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, errors.New("query count timeout")
		}
		return nil, err
	}

	return &QueryResponse{
		Data:          data,
		Total:         total,
		Page:          page,
		PageSize:      pageSize,
		ExecutionTime: 0,
		Aggregations:  aggregations,
	}, nil
}

func (q *queryExecutor) buildAggregationSelects(aggregations map[string]Aggregation) (string, []string) {
	if len(aggregations) == 0 {
		return "", nil
	}

	aggregationSelects := make([]string, 0, len(aggregations))
	aggregationAliases := make([]string, 0, len(aggregations))
	for alias, agg := range aggregations {
		aggregationSelects = append(aggregationSelects, fmt.Sprintf("%s(%s) AS `%s`", agg.Function, agg.Field, alias))
		aggregationAliases = append(aggregationAliases, alias)
	}

	if len(aggregationSelects) == 0 {
		return "", aggregationAliases
	}

	return strings.Join(aggregationSelects, ", "), aggregationAliases
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

func (q *queryExecutor) buildWhereClause(filters []Filter) (string, []interface{}, error) {
	if len(filters) == 0 {
		return "", nil, nil
	}

	var conditions []string
	var args []interface{}
	for _, filter := range filters {
		op := q.translateOperator(filter.Operator)
		if op == "IN" {
			inValues, ok := normalizeINValues(filter.Value)
			if !ok || len(inValues) == 0 {
				return "", nil, fmt.Errorf("field %s expects non-empty array for IN", filter.Field)
			}

			placeholders := make([]string, len(inValues))
			for i := range inValues {
				placeholders[i] = "?"
			}

			conditions = append(conditions, fmt.Sprintf("`%s` IN (%s)", filter.Field, strings.Join(placeholders, ", ")))
			args = append(args, inValues...)
			continue
		}

		conditions = append(conditions, fmt.Sprintf("`%s` %s ?", filter.Field, op))
		args = append(args, filter.Value)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND ")), args, nil
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

func normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	return page, pageSize
}

func (q *queryExecutor) buildLimitClause(page, pageSize int) (string, int, int) {
	page, pageSize = normalizePagination(page, pageSize)
	offset := (page - 1) * pageSize
	return fmt.Sprintf("LIMIT %d OFFSET %d", pageSize, offset), page, pageSize
}

func (q *queryExecutor) countQueryResults(ctx context.Context, db *sql.DB, baseQuery, whereClause string, whereArgs []interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) AS total FROM (%s) AS dataset_query %s", baseQuery, whereClause)
	var total int64
	err := db.QueryRowContext(ctx, query, whereArgs...).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func normalizeINValues(value interface{}) ([]interface{}, bool) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return nil, false
	}

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, false
	}

	out := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		out = append(out, v.Index(i).Interface())
	}

	return out, true
}

func (q *queryExecutor) getDBConnection(datasource *models.DataSource) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		datasource.Username,
		datasource.Password,
		datasource.Host,
		datasource.Port,
		datasource.Database,
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

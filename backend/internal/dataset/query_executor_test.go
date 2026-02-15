package dataset

import (
	"testing"

	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryExecutor(t *testing.T) {
	executor := NewQueryExecutor(nil, nil, nil, nil, nil)
	assert.NotNil(t, executor)
}

func TestQueryExecutor_BuildAggregationSelects(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name            string
		aggregations    map[string]Aggregation
		expectEmpty     bool
		expectedAliases []string
	}{
		{
			name:            "empty aggregations",
			aggregations:    nil,
			expectEmpty:     true,
			expectedAliases: nil,
		},
		{
			name: "single aggregation",
			aggregations: map[string]Aggregation{
				"total": {Function: "SUM", Field: "amount"},
			},
			expectEmpty:     false,
			expectedAliases: []string{"total"},
		},
		{
			name: "multiple aggregations",
			aggregations: map[string]Aggregation{
				"total":   {Function: "SUM", Field: "amount"},
				"average": {Function: "AVG", Field: "price"},
				"count":   {Function: "COUNT", Field: "id"},
			},
			expectEmpty:     false,
			expectedAliases: []string{"total", "average", "count"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clause, aliases := executor.buildAggregationSelects(tt.aggregations)
			if tt.expectEmpty {
				assert.Empty(t, clause)
			} else {
				assert.NotEmpty(t, clause)
			}
			if tt.expectedAliases != nil {
				assert.ElementsMatch(t, tt.expectedAliases, aliases)
			} else {
				assert.Nil(t, aliases)
			}
		})
	}
}

func TestQueryExecutor_BuildSelectClause(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name           string
		dataset        *models.Dataset
		fields         []string
		expectedResult string
	}{
		{
			name: "empty fields returns star",
			dataset: &models.Dataset{
				Fields: []models.DatasetField{},
			},
			fields:         []string{},
			expectedResult: "*",
		},
		{
			name: "single regular field",
			dataset: &models.Dataset{
				Fields: []models.DatasetField{
					{Name: "id"},
					{Name: "name"},
				},
			},
			fields:         []string{"id"},
			expectedResult: "`id`",
		},
		{
			name: "multiple regular fields",
			dataset: &models.Dataset{
				Fields: []models.DatasetField{
					{Name: "id"},
					{Name: "name"},
				},
			},
			fields:         []string{"id", "name"},
			expectedResult: "`id`, `name`",
		},
		{
			name: "non-existent field",
			dataset: &models.Dataset{
				Fields: []models.DatasetField{
					{Name: "id"},
				},
			},
			fields:         []string{"unknown"},
			expectedResult: "`unknown`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.buildSelectClause(tt.dataset, tt.fields)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestQueryExecutor_BuildWhereClause(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name         string
		filters      []Filter
		expectClause bool
		expectedArg  interface{}
		expectError  bool
	}{
		{
			name:         "empty filters",
			filters:      []Filter{},
			expectClause: false,
		},
		{
			name: "single equals filter",
			filters: []Filter{
				{Field: "status", Operator: "eq", Value: "active"},
			},
			expectClause: true,
			expectedArg:  "active",
		},
		{
			name: "greater than filter",
			filters: []Filter{
				{Field: "amount", Operator: "gt", Value: 100},
			},
			expectClause: true,
			expectedArg:  100,
		},
		{
			name: "IN filter with array",
			filters: []Filter{
				{Field: "id", Operator: "in", Value: []interface{}{1, 2, 3}},
			},
			expectClause: true,
		},
		{
			name: "IN filter with invalid value",
			filters: []Filter{
				{Field: "id", Operator: "in", Value: "not-an-array"},
			},
			expectError: true,
		},
		{
			name: "IN filter with empty array",
			filters: []Filter{
				{Field: "id", Operator: "in", Value: []interface{}{}},
			},
			expectError: true,
		},
		{
			name: "multiple filters",
			filters: []Filter{
				{Field: "status", Operator: "eq", Value: "active"},
				{Field: "amount", Operator: "gte", Value: 50},
			},
			expectClause: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clause, args, err := executor.buildWhereClause(tt.filters)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			if tt.expectClause {
				assert.Contains(t, clause, "WHERE")
				if tt.expectedArg != nil {
					assert.Contains(t, args, tt.expectedArg)
				}
			} else {
				assert.Empty(t, clause)
				assert.Nil(t, args)
			}
		})
	}
}

func TestQueryExecutor_BuildGroupByClause(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name     string
		groupBy  []string
		expected string
	}{
		{
			name:     "empty group by",
			groupBy:  []string{},
			expected: "",
		},
		{
			name:     "single field",
			groupBy:  []string{"category"},
			expected: "GROUP BY `category`",
		},
		{
			name:     "multiple fields",
			groupBy:  []string{"category", "region"},
			expected: "GROUP BY `category`, `region`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.buildGroupByClause(tt.groupBy)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueryExecutor_BuildOrderByClause(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		expected  string
	}{
		{
			name:      "empty sort by",
			sortBy:    "",
			sortOrder: "asc",
			expected:  "",
		},
		{
			name:      "default order",
			sortBy:    "name",
			sortOrder: "",
			expected:  "ORDER BY `name` ASC",
		},
		{
			name:      "explicit asc",
			sortBy:    "name",
			sortOrder: "asc",
			expected:  "ORDER BY `name` ASC",
		},
		{
			name:      "explicit desc",
			sortBy:    "amount",
			sortOrder: "desc",
			expected:  "ORDER BY `amount` DESC",
		},
		{
			name:      "invalid order defaults to asc",
			sortBy:    "id",
			sortOrder: "invalid",
			expected:  "ORDER BY `id` ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := executor.buildOrderByClause(tt.sortBy, tt.sortOrder)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueryExecutor_TranslateOperator(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		operator string
		expected string
	}{
		{"eq", "="},
		{"neq", "<>"},
		{"gt", ">"},
		{"gte", ">="},
		{"lt", "<"},
		{"lte", "<="},
		{"like", "LIKE"},
		{"in", "IN"},
		{"unknown", "="},
		{"", "="},
	}

	for _, tt := range tests {
		t.Run(tt.operator, func(t *testing.T) {
			result := executor.translateOperator(tt.operator)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		expectedPage int
		expectedSize int
	}{
		{"zero values", 0, 0, 1, 10},
		{"negative values", -1, -5, 1, 10},
		{"pageSize over limit", 1, 2000, 1, 1000},
		{"valid values", 5, 20, 5, 20},
		{"large valid pageSize", 1, 1000, 1, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, pageSize := normalizePagination(tt.page, tt.pageSize)
			assert.Equal(t, tt.expectedPage, page)
			assert.Equal(t, tt.expectedSize, pageSize)
		})
	}
}

func TestNormalizeINValues(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected []interface{}
		ok       bool
	}{
		{
			name:     "string slice",
			value:    []string{"a", "b", "c"},
			expected: []interface{}{"a", "b", "c"},
			ok:       true,
		},
		{
			name:     "int slice",
			value:    []int{1, 2, 3},
			expected: []interface{}{1, 2, 3},
			ok:       true,
		},
		{
			name:     "interface slice",
			value:    []interface{}{1, "a", true},
			expected: []interface{}{1, "a", true},
			ok:       true,
		},
		{
			name:     "empty slice",
			value:    []string{},
			expected: []interface{}{},
			ok:       true,
		},
		{
			name:     "string not slice",
			value:    "not a slice",
			expected: nil,
			ok:       false,
		},
		{
			name:     "nil value",
			value:    nil,
			expected: nil,
			ok:       false,
		},
		{
			name:     "int not slice",
			value:    42,
			expected: nil,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := normalizeINValues(tt.value)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestQueryExecutor_BuildLimitClause(t *testing.T) {
	executor := &queryExecutor{}

	tests := []struct {
		name         string
		page         int
		pageSize     int
		expectedPage int
		expectedSize int
		expectClause string
	}{
		{
			name:         "zero values normalized",
			page:         0,
			pageSize:     0,
			expectedPage: 1,
			expectedSize: 10,
			expectClause: "LIMIT 10 OFFSET 0",
		},
		{
			name:         "valid pagination",
			page:         2,
			pageSize:     20,
			expectedPage: 2,
			expectedSize: 20,
			expectClause: "LIMIT 20 OFFSET 20",
		},
		{
			name:         "over limit normalized",
			page:         1,
			pageSize:     2000,
			expectedPage: 1,
			expectedSize: 1000,
			expectClause: "LIMIT 1000 OFFSET 0",
		},
		{
			name:         "third page",
			page:         3,
			pageSize:     50,
			expectedPage: 3,
			expectedSize: 50,
			expectClause: "LIMIT 50 OFFSET 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clause, page, pageSize := executor.buildLimitClause(tt.page, tt.pageSize)
			assert.Equal(t, tt.expectedPage, page)
			assert.Equal(t, tt.expectedSize, pageSize)
			assert.Equal(t, tt.expectClause, clause)
		})
	}
}

func TestQueryExecutor_BuildSelectClause_ComputedField(t *testing.T) {
	cache := NewComputedFieldCache()
	sqlBuilder := NewSQLExpressionBuilder()
	executor := &queryExecutor{
		cache:      cache,
		sqlBuilder: sqlBuilder,
	}

	expr := "price * quantity"
	dataset := &models.Dataset{
		Fields: []models.DatasetField{
			{ID: "f1", Name: "price"},
			{ID: "f2", Name: "quantity"},
			{ID: "f3", Name: "total", IsComputed: true, Expression: &expr},
		},
	}

	result := executor.buildSelectClause(dataset, []string{"total"})
	assert.Contains(t, result, "AS `total`")
}

func TestQueryExecutor_BuildSelectClause_ComputedField_Cached(t *testing.T) {
	cache := NewComputedFieldCache()
	sqlBuilder := NewSQLExpressionBuilder()
	executor := &queryExecutor{
		cache:      cache,
		sqlBuilder: sqlBuilder,
	}

	expr := "price * quantity"
	dataset := &models.Dataset{
		Fields: []models.DatasetField{
			{ID: "f1", Name: "price"},
			{ID: "f2", Name: "quantity"},
			{ID: "f3", Name: "total", IsComputed: true, Expression: &expr},
		},
	}

	cache.SetSQL("f3", "price * quantity", 3600000000000)

	result := executor.buildSelectClause(dataset, []string{"total"})
	assert.Contains(t, result, "AS `total`")
}

func TestQueryExecutor_BuildSelectClause_ComputedField_NoExpression(t *testing.T) {
	cache := NewComputedFieldCache()
	sqlBuilder := NewSQLExpressionBuilder()
	executor := &queryExecutor{
		cache:      cache,
		sqlBuilder: sqlBuilder,
	}

	dataset := &models.Dataset{
		Fields: []models.DatasetField{
			{ID: "f1", Name: "price"},
			{ID: "f2", Name: "quantity"},
			{ID: "f3", Name: "total", IsComputed: true, Expression: nil},
		},
	}

	result := executor.buildSelectClause(dataset, []string{"total"})
	assert.Contains(t, result, "`total`")
}

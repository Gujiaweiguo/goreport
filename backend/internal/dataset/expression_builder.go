package dataset

import (
	"fmt"
	"regexp"
	"strings"
)

type SQLExpressionBuilder interface {
	Build(expression string, fields []string) (string, error)
	Validate(expression string, fields []string) error
	SubstituteFieldReferences(expression string, fieldMapping map[string]string) string
	TranslateFunction(expression string, databaseType string) (string, error)
}

type sqlExpressionBuilder struct {
	functionMap map[string]string
}

func NewSQLExpressionBuilder() SQLExpressionBuilder {
	return &sqlExpressionBuilder{
		functionMap: map[string]string{
			"CONCAT":      "CONCAT",
			"SUBSTRING":   "SUBSTRING",
			"LENGTH":      "LENGTH",
			"UPPER":       "UPPER",
			"LOWER":       "LOWER",
			"TRIM":        "TRIM",
			"DATE_FORMAT": "DATE_FORMAT",
			"DATE_ADD":    "DATE_ADD",
			"DATE_SUB":    "DATE_SUB",
			"DATEDIFF":    "DATEDIFF",
			"NOW":         "NOW",
			"CURDATE":     "CURDATE",
			"CURTIME":     "CURTIME",
			"YEAR":        "YEAR",
			"MONTH":       "MONTH",
			"DAY":         "DAY",
			"HOUR":        "HOUR",
			"MINUTE":      "MINUTE",
			"SECOND":      "SECOND",
			"ROUND":       "ROUND",
			"CEIL":        "CEIL",
			"FLOOR":       "FLOOR",
			"ABS":         "ABS",
			"SUM":         "SUM",
			"AVG":         "AVG",
			"COUNT":       "COUNT",
			"MAX":         "MAX",
			"MIN":         "MIN",
			"IF":          "IF",
			"CASE":        "CASE",
		},
	}
}

func (b *sqlExpressionBuilder) Build(expression string, fields []string) (string, error) {
	translated, err := b.TranslateFunction(expression, "mysql")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s)", translated), nil
}

func (b *sqlExpressionBuilder) Validate(expression string, fields []string) error {
	if expression == "" {
		return fmt.Errorf("expression cannot be empty")
	}

	fieldPattern := regexp.MustCompile(`\[(\w+)\]`)
	matches := fieldPattern.FindAllString(expression, -1)

	for _, match := range matches {
		fieldName := strings.Trim(match, "[]")
		found := false
		for _, field := range fields {
			if field == fieldName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("field reference '%s' not found in dataset fields", match)
		}
	}

	return nil
}

func (b *sqlExpressionBuilder) SubstituteFieldReferences(expression string, fieldMapping map[string]string) string {
	result := expression
	for fieldRef, columnName := range fieldMapping {
		result = strings.ReplaceAll(result, fmt.Sprintf("[%s]", fieldRef), columnName)
	}
	return result
}

func (b *sqlExpressionBuilder) TranslateFunction(expression string, databaseType string) (string, error) {
	result := expression

	for genericFunc, sqlFunc := range b.functionMap {
		pattern := fmt.Sprintf(`\b%s\b`, genericFunc)
		if databaseType == "mysql" {
			result = strings.ReplaceAll(result, pattern, sqlFunc)
		}
	}

	return result, nil
}

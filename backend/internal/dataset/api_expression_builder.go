package dataset

import (
	"fmt"
	"regexp"
	"strings"
)

type APIExpressionBuilder interface {
	Build(expression string, fields []string) (string, error)
	Validate(expression string, fields []string) error
	Evaluate(expression string, row map[string]interface{}) (interface{}, error)
}

type apiExpressionBuilder struct {
	functionMap map[string]func(...interface{}) interface{}
}

func NewAPIExpressionBuilder() APIExpressionBuilder {
	return &apiExpressionBuilder{
		functionMap: map[string]func(...interface{}) interface{}{
			"CONCAT": func(args ...interface{}) interface{} {
				var result strings.Builder
				for _, arg := range args {
					result.WriteString(fmt.Sprintf("%v", arg))
				}
				return result.String()
			},
			"SUBSTRING": func(args ...interface{}) interface{} {
				if len(args) < 2 {
					return nil
				}
				str := fmt.Sprintf("%v", args[0])
				start := 0
				if len(args) >= 2 {
					if s, ok := args[1].(int); ok {
						start = s
					}
				}
				if start < 0 || start >= len(str) {
					return ""
				}
				if len(args) >= 3 {
					length := 0
					if l, ok := args[2].(int); ok {
						length = l
					}
					end := start + length
					if end > len(str) {
						end = len(str)
					}
					return str[start:end]
				}
				return str[start:]
			},
			"LENGTH": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				return len(fmt.Sprintf("%v", args[0]))
			},
			"UPPER": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return ""
				}
				return strings.ToUpper(fmt.Sprintf("%v", args[0]))
			},
			"LOWER": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return ""
				}
				return strings.ToLower(fmt.Sprintf("%v", args[0]))
			},
			"TRIM": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return ""
				}
				return strings.TrimSpace(fmt.Sprintf("%v", args[0]))
			},
			"ROUND": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				return args[0]
			},
			"CEIL": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				if val, ok := args[0].(int); ok {
					if val >= 0 {
						return val
					}
					return val - 1
				}
				return args[0]
			},
			"FLOOR": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				if val, ok := args[0].(int); ok {
					if val <= 0 {
						return val
					}
					return val + 1
				}
				return args[0]
			},
			"ABS": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				return args[0]
			},
			"SUM": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				return args[0]
			},
			"AVG": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return 0
				}
				return args[0]
			},
			"COUNT": func(args ...interface{}) interface{} {
				return 0
			},
			"MAX": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return nil
				}
				return args[0]
			},
			"MIN": func(args ...interface{}) interface{} {
				if len(args) == 0 {
					return nil
				}
				return args[0]
			},
			"IF": func(args ...interface{}) interface{} {
				if len(args) < 2 {
					return nil
				}
				condition := fmt.Sprintf("%v", args[0])
				if condition == "true" || condition == "1" {
					return args[1]
				}
				if len(args) >= 3 {
					return args[2]
				}
				return nil
			},
		},
	}
}

func (b *apiExpressionBuilder) Build(expression string, fields []string) (string, error) {
	if err := b.Validate(expression, fields); err != nil {
		return "", err
	}

	return expression, nil
}

func (b *apiExpressionBuilder) Validate(expression string, fields []string) error {
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

func (b *apiExpressionBuilder) Evaluate(expression string, row map[string]interface{}) (interface{}, error) {
	result := expression

	for fieldName, value := range row {
		result = strings.ReplaceAll(result, fmt.Sprintf("[%s]", fieldName), fmt.Sprintf("%v", value))
	}

	if strings.Contains(result, "[") || strings.Contains(result, "]") {
		return nil, fmt.Errorf("unresolved field references in expression")
	}

	return result, nil
}

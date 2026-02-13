package dataset

import (
	"testing"
)

func TestNewAPIExpressionBuilder(t *testing.T) {
	builder := NewAPIExpressionBuilder()
	if builder == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestAPIExpressionBuilder_Build(t *testing.T) {
	builder := NewAPIExpressionBuilder()

	tests := []struct {
		name       string
		expression string
		fields     []string
		wantErr    bool
	}{
		{
			name:       "simple expression",
			expression: "[amount] * [quantity]",
			fields:     []string{"amount", "quantity"},
			wantErr:    false,
		},
		{
			name:       "empty expression",
			expression: "",
			fields:     []string{"amount"},
			wantErr:    true,
		},
		{
			name:       "invalid field reference",
			expression: "[amount] + [unknown]",
			fields:     []string{"amount"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builder.Build(tt.expression, tt.fields)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == "" {
				t.Error("expected non-empty result")
			}
		})
	}
}

func TestAPIExpressionBuilder_Validate(t *testing.T) {
	builder := NewAPIExpressionBuilder()

	tests := []struct {
		name       string
		expression string
		fields     []string
		wantErr    bool
	}{
		{
			name:       "valid expression",
			expression: "[amount] + [quantity]",
			fields:     []string{"amount", "quantity"},
			wantErr:    false,
		},
		{
			name:       "empty expression",
			expression: "",
			fields:     []string{"amount"},
			wantErr:    true,
		},
		{
			name:       "invalid field reference",
			expression: "[amount] + [unknown]",
			fields:     []string{"amount"},
			wantErr:    true,
		},
		{
			name:       "no field references",
			expression: "1 + 1",
			fields:     []string{},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := builder.Validate(tt.expression, tt.fields)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAPIExpressionBuilder_Evaluate(t *testing.T) {
	builder := NewAPIExpressionBuilder()

	tests := []struct {
		name       string
		expression string
		row        map[string]interface{}
		want       interface{}
		wantErr    bool
	}{
		{
			name:       "simple field substitution",
			expression: "[name]",
			row:        map[string]interface{}{"name": "test"},
			want:       "test",
			wantErr:    false,
		},
		{
			name:       "expression with multiple fields",
			expression: "[first] [last]",
			row:        map[string]interface{}{"first": "John", "last": "Doe"},
			want:       "John Doe",
			wantErr:    false,
		},
		{
			name:       "unresolved field reference",
			expression: "[name] [unknown]",
			row:        map[string]interface{}{"name": "test"},
			want:       nil,
			wantErr:    true,
		},
		{
			name:       "numeric field",
			expression: "[count]",
			row:        map[string]interface{}{"count": 42},
			want:       "42",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builder.Evaluate(tt.expression, tt.row)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result != tt.want {
				t.Errorf("Evaluate() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestAPIExpressionBuilder_Functions(t *testing.T) {
	builder := NewAPIExpressionBuilder()
	fnMap := builder.(*apiExpressionBuilder).functionMap

	t.Run("CONCAT", func(t *testing.T) {
		result := fnMap["CONCAT"]("a", "b", "c")
		if result != "abc" {
			t.Errorf("CONCAT() = %v, want abc", result)
		}
	})

	t.Run("UPPER", func(t *testing.T) {
		result := fnMap["UPPER"]("hello")
		if result != "HELLO" {
			t.Errorf("UPPER() = %v, want HELLO", result)
		}
	})

	t.Run("LOWER", func(t *testing.T) {
		result := fnMap["LOWER"]("HELLO")
		if result != "hello" {
			t.Errorf("LOWER() = %v, want hello", result)
		}
	})

	t.Run("LENGTH", func(t *testing.T) {
		result := fnMap["LENGTH"]("hello")
		if result != 5 {
			t.Errorf("LENGTH() = %v, want 5", result)
		}
	})

	t.Run("TRIM", func(t *testing.T) {
		result := fnMap["TRIM"]("  hello  ")
		if result != "hello" {
			t.Errorf("TRIM() = %v, want hello", result)
		}
	})

	t.Run("IF true condition", func(t *testing.T) {
		result := fnMap["IF"]("true", "yes", "no")
		if result != "yes" {
			t.Errorf("IF() = %v, want yes", result)
		}
	})

	t.Run("IF false condition", func(t *testing.T) {
		result := fnMap["IF"]("false", "yes", "no")
		if result != "no" {
			t.Errorf("IF() = %v, want no", result)
		}
	})

	t.Run("IF with 1 condition", func(t *testing.T) {
		result := fnMap["IF"]("1", "yes", "no")
		if result != "yes" {
			t.Errorf("IF() = %v, want yes", result)
		}
	})
}

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

	t.Run("SUBSTRING with length", func(t *testing.T) {
		result := fnMap["SUBSTRING"]("hello", 0, 3)
		if result != "hel" {
			t.Errorf("SUBSTRING() = %v, want hel", result)
		}
	})

	t.Run("SUBSTRING without length", func(t *testing.T) {
		result := fnMap["SUBSTRING"]("hello", 2)
		if result != "llo" {
			t.Errorf("SUBSTRING() = %v, want llo", result)
		}
	})

	t.Run("SUBSTRING out of range", func(t *testing.T) {
		result := fnMap["SUBSTRING"]("hello", 100, 3)
		if result != "" {
			t.Errorf("SUBSTRING() = %v, want empty", result)
		}
	})

	t.Run("SUBSTRING negative start", func(t *testing.T) {
		result := fnMap["SUBSTRING"]("hello", -1, 3)
		if result != "" {
			t.Errorf("SUBSTRING() = %v, want empty", result)
		}
	})

	t.Run("CEIL positive", func(t *testing.T) {
		result := fnMap["CEIL"](3)
		if result != 3 {
			t.Errorf("CEIL() = %v, want 3", result)
		}
	})

	t.Run("CEIL negative", func(t *testing.T) {
		result := fnMap["CEIL"](-3)
		if result != -4 {
			t.Errorf("CEIL() = %v, want -4", result)
		}
	})

	t.Run("FLOOR positive", func(t *testing.T) {
		result := fnMap["FLOOR"](3)
		if result != 4 {
			t.Errorf("FLOOR() = %v, want 4", result)
		}
	})

	t.Run("FLOOR negative", func(t *testing.T) {
		result := fnMap["FLOOR"](-3)
		if result != -3 {
			t.Errorf("FLOOR() = %v, want -3", result)
		}
	})

	t.Run("ABS", func(t *testing.T) {
		result := fnMap["ABS"](-5)
		if result != -5 {
			t.Errorf("ABS() = %v, want -5", result)
		}
	})

	t.Run("SUM", func(t *testing.T) {
		result := fnMap["SUM"](10)
		if result != 10 {
			t.Errorf("SUM() = %v, want 10", result)
		}
	})

	t.Run("AVG", func(t *testing.T) {
		result := fnMap["AVG"](10)
		if result != 10 {
			t.Errorf("AVG() = %v, want 10", result)
		}
	})

	t.Run("COUNT", func(t *testing.T) {
		result := fnMap["COUNT"]()
		if result != 0 {
			t.Errorf("COUNT() = %v, want 0", result)
		}
	})

	t.Run("MAX", func(t *testing.T) {
		result := fnMap["MAX"](10)
		if result != 10 {
			t.Errorf("MAX() = %v, want 10", result)
		}
	})

	t.Run("MIN", func(t *testing.T) {
		result := fnMap["MIN"](10)
		if result != 10 {
			t.Errorf("MIN() = %v, want 10", result)
		}
	})

	t.Run("ROUND", func(t *testing.T) {
		result := fnMap["ROUND"](3.14)
		if result != 3.14 {
			t.Errorf("ROUND() = %v, want 3.14", result)
		}
	})

	t.Run("functions with no args", func(t *testing.T) {
		if fnMap["LENGTH"]() != 0 {
			t.Error("LENGTH with no args should return 0")
		}
		if fnMap["UPPER"]() != "" {
			t.Error("UPPER with no args should return empty")
		}
		if fnMap["LOWER"]() != "" {
			t.Error("LOWER with no args should return empty")
		}
		if fnMap["TRIM"]() != "" {
			t.Error("TRIM with no args should return empty")
		}
		if fnMap["ROUND"]() != 0 {
			t.Error("ROUND with no args should return 0")
		}
		if fnMap["CEIL"]() != 0 {
			t.Error("CEIL with no args should return 0")
		}
		if fnMap["FLOOR"]() != 0 {
			t.Error("FLOOR with no args should return 0")
		}
		if fnMap["ABS"]() != 0 {
			t.Error("ABS with no args should return 0")
		}
		if fnMap["SUM"]() != 0 {
			t.Error("SUM with no args should return 0")
		}
		if fnMap["AVG"]() != 0 {
			t.Error("AVG with no args should return 0")
		}
		if fnMap["MAX"]() != nil {
			t.Error("MAX with no args should return nil")
		}
		if fnMap["MIN"]() != nil {
			t.Error("MIN with no args should return nil")
		}
	})
}

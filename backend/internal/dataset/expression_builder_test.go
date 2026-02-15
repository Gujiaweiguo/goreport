package dataset

import (
	"testing"
)

func TestNewSQLExpressionBuilder(t *testing.T) {
	builder := NewSQLExpressionBuilder()
	if builder == nil {
		t.Fatal("expected non-nil builder")
	}
}

func TestSQLExpressionBuilder_Build(t *testing.T) {
	builder := NewSQLExpressionBuilder()

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
			name:       "expression with function",
			expression: "SUM([amount])",
			fields:     []string{"amount"},
			wantErr:    false,
		},
		{
			name:       "complex expression",
			expression: "[price] * (1 + [tax_rate])",
			fields:     []string{"price", "tax_rate"},
			wantErr:    false,
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

func TestSQLExpressionBuilder_Validate(t *testing.T) {
	builder := NewSQLExpressionBuilder()

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

func TestSQLExpressionBuilder_SubstituteFieldReferences(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	tests := []struct {
		name         string
		expression   string
		fieldMapping map[string]string
		want         string
	}{
		{
			name:       "simple substitution",
			expression: "[amount] * [quantity]",
			fieldMapping: map[string]string{
				"amount":   "order_amount",
				"quantity": "order_qty",
			},
			want: "order_amount * order_qty",
		},
		{
			name:       "partial substitution",
			expression: "[amount] + [unknown]",
			fieldMapping: map[string]string{
				"amount": "order_amount",
			},
			want: "order_amount + [unknown]",
		},
		{
			name:         "no substitution needed",
			expression:   "1 + 1",
			fieldMapping: map[string]string{},
			want:         "1 + 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := builder.SubstituteFieldReferences(tt.expression, tt.fieldMapping)
			if result != tt.want {
				t.Errorf("SubstituteFieldReferences() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestSQLExpressionBuilder_TranslateFunction(t *testing.T) {
	builder := NewSQLExpressionBuilder()

	tests := []struct {
		name         string
		expression   string
		databaseType string
		wantErr      bool
	}{
		{
			name:         "mysql function",
			expression:   "SUM([amount])",
			databaseType: "mysql",
			wantErr:      false,
		},
		{
			name:         "expression with multiple functions",
			expression:   "CONCAT([first_name], ' ', [last_name])",
			databaseType: "mysql",
			wantErr:      false,
		},
		{
			name:         "date function",
			expression:   "DATE_FORMAT([created_at], '%Y-%m-%d')",
			databaseType: "mysql",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := builder.TranslateFunction(tt.expression, tt.databaseType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("TranslateFunction() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == "" {
				t.Error("expected non-empty result")
			}
		})
	}
}

package dataset

import "testing"

func TestValidateSQLSafety(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "valid select",
			query:   "SELECT id, name FROM orders",
			wantErr: false,
		},
		{
			name:    "reject dangerous statement",
			query:   "SELECT * FROM users; DROP TABLE users",
			wantErr: true,
		},
		{
			name:    "reject multiple statements",
			query:   "SELECT 1; SELECT 2",
			wantErr: true,
		},
		{
			name: "reject too many joins",
			query: "SELECT * FROM t1 JOIN t2 ON t1.id=t2.id JOIN t3 ON t1.id=t3.id JOIN t4 ON t1.id=t4.id " +
				"JOIN t5 ON t1.id=t5.id JOIN t6 ON t1.id=t6.id JOIN t7 ON t1.id=t7.id",
			wantErr: true,
		},
		{
			name:    "reject too many nested selects",
			query:   "SELECT * FROM (SELECT * FROM (SELECT * FROM (SELECT * FROM (SELECT * FROM t1)))) t",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSQLSafety(tt.query)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateSQLSafety() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package dataset

import (
	"strings"
	"testing"
)

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
		{
			name:    "reject INSERT",
			query:   "INSERT INTO users VALUES (1, 'test')",
			wantErr: true,
		},
		{
			name:    "reject UPDATE",
			query:   "UPDATE users SET name = 'hacked'",
			wantErr: true,
		},
		{
			name:    "reject DELETE",
			query:   "DELETE FROM users",
			wantErr: true,
		},
		{
			name:    "reject DROP",
			query:   "DROP TABLE users",
			wantErr: true,
		},
		{
			name:    "reject TRUNCATE",
			query:   "TRUNCATE TABLE users",
			wantErr: true,
		},
		{
			name:    "reject ALTER",
			query:   "ALTER TABLE users ADD COLUMN test VARCHAR(100)",
			wantErr: true,
		},
		{
			name:    "reject CREATE",
			query:   "CREATE TABLE evil (id INT)",
			wantErr: true,
		},
		{
			name:    "reject GRANT",
			query:   "GRANT ALL ON *.* TO 'evil'@'%'",
			wantErr: true,
		},
		{
			name:    "empty query",
			query:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only query",
			query:   "   ",
			wantErr: true,
		},
		{
			name:    "valid query with WHERE clause",
			query:   "SELECT * FROM orders WHERE status = 'active'",
			wantErr: false,
		},
		{
			name:    "valid query with ORDER BY",
			query:   "SELECT * FROM orders ORDER BY created_at DESC",
			wantErr: false,
		},
		{
			name:    "valid query with GROUP BY",
			query:   "SELECT category, COUNT(*) FROM orders GROUP BY category",
			wantErr: false,
		},
		{
			name:    "valid query with valid number of joins",
			query:   "SELECT * FROM orders o JOIN customers c ON o.customer_id = c.id JOIN products p ON o.product_id = p.id",
			wantErr: false,
		},
		{
			name:    "valid query with valid nested selects",
			query:   "SELECT * FROM (SELECT * FROM (SELECT * FROM orders) AS sub1) AS sub2",
			wantErr: false,
		},
		{
			name:    "reject query too long",
			query:   "SELECT * FROM orders WHERE " + strings.Repeat("id = 1 AND ", 3000),
			wantErr: true,
		},
		{
			name:    "reject MERGE statement",
			query:   "MERGE INTO target USING source ON (target.id = source.id)",
			wantErr: true,
		},
		{
			name:    "reject REPLACE statement",
			query:   "REPLACE INTO users VALUES (1, 'test')",
			wantErr: true,
		},
		{
			name:    "reject CALL statement",
			query:   "CALL my_procedure()",
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

func TestContainsMultipleStatements(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{
			name:  "single statement with semicolon",
			query: "SELECT * FROM users;",
			want:  false,
		},
		{
			name:  "single statement without semicolon",
			query: "SELECT * FROM users",
			want:  false,
		},
		{
			name:  "multiple statements",
			query: "SELECT * FROM users; DROP TABLE users",
			want:  true,
		},
		{
			name:  "multiple statements with trailing semicolon",
			query: "SELECT 1; SELECT 2;",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsMultipleStatements(tt.query); got != tt.want {
				t.Errorf("containsMultipleStatements() = %v, want %v", got, tt.want)
			}
		})
	}
}

package dataset

import "testing"

func TestBuildLimitClauseNormalizesPagination(t *testing.T) {
	executor := &queryExecutor{}

	clause, page, pageSize := executor.buildLimitClause(0, 2001)
	if page != 1 {
		t.Fatalf("expected normalized page=1, got %d", page)
	}
	if pageSize != 1000 {
		t.Fatalf("expected normalized pageSize=1000, got %d", pageSize)
	}
	if clause != "LIMIT 1000 OFFSET 0" {
		t.Fatalf("unexpected limit clause: %s", clause)
	}

	clause, page, pageSize = executor.buildLimitClause(3, 25)
	if page != 3 || pageSize != 25 {
		t.Fatalf("expected page=3 pageSize=25, got page=%d pageSize=%d", page, pageSize)
	}
	if clause != "LIMIT 25 OFFSET 50" {
		t.Fatalf("unexpected limit clause: %s", clause)
	}
}

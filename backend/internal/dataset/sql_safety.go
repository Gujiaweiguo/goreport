package dataset

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	disallowedSQLKeywordPattern = regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|TRUNCATE|ALTER|CREATE|GRANT|REVOKE|EXEC|CALL|MERGE|REPLACE)\b`)
	joinPattern                 = regexp.MustCompile(`(?i)\bJOIN\b`)
	nestedSelectPattern         = regexp.MustCompile(`(?i)\(\s*SELECT\b`)
)

const (
	maxQueryLength    = 20000
	maxQueryJoinCount = 5
	maxNestedSelects  = 3
)

func validateSQLSafety(query string) error {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return errors.New("query is required")
	}

	if len(trimmedQuery) > maxQueryLength {
		return fmt.Errorf("query is too long")
	}

	if containsMultipleStatements(trimmedQuery) {
		return errors.New("multiple SQL statements are not allowed")
	}

	if disallowedSQLKeywordPattern.MatchString(trimmedQuery) {
		return errors.New("query contains disallowed SQL operation")
	}

	joinCount := len(joinPattern.FindAllStringIndex(trimmedQuery, -1))
	if joinCount > maxQueryJoinCount {
		return fmt.Errorf("query exceeds max join count")
	}

	nestedSelectCount := len(nestedSelectPattern.FindAllStringIndex(trimmedQuery, -1))
	if nestedSelectCount > maxNestedSelects {
		return fmt.Errorf("query exceeds max nested subquery count")
	}

	return nil
}

func containsMultipleStatements(query string) bool {
	trimmedQuery := strings.TrimSuffix(strings.TrimSpace(query), ";")
	return strings.Contains(trimmedQuery, ";")
}

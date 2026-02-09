package report

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {
	handler := NewHandler(nil)
	assert.NotNil(t, handler)
}

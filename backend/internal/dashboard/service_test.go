package dashboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	service := NewService(nil)
	assert.NotNil(t, service)
}

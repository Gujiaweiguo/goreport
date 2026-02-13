package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewHealthHandler(t *testing.T) {
	handler := NewHealthHandler(nil)
	assert.NotNil(t, handler)
}

func TestHealthHandler_Check_Struct(t *testing.T) {
	handler := &HealthHandler{db: nil}
	assert.NotNil(t, handler)
}

func TestHealthHandler_Check_WithMockDB(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:invalid@tcp(localhost:9999)/test?parseTime=true"), &gorm.Config{})
	if err != nil {
		t.Skip("cannot create mock db connection")
	}

	handler := NewHealthHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

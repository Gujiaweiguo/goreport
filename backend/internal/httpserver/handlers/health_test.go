package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func TestHealthHandler_Check_NilDB(t *testing.T) {
	handler := NewHealthHandler(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

func TestHealthHandler_Check_WithDB(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DB_DSN")
	}
	if dsn == "" {
		t.Skip("TEST_DB_DSN or DB_DSN not set")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("cannot connect to database")
	}

	handler := NewHealthHandler(db)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.Check(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
	assert.Contains(t, w.Body.String(), "database")
}

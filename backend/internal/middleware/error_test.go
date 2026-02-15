package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "success", resp["message"])
}

func TestErrorHandler_WithTimeoutError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(http.ErrHandlerTimeout)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusRequestTimeout, w.Code)
	var resp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "request timeout", resp.Message)
}

func TestErrorHandler_WithGenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.Error(assert.AnError)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp ErrorResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "internal server error", resp.Message)
}

func TestErrorHandler_AlreadyWritten(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(ErrorHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "already written"})
		c.Error(assert.AnError)
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecoveryHandler_NoPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RecoveryHandler())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "no panic"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "no panic", resp["message"])
}

func TestRecoveryHandler_WithPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RecoveryHandler())
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	done := make(chan bool)
	go func() {
		router.ServeHTTP(w, req)
		done <- true
	}()

	select {
	case <-done:
		body := w.Body.String()
		assert.Contains(t, body, "internal server error")
		assert.Contains(t, body, "success")
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out - panic was not recovered")
	}
}

func TestRecoveryHandler_PanicWithStruct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(RecoveryHandler())
	router.GET("/test", func(c *gin.Context) {
		panic(struct{ msg string }{"structured panic"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	done := make(chan bool)
	go func() {
		router.ServeHTTP(w, req)
		done <- true
	}()

	select {
	case <-done:
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var resp ErrorResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.False(t, resp.Success)
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestErrorResponse_Structure(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Message: "test error",
		Code:    400,
	}

	assert.False(t, resp.Success)
	assert.Equal(t, "test error", resp.Message)
	assert.Equal(t, 400, resp.Code)
}

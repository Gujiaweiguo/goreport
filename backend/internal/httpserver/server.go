package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gujiaweiguo/goreport/internal/auth"
	"github.com/gujiaweiguo/goreport/internal/cache"
	"github.com/gujiaweiguo/goreport/internal/config"
	"github.com/gujiaweiguo/goreport/internal/dashboard"
	"github.com/gujiaweiguo/goreport/internal/dataset"
	"github.com/gujiaweiguo/goreport/internal/datasource"
	"github.com/gujiaweiguo/goreport/internal/httpserver/handlers"
	"github.com/gujiaweiguo/goreport/internal/middleware"
	"github.com/gujiaweiguo/goreport/internal/render"
	"github.com/gujiaweiguo/goreport/internal/report"
	"github.com/gujiaweiguo/goreport/internal/repository"
	"gorm.io/gorm"
)

// Server HTTP 服务器
type Server struct {
	Engine *gin.Engine
	Server *http.Server
	Cache  *cache.Cache
}

// NewServer 创建新的 HTTP 服务器
func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {
	auth.InitJWT(&cfg.JWT)

	cache, err := cache.New(cfg.Cache)
	if err != nil {
		return nil, err
	}
	auth.InitBlacklist(cache)

	r := gin.Default()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.RecoveryHandler())
	r.Use(auth.AuthMiddleware())

	// 健康检查
	healthHandler := handlers.NewHealthHandler(db)
	r.GET("/health", healthHandler.Check)

	// 认证路由
	authHandler := handlers.NewAuthHandler(db)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
	}

	// 用户与租户路由
	userHandler := handlers.NewUserHandler(repository.NewUserRepository(db))
	users := r.Group("/api/v1/users")
	{
		users.GET("/me", userHandler.GetMe)
	}

	tenantHandler := handlers.NewTenantHandler(repository.NewTenantRepository(db))
	tenants := r.Group("/api/v1/tenants")
	{
		tenants.GET("", tenantHandler.List)
		tenants.GET("/current", tenantHandler.GetCurrent)
	}

	// 数据源路由（新的 datasource 包）
	datasourceRepo := repository.NewDatasourceRepository(db)
	datasourceService := datasource.NewService(datasourceRepo)
	datasourceHandler := datasource.NewHandlerWithMetadata(datasourceService, datasource.NewCachedMetadataService(cache))
	datasources := r.Group("/api/v1/datasources")
	{
		datasources.GET("", datasourceHandler.List)
		datasources.POST("", datasourceHandler.Create)
		datasources.GET("/:id", datasourceHandler.Get)
		datasources.PUT("/:id", datasourceHandler.Update)
		datasources.DELETE("/:id", datasourceHandler.Delete)
		datasources.GET("/:id/tables", datasourceHandler.GetTables)
		datasources.GET("/:id/tables/:table/fields", datasourceHandler.GetFields)
		datasources.POST("/copy/:id", datasourceHandler.Copy)
		datasources.POST("/move", datasourceHandler.Move)
		datasources.PUT("/:id/rename", datasourceHandler.Rename)
		datasources.GET("/search", datasourceHandler.Search)
		datasources.POST("/test", datasourceHandler.TestConnection)
		datasources.POST("/:id/test", datasourceHandler.TestSavedConnection)
		datasources.GET("/profiles", datasourceHandler.ListProfiles)
	}

	// 缓存指标路由
	cacheHandler := handlers.NewCacheHandler(cache)
	r.GET("/api/v1/cache/metrics", cacheHandler.GetMetrics)

	// 报表路由
	reportRepo := report.NewRepository(db)
	reportEngine := render.NewEngine(db, cache)
	reportService := report.NewService(reportRepo, reportEngine, cache)
	reportHandler := report.NewHandler(reportService)
	reports := r.Group("/api/v1/jmreport")
	{
		reports.GET("/list", reportHandler.List)
		reports.GET("/get", reportHandler.Get)
		reports.POST("/create", reportHandler.Create)
		reports.POST("/update", reportHandler.Update)
		reports.DELETE("/delete", reportHandler.Delete)
		reports.POST("/preview", reportHandler.Preview)
	}

	// 仪表盘路由
	dashboardRepo := dashboard.NewRepository(db)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)
	dashboards := r.Group("/api/v1/dashboard")
	{
		dashboards.GET("/list", dashboardHandler.List)
		dashboards.POST("/create", dashboardHandler.Create)
		dashboards.GET("/:id", dashboardHandler.Get)
		dashboards.PUT("/:id", dashboardHandler.Update)
		dashboards.DELETE("/:id", dashboardHandler.Delete)
	}

	// 数据集路由
	datasetRepo := repository.NewDatasetRepository(db)
	fieldRepo := repository.NewDatasetFieldRepository(db)
	sourceRepo := repository.NewDatasetSourceRepository(db)
	datasetService := dataset.NewService(datasetRepo, fieldRepo, sourceRepo, datasourceRepo)
	queryExecutor := dataset.NewQueryExecutor(datasetRepo, fieldRepo, datasourceRepo, dataset.NewSQLExpressionBuilder(), dataset.NewComputedFieldCache())
	datasetHandler := dataset.NewHandler(datasetService, queryExecutor)

	datasets := r.Group("/api/v1/datasets")
	{
		datasets.GET("", datasetHandler.List)
		datasets.POST("", datasetHandler.Create)
		datasets.GET("/:id", datasetHandler.Get)
		datasets.PUT("/:id", datasetHandler.Update)
		datasets.DELETE("/:id", datasetHandler.Delete)
		datasets.GET("/:id/preview", datasetHandler.Preview)
		datasets.POST("/:id/data", datasetHandler.QueryData)
		datasets.GET("/:id/dimensions", datasetHandler.GetDimensions)
		datasets.GET("/:id/measures", datasetHandler.GetMeasures)
		datasets.GET("/:id/schema", datasetHandler.GetSchema)
		datasets.POST("/:id/fields", datasetHandler.CreateComputedField)
		datasets.PATCH("/:id/fields", datasetHandler.BatchUpdateFields)
		datasets.PUT("/:id/fields/:fieldId", datasetHandler.UpdateField)
		datasets.DELETE("/:id/fields/:fieldId", datasetHandler.DeleteField)
	}

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}

	return &Server{
		Engine: r,
		Server: srv,
		Cache:  cache,
	}, nil
}

// Run 启动 HTTP 服务器
func (s *Server) Run(addr string) error {
	s.Server.Addr = addr
	return s.Server.ListenAndServe()
}

// Shutdown 关闭 HTTP 服务器
func (s *Server) Shutdown(ctx context.Context) error {
	if s.Cache != nil {
		_ = s.Cache.Close()
	}
	return s.Server.Shutdown(ctx)
}

// GetEngine 获取 Gin Engine（用于测试）
func (s *Server) GetEngine() *gin.Engine {
	return s.Engine
}

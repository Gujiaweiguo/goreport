package httpserver

import (
	"context"
	"net/http"

	"github.com/jeecg/jimureport-go/internal/auth"
	"github.com/jeecg/jimureport-go/internal/cache"
	"github.com/jeecg/jimureport-go/internal/config"
	"github.com/jeecg/jimureport-go/internal/dashboard"
	"github.com/jeecg/jimureport-go/internal/httpserver/handlers"
	"github.com/jeecg/jimureport-go/internal/middleware"
	"github.com/jeecg/jimureport-go/internal/render"
	"github.com/jeecg/jimureport-go/internal/report"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
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

	// 数据源路由
	datasourceHandler := handlers.NewDataSourceHandler(db, cache)
	datasources := r.Group("/api/v1/datasource")
	{
		datasources.GET("/list", datasourceHandler.ListDatasources)
		datasources.POST("/create", datasourceHandler.CreateDatasource)
		datasources.POST("/test", datasourceHandler.TestDatasource)
		datasources.GET("/:id/tables", datasourceHandler.GetTables)
		datasources.GET("/:id/tables/:table/fields", datasourceHandler.GetFields)
		datasources.PUT("/:id", datasourceHandler.UpdateDatasource)
		datasources.DELETE("/:id", datasourceHandler.DeleteDatasource)
	}

	// 缓存指标路由
	cacheHandler := handlers.NewCacheHandler(cache)
	r.GET("/api/v1/cache/metrics", cacheHandler.GetMetrics)

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

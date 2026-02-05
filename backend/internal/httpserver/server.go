package httpserver

import (
	"context"
	"net/http"

	"github.com/jeecg/jimureport-go/internal/auth"
	"github.com/jeecg/jimureport-go/internal/config"
	"github.com/jeecg/jimureport-go/internal/httpserver/handlers"
	"github.com/jeecg/jimureport-go/internal/render"
	"github.com/jeecg/jimureport-go/internal/report"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Server HTTP 服务器
type Server struct {
	Engine *gin.Engine
	Server *http.Server
}

// NewServer 创建新的 HTTP 服务器
func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	auth.InitJWT(&cfg.JWT)

	r := gin.Default()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
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
	datasourceHandler := handlers.NewDataSourceHandler(db)
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

	reportRepo := report.NewRepository(db)
	reportEngine := render.NewEngine(db)
	reportService := report.NewService(reportRepo, reportEngine)
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

	srv := &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}

	return &Server{
		Engine: r,
		Server: srv,
	}
}

// Run 启动 HTTP 服务器
func (s *Server) Run(addr string) error {
	s.Server.Addr = addr
	return s.Server.ListenAndServe()
}

// Shutdown 关闭 HTTP 服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// GetEngine 获取 Gin Engine（用于测试）
func (s *Server) GetEngine() *gin.Engine {
	return s.Engine
}

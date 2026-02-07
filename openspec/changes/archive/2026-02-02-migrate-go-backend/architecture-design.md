# Go 后端架构设计

## 1. 概述

本文档定义了 goReport Go 后端的技术栈选择、模块划分和架构设计。

## 2. 技术栈选择

### 2.1 HTTP 框架: Gin

**选择理由**:
- 高性能：基于 httprouter，性能优异
- 丰富的中间件生态
- 广泛的社区支持和文档
- 易于学习和使用

**替代方案考虑**:
- **Echo**: 更轻量，但中间件生态不如 Gin
- **Fiber**: 基于 Fasthttp，性能更好，但兼容性稍差
- **Chi**: 轻量级，但功能较少
- **标准库**: 过于底层，开发效率低

### 2.2 数据库访问: GORM

**选择理由**:
- 功能完整的 ORM
- 自动迁移支持
- 关联关系处理强大
- 支持多种数据库（MySQL, PostgreSQL, SQLite, etc.）

**配置**:
```go
import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### 2.3 JWT 验证: golang-jwt/jwt

**选择理由**:
- 官方维护的 JWT 库
- 支持多种签名算法
- 灵活的 claim 处理

### 2.4 配置管理: Viper

**选择理由**:
- 支持多种配置格式（YAML, JSON, TOML, etc.）
- 支持环境变量覆盖
- 支持配置热重载
- 与 Cobra CLI 框架集成良好

### 2.5 日志: Zerolog

**选择理由**:
- 高性能
- 结构化日志
- 支持 JSON 输出
- 易于集成日志聚合服务

### 2.6 依赖注入: Wire (可选)

**选择理由**:
- 编译时依赖注入
- 无反射开销
- 代码清晰
- 易于测试

**替代方案**:
- **Fx**: 运行时 DI，功能更强大
- **手动注入**: 最简单，适合小型项目

### 2.7 验证: Go-playground/validator

**选择理由**:
- 功能完整
- 支持自定义验证规则
- 良好的错误消息

### 2.8 文档生成: Swag (可选)

**选择理由**:
- 基于 Swagger/OpenAPI
- 注解驱动
- 自动生成文档

## 3. 中间件栈

### 3.1 中间件执行顺序

```
1. Recovery (Panic 恢复)
2. Logger (请求日志)
3. CORS (跨域处理)
4. Gzip (响应压缩)
5. RequestID (请求 ID 生成)
6. Timeout (请求超时)
7. Auth (JWT 认证)
8. Tenant (租户识别)
9. RateLimit (速率限制，可选)
10. Router (路由处理)
```

### 3.2 中间件实现

#### 3.2.1 Recovery 中间件

```go
func RecoveryMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                log.Error().Err(err.(error)).Msg("Panic recovered")
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal Server Error",
                })
                c.Abort()
            }
        }()
        c.Next()
    }
}
```

#### 3.2.2 Logger 中间件

```go
func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        status := c.Writer.Status()

        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Str("query", query).
            Int("status", status).
            Dur("latency", latency).
            Str("ip", c.ClientIP()).
            Str("user-agent", c.Request.UserAgent()).
            Msg("Request completed")
    }
}
```

#### 3.2.3 CORS 中间件

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Access-Token, X-Tenant-Id")
        c.Header("Access-Control-Expose-Headers", "Content-Length")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "86400")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}
```

#### 3.2.4 Gzip 中间件

使用官方 `gin-gzip` 包:
```go
import "github.com/gin-contrib/gzip"

router.Use(gzip.Gzip(gzip.DefaultCompression))
```

#### 3.2.5 RequestID 中间件

```go
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        c.Set("request_id", requestID)
        c.Header("X-Request-ID", requestID)
        c.Next()
    }
}
```

#### 3.2.6 Auth 中间件

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := GetToken(c.Request)
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Missing token",
            })
            c.Abort()
            return
        }

        claims, err := ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            c.Abort()
            return
        }

        c.Set("user", claims)
        c.Next()
    }
}
```

#### 3.2.7 Tenant 中间件

```go
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := GetTenantID(c.Request)
        c.Set("tenant_id", tenantID)
        c.Next()
    }
}
```

## 4. 模块划分

### 4.1 模块架构图

```
┌─────────────────────────────────────────────────────────┐
│                      HTTP Layer                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ /jmreport │  │  /drag   │  │  /login  │  │
│  └──────────┘  └──────────┘  └──────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                    Middleware Layer                 │
│  ┌──────┐┌─────┐┌─────┐┌─────┐┌─────┐  │
│  │Recovery││Logger││CORS ││Auth ││Tenant│  │
│  └──────┘└─────┘└─────┘└─────┘└─────┘  │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                     Service Layer                    │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  Report  │  │Dashboard │  │  Auth    │  │
│  │  Service │  │  Service │  │  Service │  │
│  └──────────┘  └──────────┘  └──────────┘  │
│  ┌──────────┐  ┌──────────┐                  │
│  │DataSource │  │  Export  │                  │
│  │  Service │  │  Service │                  │
│  └──────────┘  └──────────┘                  │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                   Repository Layer                   │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  Report  │  │Dashboard │  │  User    │  │
│  │  Repo    │  │   Repo   │  │   Repo   │  │
│  └──────────┘  └──────────┘  └──────────┘  │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│                      Database                       │
│                   (MySQL 5.7+)                      │
└─────────────────────────────────────────────────────────┘
```

### 4.2 模块职责

#### 4.2.1 Auth 模块

**职责**:
- JWT token 生成和验证
- 用户认证
- 角色和权限检查
- 租户识别

**主要接口**:
```go
type AuthService interface {
    GenerateToken(user *User) (string, error)
    ValidateToken(token string) (*Claims, error)
    GetCurrentUser(c *gin.Context) (*User, error)
    HasRole(role string) bool
    HasPermission(permission string) bool
}
```

#### 4.2.2 Report 模块

**职责**:
- 报表模板 CRUD
- 报表渲染
- 报表预览
- 报表分享

**主要接口**:
```go
type ReportService interface {
    Create(report *Report) error
    Update(report *Report) error
    Delete(id string) error
    Get(id string) (*Report, error)
    List(params *ListParams) (*ReportList, error)
    Render(id string, params *RenderParams) (*RenderedReport, error)
}
```

#### 4.2.3 Dashboard 模块

**职责**:
- 仪表盘页面 CRUD
- 组件管理
- 数据集管理
- 仪表盘分享

**主要接口**:
```go
type DashboardService interface {
    CreatePage(page *Page) error
    UpdatePage(page *Page) error
    DeletePage(id string) error
    GetPage(id string) (*Page, error)
    ListPages(params *ListParams) (*PageList, error)
    CreateDataset(dataset *Dataset) error
    UpdateDataset(dataset *Dataset) error
    DeleteDataset(id string) error
}
```

#### 4.2.4 DataSource 模块

**职责**:
- 数据源 CRUD
- 数据源连接测试
- 数据库连接池管理

**主要接口**:
```go
type DataSourceService interface {
    Create(ds *DataSource) error
    Update(ds *DataSource) error
    Delete(id string) error
    Get(id string) (*DataSource, error)
    List(params *ListParams) (*DataSourceList, error)
    TestConnection(ds *DataSource) error
}
```

#### 4.2.5 Export 模块

**职责**:
- 报表导出（Excel, PDF, Word, Image）
- 导出任务管理
- 导出日志记录

**主要接口**:
```go
type ExportService interface {
    Export(reportID string, format string, params *ExportParams) ([]byte, error)
    CreateExportJob(job *ExportJob) error
    GetExportJob(id string) (*ExportJob, error)
}
```

## 5. 项目结构

### 5.1 目录结构

```
jimureport-go/
├── cmd/
│   └── server/
│       └── main.go           # 应用入口
├── internal/
│   ├── auth/                # 认证模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── jwt.go
│   ├── report/              # 报表模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── models.go
│   ├── dashboard/           # 仪表盘模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── models.go
│   ├── datasource/          # 数据源模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── models.go
│   ├── export/              # 导出模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── models.go
│   ├── middleware/          # 中间件
│   │   ├── recovery.go
│   │   ├── logger.go
│   │   ├── cors.go
│   │   ├── auth.go
│   │   └── tenant.go
│   ├── models/              # 共享模型
│   │   ├── user.go
│   │   └── base.go
│   └── config/             # 配置
│       └── config.go
├── pkg/                   # 可复用包
│   ├── database/
│   │   └── gorm.go
│   ├── logger/
│   │   └── zerolog.go
│   └── utils/
│       ├── response.go
│       └── validator.go
├── static/                 # 静态资源
│   ├── login/
│   │   └── login.html
│   └── favicon.ico
├── configs/
│   ├── config.yaml
│   └── config.dev.yaml
├── go.mod
├── go.sum
└── Makefile
```

### 5.2 包依赖关系

```
cmd/server
    └── internal/
        ├── auth
        ├── report
        ├── dashboard
        ├── datasource
        ├── export
        └── middleware
    └── pkg/
        ├── database
        ├── logger
        └── utils
```

## 6. 配置管理

### 6.1 配置文件 (config.yaml)

```yaml
# Server
server:
  port: 8080
  mode: release  # debug, release

# Database
database:
  host: 127.0.0.1
  port: 3306
  name: jimureport
  username: root
  password: root
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

# JWT
jwt:
  secret: your-secret-key
  expiration: 2592000  # 30 days in seconds

# Logging
logging:
  level: info
  format: json

# CORS
cors:
  allowed_origins:
    - "*"
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
  allowed_headers:
    - Content-Type
    - Authorization
    - X-Access-Token
    - X-Tenant-Id
```

### 6.2 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| SERVER_PORT | 服务端口 | 8080 |
| DB_HOST | 数据库主机 | 127.0.0.1 |
| DB_PORT | 数据库端口 | 3306 |
| DB_NAME | 数据库名 | jimureport |
| DB_USERNAME | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | root |
| JWT_SECRET | JWT 密钥 | - |
| LOG_LEVEL | 日志级别 | info |

## 7. 路由设计

### 7.1 路由分组

```go
func SetupRoutes(router *gin.Engine) {
    // API v1
    v1 := router.Group("/api/v1")
    {
        // Auth routes (public)
        auth := v1.Group("/auth")
        {
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
        }

        // Report routes (protected)
        reports := v1.Group("/jmreport")
        reports.Use(AuthMiddleware())
        reports.Use(TenantMiddleware())
        {
            reports.GET("/list", reportHandler.List)
            reports.POST("/create", reportHandler.Create)
            reports.POST("/update", reportHandler.Update)
            reports.POST("/delete", reportHandler.Delete)
            reports.GET("/view/:id", reportHandler.View)
            reports.POST("/render", reportHandler.Render)
        }

        // Dashboard routes (protected)
        dashboards := v1.Group("/drag")
        dashboards.Use(AuthMiddleware())
        dashboards.Use(TenantMiddleware())
        {
            dashboards.GET("/list", dashboardHandler.ListPages)
            dashboards.POST("/page/create", dashboardHandler.CreatePage)
            dashboards.POST("/page/update", dashboardHandler.UpdatePage)
            dashboards.POST("/page/delete", dashboardHandler.DeletePage)
            dashboards.GET("/page/view/:id", dashboardHandler.ViewPage)
            dashboards.GET("/datasets", dashboardHandler.ListDatasets)
        }

        // DataSource routes (protected)
        datasources := v1.Group("/datasource")
        datasources.Use(AuthMiddleware())
        datasources.Use(TenantMiddleware())
        {
            datasources.GET("/list", datasourceHandler.List)
            datasources.POST("/create", datasourceHandler.Create)
            datasources.POST("/update", datasourceHandler.Update)
            datasources.POST("/delete", datasourceHandler.Delete)
            datasources.POST("/test", datasourceHandler.TestConnection)
        }

        // Export routes (protected)
        exports := v1.Group("/export")
        exports.Use(AuthMiddleware())
        {
            exports.POST("/excel", exportHandler.ExportExcel)
            exports.POST("/pdf", exportHandler.ExportPDF)
            exports.POST("/word", exportHandler.ExportWord)
        }
    }

    // Static files
    router.Static("/static", "./static")
    router.LoadHTMLGlob("templates/*")
    router.GET("/login", func(c *gin.Context) {
        c.HTML(http.StatusOK, "login.html", nil)
    })
}
```

## 8. 数据库连接池配置

```go
sqlDB, err := db.DB()
if err != nil {
    log.Fatal().Err(err).Msg("Failed to get database instance")
}

// Set connection pool settings
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## 9. 错误处理

### 9.1 统一错误响应

```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Message string `json:"message,omitempty"`
}

func ErrorResponse(c *gin.Context, statusCode int, error string) {
    c.JSON(statusCode, ErrorResponse{
        Error: error,
        Code:  statusCode,
    })
}
```

### 9.2 错误码定义

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 10. 测试策略

### 10.1 单元测试

使用 `testify` 框架:

```go
func TestAuthService_GenerateToken(t *testing.T) {
    service := NewAuthService(mockDB, mockConfig)
    token, err := service.GenerateToken(testUser)

    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

### 10.2 集成测试

使用 `httptest` 测试 HTTP 端点:

```go
func TestReportHandler_Create(t *testing.T) {
    router := SetupTestRouter()

    req, _ := http.NewRequest("POST", "/api/v1/jmreport/create", testBody)
    req.Header.Set("Authorization", "Bearer "+testToken)

    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)
}
```

## 11. 性能优化

### 11.1 数据库查询优化

- 使用索引
- 避免 N+1 查询
- 使用预加载 (Preload)
- 分页查询

### 11.2 缓存策略

- 使用 Redis 缓存热点数据
- 实现 cache-aside 模式
- 设置合理的缓存过期时间

### 11.3 并发控制

- 使用 goroutine pool
- 限制并发数量
- 实现超时控制

## 12. 监控和追踪

### 12.1 健康检查

```go
router.GET("/health", func(c *gin.Context) {
    db, _ := db.DB()
    if err := db.Ping(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
    })
})
```

### 12.2 指标收集

使用 Prometheus:

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    httpDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request latency",
        },
        []string{"method", "path", "status"},
    )
)
```

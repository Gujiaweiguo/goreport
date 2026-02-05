# JimuReport 系统架构设计

> **版本**：v1.0
> **日期**：2026-02-03
> **状态**：待确认

---

## 1. 整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                    客户端浏览器                              │
│                                                            │
│  ┌────────────────────────────────────────────────────┐         │
│  │        Vue 3 前端应用 (3000)           │         │
│  │   - Report Designer UI                       │         │
│  │   - Report Preview UI                        │         │
│  │   - BI Dashboard UI                         │         │
│  │   - Chart Editor UI                         │         │
│  └────────────┬───────────────────────────────────┘         │
│               │                                            │
└───────────────┼────────────────────────────────────────────┘
                │ HTTP/HTTPS
                │ JWT Token (Header/Query)
                ▼
┌─────────────────────────────────────────────────────────────────┐
│              Go 后端服务 (8085)                          │
│                                                            │
│  ┌────────────────────────────────────────────────────┐         │
│  │         HTTP Server (net/http)              │         │
│  └────────────┬───────────────────────────────────┘         │
│               │                                            │
│  ┌────────────▼────────────────────────────────────┐          │
│  │          Middleware Layer              │          │
│  │  - JWT Authentication                        │          │
│  │  - CORS                                   │          │
│  │  - Logging                                │          │
│  │  - Request/Response Tracing                │          │
│  └────────────┬───────────────────────────────────┘          │
│               │                                            │
│  ┌────────────▼────────────────────────────────────┐          │
│  │           Handler Layer                │          │
│  │  - Report Handlers                         │          │
│  │  - Datasource Handlers                     │          │
│  │  - Dashboard Handlers                      │          │
│  │  - Export Handlers                        │          │
│  └────────────┬───────────────────────────────────┘          │
│               │                                            │
│  ┌────────────▼────────────────────────────────────┐          │
│  │           Service Layer                 │          │
│  │  - Report Service                         │          │
│  │  - Datasource Service                    │          │
│  │  - Rendering Engine                       │          │
│  │  - Export Engine                         │          │
│  └────────────┬───────────────────────────────────┘          │
│               │                                            │
│  ┌────────────▼────────────────────────────────────┐          │
│  │          Repository Layer              │          │
│  │  - Report Repository                      │          │
│  │  - Datasource Repository                 │          │
│  │  - Dashboard Repository                  │          │
│  └────────────┬───────────────────────────────────┘          │
│               │                                            │
│  ┌────────────▼────────────────────────────────────┐          │
│  │          Database Layer (GORM)                   │          │
│  │  - MySQL 5.7+                             │          │
│  └──────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. 后端架构 (Go)

### 2.1 目录结构

```
jimureport-go/
├── cmd/
│   └── server/
│       └── main.go                 # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置管理
│   ├── auth/
│   │   ├── middleware.go         # JWT 中间件
│   │   ├── jwt.go              # JWT 生成/验证
│   │   └── claims.go           # Claims 定义
│   ├── models/
│   │   ├── report.go            # Report 模型
│   │   ├── datasource.go        # Datasource 模型
│   │   ├── dashboard.go         # Dashboard 模型
│   │   ├── export.go            # Export 模型
│   │   └── cell.go             # Cell 模型
│   ├── repository/
│   │   ├── report_repo.go       # Report 数据访问
│   │   ├── datasource_repo.go   # Datasource 数据访问
│   │   └── dashboard_repo.go   # Dashboard 数据访问
│   ├── service/
│   │   ├── report_service.go    # 报表业务逻辑
│   │   ├── datasource_service.go # 数据源业务逻辑
│   │   ├── render/
│   │   │   ├── engine.go       # 渲染引擎核心
│   │   │   ├── data.go         # 数据查询
│   │   │   └── html.go        # HTML 生成
│   │   └── export/
│   │       ├── excel.go        # Excel 导出
│   │       └── pdf.go          # PDF 导出
│   └── httpserver/
│       ├── server.go            # HTTP 服务器
│       ├── routes.go            # 路由注册
│       └── handlers/
│           ├── report.go        # 报表处理器
│           ├── datasource.go    # 数据源处理器
│           ├── dashboard.go     # 仪表盘处理器
│           ├── export.go        # 导出处理器
│           └── health.go        # 健康检查
├── pkg/
│   ├── database/
│   │   └── db.go              # 数据库初始化工具
│   └── response/
│       └── response.go          # 统一响应格式
└── go.mod
```

### 2.2 分层架构原则

| 层级 | 职责 | 规则 |
|------|--------|------|
| **Handler Layer** | HTTP 请求处理 | - 参数验证<br>- 调用 Service<br>- 返回 HTTP 响应 |
| **Service Layer** | 业务逻辑 | - 核心业务实现<br>- 事务管理<br>- 调用 Repository |
| **Repository Layer** | 数据访问 | - CRUD 操作<br>- GORM 查询<br>- 数据库交互 |
| **Model Layer** | 数据模型 | - GORM 模型<br>- JSON 序列化<br>- 数据验证 |

### 2.3 关键模块

#### 2.3.1 认证模块 (auth)

```go
// internal/auth/jwt.go
type JWTService interface {
    GenerateToken(userID, tenantID string, roles []string) (string, error)
    ValidateToken(token string) (*Claims, error)
}

// internal/auth/middleware.go
func JWTAuthMiddleware(jwtService JWTService) func(http.Handler) http.Handler
```

**路由保护：**
- 公开路由：`/health`, `/login`
- 受保护路由：`/api/*`, `/jmreport/*`, `/drag/*`

#### 2.3.2 报表服务 (service/report_service)

```go
type ReportService interface {
    Create(ctx context.Context, req *CreateReportRequest) (*Report, error)
    Update(ctx context.Context, req *UpdateReportRequest) (*Report, error)
    Get(ctx context.Context, id string) (*Report, error)
    List(ctx context.Context, tenantID string) ([]Report, error)
    Delete(ctx context.Context, id string) error
    Preview(ctx context.Context, id string, params map[string]interface{}) (string, error)
}
```

#### 2.3.3 渲染引擎 (service/render/)

```go
type RenderEngine interface {
    Render(ctx context.Context, reportID string, params map[string]interface{}) (string, error)
}

// 数据查询流程
// 1. 解析报表配置 (JSON)
// 2. 提取数据绑定 (datasource, table, field)
// 3. 查询数据库获取数据
// 4. 生成 HTML (填充数据到单元格)
```

---

## 3. 前端架构 (Vue 3 + TypeScript)

### 3.1 目录结构

```
frontend/
├── public/                    # 静态资源
├── src/
│   ├── api/                  # API 调用层
│   │   ├── client.ts         # Axios 客户端配置
│   │   ├── auth.ts          # 认证 API
│   │   ├── report.ts        # 报表 API
│   │   ├── datasource.ts    # 数据源 API
│   │   └── dashboard.ts     # 仪表盘 API
│   ├── components/            # 公共组件
│   │   ├── common/
│   │   │   └── AppLayout.vue    # 应用布局
│   │   ├── report/
│   │   │   ├── ReportCanvas.vue     # 报表画布
│   │   │   ├── PropertyPanel.vue    # 属性面板
│   │   │   └── CellEditor.vue      # 单元格编辑器
│   │   └── dashboard/
│   │       ├── DashboardCanvas.vue  # 仪表盘画布
│   │       └── ComponentPanel.vue  # 组件面板
│   ├── views/                # 页面视图
│   │   ├── Login.vue              # 登录页
│   │   ├── ReportList.vue         # 报表列表
│   │   ├── ReportDesigner.vue     # 报表设计器
│   │   ├── ReportPreview.vue      # 报表预览
│   │   ├── DashboardList.vue      # 仪表盘列表
│   │   └── DashboardDesigner.vue  # 仪表盘设计器
│   ├── stores/               # Pinia 状态管理
│   │   ├── auth.ts          # 认证状态
│   │   ├── report.ts        # 报表状态
│   │   └── datasource.ts    # 数据源状态
│   ├── router/               # Vue Router
│   │   └── index.ts
│   ├── types/               # TypeScript 类型
│   │   ├── report.ts        # 报表类型
│   │   ├── datasource.ts    # 数据源类型
│   │   └── api.ts          # API 响应类型
│   ├── utils/               # 工具函数
│   │   ├── request.ts       # 请求封装
│   │   └── storage.ts       # 本地存储
│   ├── App.vue              # 根组件
│   └── main.ts              # 入口文件
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
```

### 3.2 分层架构

| 层级 | 职责 | 技术栈 |
|------|--------|----------|
| **View Layer** | UI 组件和页面 | Vue 3, Element Plus |
| **State Layer** | 状态管理 | Pinia |
| **API Layer** | HTTP 请求封装 | Axios |
| **Type Layer** | TypeScript 类型 | 类型定义 |

### 3.3 状态管理 (Pinia)

```typescript
// stores/auth.ts
export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: '',
    user: null,
    tenantId: '',
  }),
  actions: {
    login(credentials: LoginRequest) { },
    logout() { },
    setToken(token: string) { },
  }
})

// stores/report.ts
export const useReportStore = defineStore('report', {
  state: () => ({
    reports: [],
    currentReport: null,
  }),
  actions: {
    fetchReports() { },
    saveReport(report: ReportConfig) { },
  }
})
```

### 3.4 API 调用层

```typescript
// api/client.ts
import axios from 'axios'

const client = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// 请求拦截器（添加 Token）
client.interceptors.request.use(config => {
  const token = useAuthStore().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
    config.headers['X-Access-Token'] = token
  }
  return config
})

// 响应拦截器（错误处理）
client.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 401) {
      useAuthStore().logout()
    }
    return Promise.reject(error)
  }
)

export default client
```

---

## 4. API 设计

### 4.1 通用原则

| 规则 | 说明 |
|------|------|
| RESTful | 使用标准 HTTP 方法（GET, POST, PUT, DELETE） |
| 版本化 | 所有 API 以 `/api/v1/` 开头 |
| 统一响应 | 使用 `{ success, result, message, timestamp }` 格式 |
| 错误码 | HTTP 状态码 + 业务错误码 |

### 4.2 统一响应格式

```json
// 成功响应
{
  "success": true,
  "result": { /* data */ },
  "message": "success",
  "timestamp": 1706952000000
}

// 失败响应
{
  "success": false,
  "result": null,
  "message": "Error message",
  "code": "ERROR_CODE",
  "timestamp": 1706952000000
}
```

### 4.3 核心端点

#### 认证

| 方法 | 路径 | 说明 |
|------|--------|------|
| POST | `/api/v1/auth/login` | 登录，返回 JWT |
| POST | `/api/v1/auth/logout` | 登出 |
| GET | `/api/v1/auth/me` | 获取当前用户信息 |

#### 报表

| 方法 | 路径 | 说明 |
|------|--------|------|
| GET | `/api/v1/jmreport/list` | 获取报表列表 |
| GET | `/api/v1/jmreport/get?id={id}` | 获取报表详情 |
| POST | `/api/v1/jmreport/create` | 创建报表 |
| POST | `/api/v1/jmreport/update` | 更新报表 |
| DELETE | `/api/v1/jmreport/delete?id={id}` | 删除报表 |
| POST | `/api/v1/jmreport/preview` | 预览报表 |
| GET | `/jmreport/list` | 报表设计器入口（返回 HTML） |

#### 数据源

| 方法 | 路径 | 说明 |
|------|--------|------|
| GET | `/datasource/list` | 获取数据源列表 |
| POST | `/datasource/create` | 创建数据源 |
| POST | `/datasource/update` | 更新数据源 |
| DELETE | `/datasource/delete/{id}` | 删除数据源 |
| POST | `/datasource/test` | 测试数据源连接 |
| GET | `/datasource/{id}/tables` | 获取数据源的所有表 |
| GET | `/datasource/{id}/tables/{table}/fields` | 获取表的字段列表 |

#### 仪表盘

| 方法 | 路径 | 说明 |
|------|--------|------|
| GET | `/api/v1/drag/list` | 获取仪表盘列表 |
| GET | `/api/v1/drag/get?id={id}` | 获取仪表盘详情 |
| POST | `/api/v1/drag/create` | 创建仪表盘 |
| POST | `/api/v1/drag/update` | 更新仪表盘 |
| DELETE | `/api/v1/drag/delete?id={id}` | 删除仪表盘 |
| GET | `/drag/list` | 仪表盘设计器入口（返回 HTML） |

#### 导出

| 方法 | 路径 | 说明 |
|------|--------|------|
| POST | `/api/v1/jmreport/export` | 导出报表（Excel/PDF） |
| GET | `/api/v1/export/job?id={id}` | 获取导出任务状态 |
| GET | `/api/v1/export/download?id={id}` | 下载导出文件 |

---

## 5. 数据流

### 5.1 报表创建流程

```
用户 → ReportDesigner UI (Vue)
  ↓
填写报表信息 → 点击保存
  ↓
POST /api/v1/jmreport/create
  ↓
后端 Handler 层
  ↓
后端 Service 层
  ↓
后端 Repository 层
  ↓
MySQL → 保存到 jimu_report 表
  ↓
返回 { success: true, result: { id: "xxx" } }
  ↓
前端显示成功提示
```

### 5.2 报表预览流程

```
用户 → ReportPreview UI (Vue)
  ↓
点击预览按钮
  ↓
POST /api/v1/jmreport/preview { id: "xxx", params: {} }
  ↓
后端 Handler 层
  ↓
后端 ReportService.Preview()
  ↓
1. 从数据库加载报表配置
  ↓
2. 解析 JSON 配置，提取数据绑定
  ↓
3. RenderEngine.QueryData() → 查询数据源获取数据
  ↓
4. RenderEngine.GenerateHTML() → 生成 HTML（填充数据）
  ↓
返回 { success: true, result: { html: "<table>...</table>" } }
  ↓
前端显示 HTML
```

### 5.3 认证流程

```
用户 → Login UI (Vue)
  ↓
输入用户名密码 → 点击登录
  ↓
POST /api/v1/auth/login { username, password }
  ↓
后端验证用户（集成 Sa-Token 或本地验证）
  ↓
生成 JWT Token (包含 userId, roles, tenantId)
  ↓
返回 { success: true, result: { token: "jwt..." } }
  ↓
前端存储 Token 到 localStorage 和 Pinia Store
  ↓
所有后续请求自动携带 Token
```

---

## 6. 技术选型确认

| 模块 | 技术 | 版本 | 说明 |
|------|------|------|
| **后端** | | |
| 语言 | Go | 1.22+ | 高性能、并发友好 |
| Web 框架 | net/http | 标准库 | 无需第三方路由库 |
| ORM | GORM | 1.25+ | 数据库抽象 |
| 数据库 | MySQL | 5.7+ | 持久化存储 |
| JWT | golang-jwt/jwt | 5.2+ | JWT 生成/验证 |
| **前端** | | |
| 框架 | Vue | 3.4+ | 响应式 UI |
| 语言 | TypeScript | 5.3+ | 类型安全 |
| 构建 | Vite | 5.0+ | 快速开发 |
| UI 组件 | Element Plus | 2.5+ | Vue 3 组件库 |
| 状态 | Pinia | 2.1+ | 官方状态管理 |
| 路由 | Vue Router | 4.2+ | 路由管理 |
| HTTP | Axios | 1.6+ | HTTP 客户端 |
| 编辑器 | Monaco Editor | 最新 | 表达式编辑器 |

---

## 7. 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                   Nginx / CDN                    │
│                                                            │
│  ┌─────────────┬────────────┬──────────────┐         │
│  │   静态资源  │  API 反向代理 │   Vue Router │         │
│  │  (frontend/  │  → 8085      │    History    │         │
│  └─────────────┴────────────┴──────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

**开发环境：**
- 前端：`http://localhost:3000` (Vite Dev Server)
- 后端：`http://localhost:8085` (Go HTTP Server)
- 数据库：`MySQL:3306`

**生产环境：**
- Nginx 统一监听 80/443
- 前端静态文件由 Nginx 直接提供服务
- API 请求代理到后端 8085

---

## 8. 架构决策确认

以下决策已由用户确认：

### 8.1 后端决策

✅ **JWT 认证方式**：独立验证（后端用户表）
- 后端维护 `users` 表，独立验证用户名密码
- 生成 JWT Token，包含 userId, roles, tenantId
- 不依赖外部认证服务

✅ **Redis 缓存层**：需要
- 使用 Redis 缓存数据源连接
- 缓存查询结果提升性能
- 使用 `go-redis/v9` 或 `redigo`

### 路由库决策

**考虑选项对比：**

| 特性 | `net/http` + `ServeMux` | Gin |
|------|------------------------|-----|
| 参数绑定 | 需要手写 | 自动 JSON/Query/Path 绑定 |
| 路由分组 | 需要手写前缀 | `Group()` 支持 |
| 中间件链 | 手动包装 Handler | `Use()` 自动链式调用 |
| JSON 响应 | 手动 json.Marshal | `c.JSON()` 自动处理 |
| 错误处理 | 手动处理 | 内置 Recovery + 错误处理 |
| 性能 | 良好 | 优秀 |
| 学习曲线 | 低（但功能受限） | 中（文档完善） |
| 社区 | 官方支持 | 活跃、生态丰富 |

**建议：使用 Gin**

**理由：**
1. **功能更完整**：
   - 路由分组：`api := r.Group("/api/v1")`
   - 参数绑定：自动解析 JSON、Query、Path 参数
   - 中间件：JWT、CORS、日志链式调用
   - JSON 响应：`c.JSON(200, gin.H{"success": true})`

2. **开发效率高**：
   - 减少样板代码
   - 内置参数验证
   - 错误处理统一

3. **性能优秀**：
   - 基于 Radix Tree 路由匹配
   - 无 GC 压力
   - 支持 HTTP/2

4. **社区生态好**：
   - 大量中间件（Rate Limit、CORS、JWT）
   - 丰富的示例和文档
   - 问题易解决

✅ **推荐：使用 Gin**

**依赖更新：**
```go
// go.mod
require (
    github.com/gin-gonic/gin v1.9.1
    golang-jwt/jwt/v5 v5.2.0
    gorm.io/driver/mysql v1.5.2
    gorm.io/gorm v1.25.5
    github.com/redis/go-redis/v9 v9.0.5
)
```

**Gin 路由示例：**
```go
package httpserver

import (
    "github.com/gin-gonic/gin"
)

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
    r := gin.Default()

    // 中间件
    r.Use(middleware.CORS())
    r.Use(middleware.Logger())
    r.Use(middleware.JWTAuth())

    // 路由分组
    api := r.Group("/api/v1")
    {
        // 报表路由
        reports := api.Group("/jmreport")
        {
            reports.GET("/list", reportHandler.List)
            reports.GET("/get", reportHandler.Get)
            reports.POST("/create", reportHandler.Create)
            reports.POST("/update", reportHandler.Update)
            reports.DELETE("/delete", reportHandler.Delete)
            reports.POST("/preview", reportHandler.Preview)
        }

        // 数据源路由
        datasources := api.Group("/datasource")
        {
            datasources.GET("/list", datasourceHandler.List)
            datasources.POST("/create", datasourceHandler.Create)
            datasources.POST("/test", datasourceHandler.Test)
            datasources.GET("/:id/tables", datasourceHandler.GetTables)
            datasources.GET("/:id/tables/:table/fields", datasourceHandler.GetFields)
        }

        // 仪表盘路由
        dashboards := api.Group("/drag")
        {
            dashboards.GET("/list", dashboardHandler.List)
            dashboards.GET("/get", dashboardHandler.Get)
            dashboards.POST("/create", dashboardHandler.Create)
            dashboards.POST("/update", dashboardHandler.Update)
            dashboards.DELETE("/delete", dashboardHandler.Delete)
        }
    }

    // 兼容旧路由（返回 HTML）
    r.GET("/jmreport/list", reportHandler.ListPage)
    r.GET("/drag/list", dashboardHandler.ListPage)
    r.GET("/health", healthHandler.Check)

    return &Server{
        Addr: cfg.Server.Addr,
        Engine: r,
    }
}
```

### 8.2 前端决策

✅ **报表画布技术**：DOM + Flex Grid
- 使用 HTML div + CSS Grid/Flexbox 实现
- 开发简单，易于调试，适合 MVP
- 后续可优化为 Canvas API

✅ **表达式编辑器**：Monaco Editor
- 完整的代码编辑器功能
- 语法高亮、自动完成
- 适合复杂表达式编辑

✅ **状态管理**：Pinia
- 全局状态管理
- 便于跨组件共享状态
- Vue 3 官方推荐

### 8.3 数据库决策

✅ **数据库 Schema**：新建 Go 版本表
- 不复用 Java 版本数据库表
- 重新设计符合 Go 后端的表结构
- 需要手动迁移现有数据（如有）

---

## 9. 下一步

确认架构后，将按以下顺序实现：

**Phase 1：基础设施**（1 天）
- [x] Go 项目脚手架
- [x] Vue 项目脚手架
- [ ] OpenSpec 变更：infrastructure-setup
- [ ] 确认技术选型和架构

**Phase 2：认证和数据源**（2 天）
- [ ] OpenSpec 变更：auth-jwt-implementation
- [ ] OpenSpec 变更：datasource-api
- [ ] 实现 JWT 中间件
- [ ] 实现数据源 CRUD

**Phase 3：报表核心**（4 天）
- [ ] OpenSpec 变更：report-designer-implementation
- [ ] OpenSpec 变更：report-renderer-implementation
- [ ] 实现报表 CRUD API
- [ ] 实现渲染引擎
- [ ] 实现报表设计器 UI
- [ ] 实现报表预览 UI

**Phase 4：高级功能**（3 天）
- [ ] OpenSpec 变更：bi-dashboard-implementation
- [ ] OpenSpec 变更：report-export-implementation
- [ ] 实现仪表盘功能
- [ ] 实现导出功能

---

## 10. 附加说明

- 所有 API 遵循 RESTful 规范
- 前后端完全分离，通过 HTTP API 交互
- 使用 JWT Token 进行身份认证
- 支持多租户（通过 JWT tenantId）
- 数据库事务管理在 Service 层
- 日志记录使用结构化日志（JSON）
- 前端使用 TypeScript 严格模式

# goReport 开发指南

## 目录

- [开发环境搭建](#开发环境搭建)
- [项目结构](#项目结构)
- [后端开发](#后端开发)
- [前端开发](#前端开发)
- [数据库设计](#数据库设计)
- [API 规范](#api-规范)
- [测试指南](#测试指南)
- [构建部署](#构建部署)

## 开发环境搭建

### 环境要求

- **Go**: 1.22+
- **Node.js**: 20+
- **Docker**: 20.10+
- **Docker Compose**: 2.0+

### 快速开始

```bash
# 克隆项目
git clone git@github.com:Gujiaweiguo/goreport.git
cd goreport

# 启动开发环境
make dev
```

### 本地开发

#### 后端

```bash
cd backend

# 安装依赖
go mod download

# 运行开发服务器
go run cmd/server/main.go

# 运行测试
go test ./... -v

# 生成覆盖率报告
make test-coverage
```

#### 前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览构建结果
npm run preview
```

## 项目结构

```
goreport/
├── backend/                 # Go 后端
│   ├── cmd/
│   │   └── server/        # 应用入口
│   ├── internal/
│   │   ├── auth/          # 认证模块
│   │   ├── config/        # 配置管理
│   │   ├── dashboard/     # 仪表盘模块
│   │   ├── datasource/    # 数据源模块
│   │   ├── models/        # 数据模型
│   │   ├── report/        # 报表模块
│   │   └── httpserver/    # HTTP 服务
│   ├── db/
│   │   └── init.sql     # 数据库初始化
│   ├── go.mod
│   └── go.sum
├── frontend/               # Vue 前端
│   ├── src/
│   │   ├── api/           # API 调用
│   │   ├── canvas/        # Canvas 画布
│   │   ├── components/    # Vue 组件
│   │   ├── stores/        # Pinia 状态
│   │   ├── types/         # TypeScript 类型
│   │   ├── utils/         # 工具函数
│   │   └── views/         # 页面视图
│   ├── index.html
│   ├── package.json
│   └── vite.config.ts
├── docs/                  # 文档
├── openspec/              # 需求规范
├── docker-compose.yml       # 开发环境
├── docker-compose.prod.yml # 生产环境
└── Makefile
```

## 后端开发

### 代码分层

后端采用分层架构：

```
Handler (HTTP 层)
    ↓
Service (业务逻辑层)
    ↓
Repository (数据访问层)
    ↓
Database (数据库)
```

### 创建新模块

1. 创建目录结构：

```bash
mkdir -p backend/internal/module
```

2. 实现各层：

```go
// models/model.go
type Model struct {
    ID   string `gorm:"primaryKey;type:varchar(36)"`
    Name string `gorm:"type:varchar(100)"`
}

// repository/repository.go
type Repository interface {
    Create(model *Model) error
    Get(id string) (*Model, error)
}

// service/service.go
type Service interface {
    Create(ctx context.Context, req *CreateRequest) (*Model, error)
}

// handler/handler.go
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request"})
        return
    }
    // ... 业务逻辑
}
```

3. 注册路由：

```go
// server.go
handler := NewHandler(service)
r.Group("/api/v1/module") {
    r.POST("/create", handler.Create)
    r.GET("/:id", handler.Get)
}
```

### 配置管理

在 `config/config.go` 中添加配置：

```go
type Config struct {
    Module ModuleConfig
}

type ModuleConfig struct {
    Enabled bool
    Option  string
}

func Load() (*Config, error) {
    return &Config{
        Module: ModuleConfig{
            Enabled: getBoolEnv("MODULE_ENABLED", false),
            Option:  getEnv("MODULE_OPTION", ""),
        },
    }, nil
}
```

### 错误处理

使用统一的错误处理中间件：

```go
c.JSON(http.StatusBadRequest, gin.H{
    "success": false,
    "message": "error message",
})
```

### 数据库迁移

```go
err := db.AutoMigrate(&models.Model{})
if err != nil {
    log.Fatalf("Failed to migrate: %v", err)
}
```

## 前端开发

### 组件开发

创建 Vue 3 组件：

```vue
<template>
  <div class="my-component">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Props {
  title: string
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Default Title'
})
</script>

<style scoped>
.my-component {
  padding: 16px;
}
</style>
```

### 状态管理

使用 Pinia 管理状态：

```typescript
// stores/module.ts
import { defineStore } from 'pinia'

export const useModuleStore = defineStore('module', () => {
  const items = ref([])

  function addItem(item: Item) {
    items.value.push(item)
  }

  return { items, addItem }
})
```

### API 调用

```typescript
// api/module.ts
import axios from 'axios'

export interface Module {
  id: string
  name: string
}

export const moduleApi = {
  list: () => axios.get<ApiResponse<Module[]>>('/api/v1/module/list'),
  create: (data: CreateModuleRequest) => axios.post<ApiResponse<Module>>('/api/v1/module/create', data),
  update: (id: string, data: UpdateModuleRequest) => axios.put<ApiResponse<Module>>(`/api/v1/module/${id}`, data),
  delete: (id: string) => axios.delete<ApiResponse<null>>(`/api/v1/module/${id}`)
}
```

### 路由配置

```typescript
// router/index.ts
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/module',
    component: ModuleView,
    meta: { title: '模块管理', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  // 认证检查
  if (to.meta.requiresAuth && !isAuthenticated()) {
    next('/login')
  } else {
    next()
  }
})
```

## 数据库设计

### 表命名规范

- 使用复数形式：`reports`, `dashboards`, `data_sources`
- 使用下划线分隔：`user_tenants`, `export_jobs`

### 字段命名规范

- 使用下划线分隔：`created_at`, `updated_at`, `deleted_at`
- 主键统一为 `id`，类型 `varchar(36)`
- 外键统一为 `xxx_id`，如 `tenant_id`, `user_id`

### 通用字段

```sql
id VARCHAR(36) PRIMARY KEY,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
deleted_at TIMESTAMP NULL,
```

### 索引设计

```sql
-- 单列索引
INDEX idx_tenant_id (tenant_id),
INDEX idx_status (status)

-- 联合索引
INDEX idx_tenant_status (tenant_id, status)
```

## API 规范

### URL 规范

```
GET    /api/v1/{resource}/list     # 列表
GET    /api/v1/{resource}/:id      # 详情
POST   /api/v1/{resource}/create  # 创建
PUT    /api/v1/{resource}/:id     # 更新
DELETE /api/v1/{resource}/:id     # 删除
```

### 响应格式

成功响应：

```json
{
  "success": true,
  "result": { /* 数据 */ },
  "message": "success"
}
```

失败响应：

```json
{
  "success": false,
  "message": "error message",
  "code": 400
}
```

### 认证

使用 JWT 认证：

```go
// 生成 Token
token, err := GenerateToken(user)

// 验证 Token（中间件自动处理）
claims, err := ValidateToken(token)

// 获取用户信息
userID := GetUserID(c)
tenantID := GetTenantID(c)
```

## 测试指南

### 单元测试

```go
func TestService_Create(t *testing.T) {
    repo := &mockRepository{}
    service := NewService(repo)

    req := &CreateRequest{
        Name: "Test",
    }

    result, err := service.Create(context.Background(), req)
    if err != nil {
        t.Fatalf("Create failed: %v", err)
    }

    if result == nil {
        t.Fatal("Expected non-nil result")
    }
}
```

### 集成测试

```go
func TestIntegration_Create(t *testing.T) {
    router := setupTestRouter(t)

    req := map[string]interface{}{
        "name": "Test",
    }
    body, _ := json.Marshal(req)

    httpReq, _ := http.NewRequest("POST", "/api/v1/resource/create", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    router.ServeHTTP(w, httpReq)

    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
}
```

### 运行测试

```bash
# 运行所有测试
go test ./... -v

# 运行特定包
go test ./internal/module -v

# 生成覆盖率
make test-coverage
```

## 构建部署

### 前端构建

```bash
# 构建生产版本
cd frontend
npm run build

# 输出目录：dist/
```

### 后端构建

```bash
# 构建二进制
cd backend
go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# 输出文件：bin/server
```

### Docker 构建

```bash
# 开发环境
make dev

# 生产环境
make build-prod
```

### 环境变量

创建 `.env` 文件：

```bash
DB_DSN=root:password@tcp(mysql:3306)/goreport?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=your-secret-key
CACHE_ENABLED=true
```

## 常见问题

### 端口冲突

修改 `frontend/vite.config.ts` 中的端口：

```typescript
server: {
  port: 3001  // 修改为其他端口
}
```

### 依赖问题

```bash
# 后端
go mod tidy

# 前端
rm -rf node_modules package-lock.json
npm install
```

### 数据库连接

检查 Docker 容器状态：

```bash
make ps

# 查看日志
make logs
```

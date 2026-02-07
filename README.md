# JimuReport 积木报表

基于 OpenSpec 规范构建的新一代报表系统。

## 技术栈

### 后端
- **Go 1.22+** - 高性能后端语言
- **Gin** - Web 框架（路由、中间件、参数绑定）
- **GORM** - ORM 框架（MySQL 支持）
- **Redis** - 缓存（数据源连接、查询结果）
- **JWT** - 认证（golang-jwt/jwt）
- **Docker** - 容器化部署

### 前端
- **Vue 3** - 前端框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Element Plus** - UI 组件库
- **Canvas API** - 报表设计器（高性能画布）
- **Monaco Editor** - 表达式编辑器
- **Pinia** - 状态管理

## 项目结构

```
jimureport/
├── backend/                    # Go 后端
│   ├── cmd/server/            # 应用入口
│   ├── internal/              # 内部代码
│   │   ├── auth/             # JWT 认证
│   │   ├── config/           # 配置管理
│   │   ├── models/           # 数据模型
│   │   ├── repository/       # 数据访问层
│   │   ├── service/          # 业务逻辑层
│   │   │   ├── render/      # 渲染引擎
│   │   │   └── export/      # 导出服务
│   │   ├── httpserver/       # HTTP 服务
│   │   │   └── handlers/    # 请求处理器
│   │   └── middleware/       # 中间件
│   ├── pkg/                   # 公共包
│   │   ├── database/        # 数据库工具
│   │   └── cache/           # Redis 缓存
│   ├── db/
│   │   └── init.sql         # 数据库初始化脚本
│   ├── Dockerfile.dev       # 开发环境 Dockerfile
│   ├── .air.toml            # 热重载配置
│   └── go.mod               # Go 模块
├── frontend/                   # Vue 前端
│   ├── src/
│   │   ├── views/           # 页面视图
│   │   │   ├── report/     # 报表相关
│   │   │   └── dashboard/  # 仪表盘相关
│   │   ├── components/      # 组件
│   │   │   ├── report/     # 报表组件
│   │   │   ├── dashboard/  # 仪表盘组件
│   │   │   └── common/     # 公共组件
│   │   ├── canvas/          # Canvas 画布
│   │   ├── api/            # API 调用
│   │   ├── stores/         # Pinia 状态
│   │   ├── types/          # TypeScript 类型
│   │   └── utils/          # 工具函数
│   ├── public/              # 静态资源
│   ├── index.html
│   ├── package.json
│   ├── vite.config.ts
│   └── Dockerfile.dev
├── docker-compose.yml         # 开发环境配置
├── Makefile                 # 常用命令
└── openspec/                # 需求规范
    ├── specs/              # 规格定义
    └── changes/            # 变更提案
```

## 快速开始

### 方式一：Docker Compose（推荐）

一键启动完整开发环境：

```bash
# 1. 克隆项目
git clone <repository>
cd jimureport

# 2. 启动开发环境
make dev

# 3. 查看服务状态
make ps

# 4. 查看日志
make logs
```

访问地址：
- 前端：http://localhost:3000
- 后端 API：http://localhost:8085
- MySQL：localhost:3306 (root/root)
- Redis：localhost:6379

### 缓存配置

系统支持 Redis 缓存，可提升性能。配置参数：

| 参数 | 说明 | 默认值 |
|------|------|---------|
| CACHE_ENABLED | 是否启用缓存 | false |
| CACHE_ADDR | Redis 地址 | localhost:6379 |
| CACHE_PASSWORD | Redis 密码 | （空） |
| CACHE_DB | Redis DB | 0 |
| CACHE_DEFAULT_TTL | 默认 TTL（秒） | 3600 |

缓存观测端点：
- GET /api/v1/cache/metrics - 查看缓存命中率、失败次数等指标

### 方式二：本地开发

#### 后端

```bash
cd backend

# 安装依赖
go mod download

# 启动开发服务器（支持热重载）
air

# 或者
# go run cmd/server/main.go
```

#### 前端

```bash
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 常用命令

```bash
# 开发环境
make dev          # 启动所有服务
make dev-down     # 停止所有服务
make dev-logs     # 查看日志
make ps           # 查看容器状态

# 数据库
db-shell          # 进入 MySQL
redis-cli         # 进入 Redis

# 构建和测试
make build        # 构建生产镜像
make test         # 运行测试
make clean        # 清理容器和卷
```

## 开发规范

### 后端规范

1. **代码结构**
   - Handler：HTTP 请求处理，参数验证
   - Service：业务逻辑，事务管理
   - Repository：数据访问，GORM 查询
   - Model：数据模型，JSON 标签

2. **API 设计**
   - RESTful API，版本化（/api/v1/）
   - 统一响应格式：`{ success, result, message, timestamp }`
   - 错误码：HTTP 状态码 + 业务错误码

3. **数据库**
   - 表名：复数形式（users, reports）
   - 字段名：下划线命名（created_at）
   - ID：UUID（VARCHAR(36)）
   - 软删除：deleted_at 字段
   - 时间戳：created_at, updated_at

### 前端规范

1. **代码结构**
   - Views：页面级组件
   - Components：可复用组件
   - API：HTTP 请求封装
   - Stores：Pinia 状态管理
   - Types：TypeScript 类型定义

2. **Canvas 开发**
   - 使用 requestAnimationFrame 渲染
   - 事件委托处理交互
   - 虚拟滚动优化大数据
   - 高清屏适配（devicePixelRatio）

## 功能模块

### Phase 1：基础设施
- [x] 项目脚手架
- [x] Docker 开发环境
- [x] 数据库初始化
- [ ] JWT 认证
- [x] Redis 缓存

### Phase 2：认证和数据源
- [ ] 用户/租户管理
- [ ] JWT 认证中间件
- [ ] 数据源 CRUD
- [ ] 数据源连接测试
- [ ] 元数据查询（表/字段）

### Phase 3：报表核心
- [ ] 报表 CRUD API
- [ ] Canvas 报表设计器
- [ ] 单元格操作（选择、编辑、样式）
- [ ] 数据绑定（数据源、表、字段）
- [ ] 渲染引擎（数据查询 + HTML 生成）
- [ ] 报表预览

### Phase 4：高级功能
- [ ] BI 仪表盘
- [ ] 图表组件
- [ ] 导出功能（Excel、PDF）
- [ ] 表达式编辑器
- [ ] 报表参数

## 文档

- [TECHNICAL_DECISIONS.md](./TECHNICAL_DECISIONS.md) - 技术选型对比
- [ARCHITECTURE.md](./ARCHITECTURE.md) - 系统架构设计
- [docs/USER_GUIDE.md](./docs/USER_GUIDE.md) - 用户指南
- [docs/DEVELOPMENT_GUIDE.md](./docs/DEVELOPMENT_GUIDE.md) - 开发指南
- [docs/MIGRATION_GUIDE.md](./docs/MIGRATION_GUIDE.md) - 迁移指南
- [docs/CONTRIBUTING.md](./docs/CONTRIBUTING.md) - 贡献指南
- [docs/BROWSER_COMPATIBILITY_TEST.md](./docs/BROWSER_COMPATIBILITY_TEST.md) - 浏览器兼容性测试指南
- [docs/UX_OPTIMIZATION_GUIDE.md](./docs/UX_OPTIMIZATION_GUIDE.md) - 用户体验优化指南
- [docs/UAT_GUIDE.md](./docs/UAT_GUIDE.md) - 用户验收测试指南
- [openspec/](./openspec/) - 需求规范（OpenSpec）

## 贡献

1. Fork 项目
2. 创建特性分支
3. 提交代码
4. 创建 Pull Request

## 许可证

LGPL-3.0

## 联系方式

- 问题反馈：GitHub Issues
- 技术支持：<support@jimureport.com>

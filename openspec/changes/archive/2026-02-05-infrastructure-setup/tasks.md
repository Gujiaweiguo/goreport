# Infrastructure Setup - 任务清单

> **预计耗时**：1 天
> **开始日期**：2026-02-03

---

## ✅ 阶段一：Go 后端基础（半天）

### 任务 1.1：项目结构初始化
**文件**：`cmd/server/main.go`, `internal/`, `go.mod`

**目标**：创建标准 Go 项目结构

**内容**：
- 创建 cmd/server/main.go 入口文件
- 创建 internal/{config,database,auth,models,httpserver} 目录
- 创建 pkg/ 目录
- 配置 go.mod 依赖

**验收标准**：
- [x] Go 项目结构符合标准布局
- [x] go mod tidy 成功执行

---

### 任务 1.2：配置管理
**文件**：`internal/config/config.go`

**目标**：实现配置加载和环境变量支持

**内容**：
- Config 结构体（Server、Database、JWT）
- 环境变量默认值
- Load() 函数

**验收标准**：
- [x] 支持环境变量配置
- [x] 提供合理的默认值

---

### 任务 1.3：数据库连接
**文件**：`internal/database/database.go`

**目标**：初始化 GORM 数据库连接

**内容**：
- Init() 函数
- MySQL DSN 配置
- 连接成功日志

**验收标准**：
- [x] 能成功连接 MySQL
- [x] 返回 *gorm.DB 实例

---

### 任务 1.4：HTTP 服务器
**文件**：`internal/httpserver/server.go`, `internal/httpserver/handlers/health.go`

**目标**：创建 HTTP 服务器和健康检查

**内容**：
- Server 结构体
- NewServer() 函数
- /health 健康检查端点
- 优雅关闭

**验收标准**：
- [x] 服务器能启动
- [x] /health 返回状态
- [x] 支持优雅关闭

---

## 🎨 阶段二：Vue 前端基础（半天）

### 任务 2.1：项目初始化
**文件**：`package.json`, `vite.config.ts`, `tsconfig.json`

**目标**：创建 Vue 3 + TypeScript 项目

**内容**：
- package.json 依赖（Vue 3、Vite、Element Plus）
- Vite 配置
- TypeScript 配置
- 路径别名 @/

**验收标准**：
- [x] Vite 开发服务器能启动
- [x] TypeScript 编译无错误

---

### 任务 2.2：应用入口
**文件**：`frontend/src/main.ts`, `frontend/src/App.vue`

**目标**：创建应用入口和根组件

**内容**：
- main.ts 应用初始化
- App.vue 根组件
- Pinia 状态管理
- Element Plus 集成

**验收标准**：
- [x] 应用能启动
- [x] Element Plus 组件可用

---

### 任务 2.3：路由和视图
**文件**：`frontend/src/router/index.ts`, `frontend/src/views/Home.vue`

**目标**：创建路由系统和首页

**内容**：
- Vue Router 配置
- Home 首页视图
- 健康检查测试按钮

**验收标准**：
- [x] 路由正常工作
- [x] 能调用后端 /health

---

## 📊 总计

- **总任务数**：6
- **预计总耗时**：0.5 天
- **参与模块**：Go 后端 4 个，Vue 前端 3 个
- **风险等级**：低

---

## 🔄 任务状态

| ID | 阶段 | 任务 | 状态 |
|----|------|------|------|
| 1.1 | Go 后端 | 项目结构初始化 | ✅ 已完成 |
| 1.2 | Go 后端 | 配置管理 | ✅ 已完成 |
| 1.3 | Go 后端 | 数据库连接 | ✅ 已完成 |
| 1.4 | Go 后端 | HTTP 服务器 | ✅ 已完成（已验证） |
| 2.1 | Vue 前端 | 项目初始化 | ✅ 已完成（已验证） |
| 2.2 | Vue 前端 | 应用入口 | ✅ 已完成（已验证） |
| 2.3 | Vue 前端 | 路由和视图 | ✅ 已完成（已验证） |

---

## 📝 备注

- 所有任务已完成并已验证
- 前后端能正常启动和通信
- 后端验证成功：修复 gin.Engine.Shutdown 问题后，服务器能启动并返回 /health（200，database: connected）。需要环境变量 DB_DSN 指向可用 MySQL（如 localhost:3307）。
- 前端验证成功：修复 Vite root 配置（frontend/vite.config.ts:6）和添加 /health 代理（frontend/vite.config.ts:14）后，根路径返回 200，/health 通过代理成功调用后端（返回 database: connected）。

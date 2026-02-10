# Changelog

所有重要的变更都将记录在此文件中。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [1.0.0] - 2026-02-06

### 新增

#### 核心功能
- 大屏设计器：支持拖拽、组件库、数据绑定、图层管理、实时预览
- 图表编辑器：集成 ECharts 5.x，支持 20+ 图表类型
- 数据源管理：支持 MySQL，测试连接，元数据查询（表/字段）
- 报表设计器：Canvas 画布、单元格操作、数据绑定
- 报表预览：参数配置、数据展示、导出功能

#### 后端模块
- Dashboard 完整 CRUD API
- 统一错误处理中间件
- 数据库连接池优化
- JWT 认证和权限控制
- 安全测试（SQL 注入、XSS、路径遍历防护）

#### 前端组件
- 通用状态组件：EmptyState、LoadingState、ErrorState、NoPermission
- 布局组件：MainLayout（侧边栏导航）
- 图表组件：ChartTypeSelector、EChartsRenderer、DataConfigPanel、ChartPropertyPanel、ChartPreview
- 大屏组件：PropertyPanel、LayerPanel、DashboardPreview、ComponentLibrary

#### 工程化
- Docker 多阶段构建优化
- 生产环境 Docker Compose 配置
- 前端代码分割和懒加载
- Makefile 常用命令
- 环境变量配置模板

#### 文档
- 用户指南（USER_GUIDE.md）
- 开发指南（DEVELOPMENT_GUIDE.md）
- 迁移指南（MIGRATION_GUIDE.md）
- 贡献指南（CONTRIBUTING.md）
- 浏览器兼容性测试指南（BROWSER_COMPATIBILITY_TEST.md）
- 用户体验优化指南（UX_OPTIMIZATION_GUIDE.md）
- 用户验收测试指南（UAT_GUIDE.md）

### 测试

- 单元测试：Dashboard Service 层 100% 覆盖
- 集成测试：HTTP API 完整流程测试
- 安全测试：SQL 注入、XSS、权限验证等
- 测试通过率：100%

### 性能优化

- 前端：代码分割（vendor、element-plus、common）
- 后端：数据库连接池（MaxOpenConns: 100, MaxIdleConns: 10）
- 构建：资源压缩、静态资源内联

### OpenSpec

- 归档变更：`update-ui-feature-visibility` → `2026-02-06-update-ui-feature-visibility`
- 创建规范：`frontend-feature-availability`

### 依赖更新

#### 后端
- Go 1.22
- Gin 1.9.1
- GORM 1.30.0
- golang-jwt/jwt/v5 5.3.1
- gorm.io/driver/sqlite 1.6.0（测试用）

#### 前端
- Vue 3.4
- TypeScript 5.3
- Vite 5.0
- Element Plus 2.5
- ECharts 5.6
- Pinia 2.1
- Vue Router 4.2

## [Unreleased]

### 计划中

- PostgreSQL 数据源支持
- 报表批量导入/导出
- 实时数据刷新
- 移动端大屏查看优化
- AI 智能报表生成

## [1.0.1-phase2-stable] - 2026-02-10

### 新增

- 用户/租户管理 API：`/api/v1/users/me`、`/api/v1/tenants`、`/api/v1/tenants/current`
- API 端到端冒烟脚本：登录 → 数据源 → 数据集 → 报表预览（`make smoke-e2e`）
- CI 前端独立任务：`npm run test:run` + `npm run build`

### 变更

- 认证链路增强：路由守卫增加 `/users/me` 会话校验，401 统一清理会话
- 数据源与鉴权测试增强：补齐 JWT/中间件/黑名单测试与数据源 handler 覆盖
- 前端构建拆包优化：`vite` 手动拆分 `zrender/echarts/element-plus/vendor` chunk

### 修复

- 修复 `frontend/src/api/dataset.ts` 重复 `getSchema` 键定义导致的构建告警

### 清理

- 删除临时产物：`.bak` 备份文件与调试截图资源

---

## 版本历史

### 版本号说明

- **主版本号**：不兼容的 API 变更
- **次版本号**：向后兼容的功能新增
- **修订号**：向后兼容的问题修复

### 标签说明

- `[Added]` 新功能
- `[Changed]` 变更
- `[Deprecated]` 即将移除
- `[Removed]` 移除
- `[Fixed]` 修复
- `[Security]` 安全

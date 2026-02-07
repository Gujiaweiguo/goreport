# Project Context

## Purpose
JimuReport 是一个基于 OpenSpec 驱动开发的报表与可视化系统，覆盖报表设计、预览渲染、导出、仪表盘与大屏能力。
当前仓库以 Go 后端 + Vue 前端实现为主，同时保留 `jimureport-example/` 作为 Java 生态兼容与参考示例。

## Tech Stack
- Go 1.22+
- Gin（HTTP 路由与中间件）
- GORM（MySQL 访问）
- JWT（认证）
- Redis（可选缓存）
- Vue 3 + TypeScript
- Vite
- Element Plus
- Pinia
- ECharts（图表）
- Monaco Editor（表达式编辑）
- Docker / Docker Compose

## Project Conventions

### Code Style
- 默认 4 空格缩进，遵循所在语言既有风格
- 优先复用现有模式，避免无关重构
- 注释与日志可中英混合，保持就近文件风格一致
- 未经确认不新增依赖

### Architecture Patterns
- 后端主目录：`backend/cmd/server`、`backend/internal/{auth,config,models,repository,service,httpserver,middleware}`
- 前端主目录：`frontend/src/{views,components,api,stores,types,utils,router}`
- OpenSpec 目录：`openspec/specs`（当前规格）与 `openspec/changes`（变更提案）
- 兼容示例：`jimureport-example/`（Spring Boot 集成参考）

### Testing Strategy
- 优先执行项目脚本（如 `make test`）进行回归验证
- 后端以 Go 测试为主，前端按现有脚本执行单测/集成测试
- OpenSpec 变更需执行严格校验：`openspec validate --strict --no-interactive`

### Git Workflow
- 保持小步提交与聚焦变更
- 未明确要求时不提交 commit

## Domain Context
- 核心能力包括：报表设计器、报表渲染器、报表导出、Dashboard 与图表编辑
- API 与 UI 实现面向 JimuReport/JimuBI 兼容场景
- 目标是先保证核心流程可用，再逐步补齐高级能力与兼容性

## Important Constraints
- 后端开发环境要求 Go 1.22+
- 前端开发环境要求 Node.js（以仓库脚本为准）
- 依赖 MySQL；Redis 作为可选组件
- 变更遵循 OpenSpec 流程：提案 -> 实施 -> 归档

## External Dependencies
- MySQL（必需）
- Redis（可选）
- Docker / Docker Compose（推荐本地联调）
- 官方文档：https://help.jimureport.com

# Change: Infrastructure Setup - 基础设施搭建

## Why
项目代码已完全重置，需要建立基础的 Go 后端和 Vue 前端脚手架，以便后续基于 OpenSpec specs 逐步实现功能模块。

## What Changes
- 创建 Go 后端基础项目结构（config、database、httpserver、auth 等）
- 创建 Vue 3 + TypeScript + Element Plus 前端项目结构
- 配置数据库连接和健康检查端点
- 配置前端开发服务器和 API 代理

## Impact
- Affected specs: `openspec/specs/auth-jwt/spec.md`, `openspec/specs/datasource-management/spec.md`
- Affected code: `cmd/server/main.go`, `internal/`, `frontend/`

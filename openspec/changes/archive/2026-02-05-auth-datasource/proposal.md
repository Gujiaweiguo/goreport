# Change: Authentication and Datasource Management

## Why
Phase 1 基础设施已完成，现在需要实现核心业务功能：JWT 认证和数据源管理，为后续报表功能提供基础支撑。

## What Changes
- 实现 JWT 认证中间件和 Token 生成
- 实现用户登录/登出 API
- 实现数据源 CRUD API（创建、查询、更新、删除）
- 实现数据源连接测试
- 实现数据源元数据查询（获取表列表、字段列表）

## Impact
- Affected specs: `openspec/specs/auth-jwt/spec.md`, `openspec/specs/datasource-management/spec.md`
- Affected code:
  - `backend/internal/auth/` - JWT 认证
  - `backend/internal/models/` - 数据模型
  - `backend/internal/repository/` - 数据访问层
  - `backend/internal/httpserver/handlers/` - HTTP 处理器

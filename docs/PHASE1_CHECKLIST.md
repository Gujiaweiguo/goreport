# Phase 1 验收清单（认证与数据源）

## 1. 接口验收

### 1.1 认证

- [ ] `POST /api/v1/auth/login` 登录成功返回 `token` 与 `user`
- [ ] `POST /api/v1/auth/login` 错误凭证返回 401
- [ ] `POST /api/v1/auth/logout` 成功拉黑当前 token
- [ ] 受保护接口在无 token 时返回 401
- [ ] 受保护接口在非法/过期 token 时返回 401

### 1.2 用户/租户

- [ ] `GET /api/v1/users/me` 返回当前登录用户信息
- [ ] `GET /api/v1/tenants/current` 返回当前租户信息
- [ ] `GET /api/v1/tenants` 返回当前用户可访问租户列表

### 1.3 数据源

- [ ] `GET /api/v1/datasource/list` 返回租户内数据源列表
- [ ] `POST /api/v1/datasource/create` 成功创建数据源
- [ ] `PUT /api/v1/datasource/:id` 成功更新数据源
- [ ] `DELETE /api/v1/datasource/:id` 成功删除数据源
- [ ] `POST /api/v1/datasource/test` 能正确测试连接
- [ ] `GET /api/v1/datasource/:id/tables` 能返回表列表
- [ ] `GET /api/v1/datasource/:id/tables/:table/fields` 能返回字段列表

## 2. 自动化测试

### 2.1 后端

- [ ] `backend/internal/auth/jwt_test.go`
- [ ] `backend/internal/auth/middleware_test.go`
- [ ] `backend/internal/auth/blacklist_test.go`
- [ ] `backend/internal/httpserver/handlers/datasource_test.go`（扩展覆盖）

### 2.2 前端

- [ ] `frontend/src/api/auth.test.ts`
- [ ] `frontend/src/api/datasource.test.ts`

## 3. 手工回归路径

- [ ] 登录页面使用合法账号登录
- [ ] 进入数据源管理页面加载列表
- [ ] 新建数据源并保存成功
- [ ] 执行数据源连接测试并验证提示
- [ ] 拉取表列表
- [ ] 拉取字段列表
- [ ] 退出登录后受保护页面跳转到登录页

## 4. 命令验收

- [x] `cd backend && go test ./... -v`
- [x] `cd frontend && npm run test:run`
- [x] `cd frontend && npm run build`

## 5. 文档同步

- [ ] 更新 `README.md` Phase 1 勾选状态与实际一致

## 6. 本次执行结果（2026-02-10）

- 后端全量测试通过；部分集成测试因缺少 `TEST_DB_DSN/DB_DSN` 与 `TEST_REDIS_ADDR/REDIS_ADDR/CACHE_ADDR` 被按设计跳过。
- 前端测试通过，生产构建通过。
- 构建过程中存在既有告警（例如 `src/api/dataset.ts` duplicate key、chunk size warning），不影响本次 Phase 1 交付链路验证。

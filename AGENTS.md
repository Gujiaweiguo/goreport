<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# AGENTS.md

**Updated:** 2026-02-14

## OVERVIEW

GoReport - 报表设计与数据可视化平台。Go 1.23 后端 + Vue 3 前端。支持多数据源、Canvas 报表设计器、ECharts 图表。

## COMMANDS

```bash
# === 开发环境 ===
make dev                    # 启动所有服务 (Docker Compose)
make dev-logs               # 查看日志
make dev-down               # 停止服务

# === 后端测试 ===
cd backend && go test ./... -cover                              # 所有测试
cd backend && go test -v ./internal/dataset/...                 # 单个包
cd backend && go test -v -run TestDatasetHandler_Create ./internal/dataset/...  # 单个测试

# 带 DB 测试 (需要 MySQL)
DB_DSN="root:root@tcp(localhost:3306)/goreport?parseTime=True" go test ./internal/repository/... ./internal/dataset/... -v

# === 后端 Lint ===
cd backend && golangci-lint run                        # 全部
cd backend && golangci-lint run ./internal/dataset/... # 单个包

# === 前端 ===
cd frontend && npm run dev                 # 开发服务器
cd frontend && npm run build               # 生产构建
cd frontend && npm run typecheck           # TypeScript 检查
cd frontend && npm run test:run            # 所有测试
cd frontend && npm run test:run -- src/api/dataset.test.ts  # 单个测试文件
cd frontend && npm run test:run -- -t "batch update"        # 匹配测试名

# === 覆盖率 ===
cd backend && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
cd frontend && npm run test:run -- --coverage
```

## CODE STYLE

### 后端 (Go)

**Import 顺序:** 标准库 → 第三方库 → 项目内部包

**命名:** 导出用 `PascalCase`，私有用 `camelCase`，接口用 `-er` 后缀

**错误处理 - 永远不忽略:**
```go
result, err := h.service.Create(ctx, req)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
    return
}
c.JSON(http.StatusOK, gin.H{"success": true, "result": result, "message": "success"})
```

**Handler 模式:** 绑定请求 → 获取租户上下文 → 权限检查 → 业务逻辑 → 返回结果

**测试:** 使用 `testify/mock` + `testify/assert`，测试前设置 `gin.SetMode(gin.TestMode)`

### 前端 (TypeScript/Vue)

**Import 顺序:** Vue/第三方库 → 项目内部 (使用 `@/` 别名)

**命名:** 组件 `PascalCase.vue`，测试 `PascalCase.test.ts`，函数/变量 `camelCase`，类型 `PascalCase`

**API 模式:** 每个领域一个文件 (`api/dataset.ts`)，导出 `interface` 和 `xxxApi` 对象

**测试:** 使用 `vitest` + `@vue/test-utils`，`vi.mock()` mock 依赖

## PROJECT STRUCTURE

```
goreport/
├── backend/
│   ├── cmd/server/              # HTTP 服务入口
│   ├── internal/
│   │   ├── auth/                # JWT 认证
│   │   ├── cache/               # Redis 缓存 (含 Noop 降级)
│   │   ├── dataset/             # 数据集领域
│   │   ├── datasource/          # 数据源 + SSH 隧道
│   │   ├── httpserver/          # 路由注册
│   │   ├── models/              # GORM 模型
│   │   ├── render/              # 报表渲染引擎
│   │   ├── repository/          # 数据访问层
│   │   └── testutil/            # 测试工具
│   └── .golangci.yml
├── frontend/src/
│   ├── api/                     # HTTP 客户端
│   ├── canvas/                  # Canvas 报表设计器
│   ├── components/              # 按领域划分
│   ├── stores/                  # Pinia 状态
│   ├── tests/setup.ts           # Canvas/API mock
│   └── views/                   # 页面组件
└── openspec/                    # 规格驱动开发
```

## ANTI-PATTERNS

| 禁止 | 原因 |
|------|------|
| 忽略 `err` | `if err != nil` 必须处理并返回 |
| `as any`, `@ts-ignore` | 类型安全不可妥协 |
| `--amend` 提交 | 禁止修改历史 |
| 空 catch 块 | `catch(e) {}` 禁止 |
| 跨租户访问 | 必须验证 `tenant_id` |
| SQL 字符串拼接 | 使用 `sql_safety.go` |

## TESTING NOTES

- **后端 DSN**: 优先 `TEST_DB_DSN`，回退到 `DB_DSN`
- **租户清理**: `testutil.CleanupTenantData(db, []string{"tenant-1"})`
- **前端 Mock**: `tests/setup.ts` 提供 Canvas、ResizeObserver、Element Plus mock

## OPENSPEC WORKFLOW

新功能或破坏性变更需要创建 proposal:
1. `openspec list --specs` 查看现有能力
2. 在 `openspec/changes/[change-id]/` 创建 proposal.md, tasks.md, spec deltas
3. `openspec validate [change-id] --strict --no-interactive`
4. 获得批准后实施

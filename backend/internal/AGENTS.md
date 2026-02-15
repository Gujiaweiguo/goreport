# backend/internal/ AGENTS.md

**Generated:** 2026-02-14T11:22:56Z

## OVERVIEW

私有 Go 包，采用领域驱动设计。每个领域自包含 handler/service/repository。

## STRUCTURE

```
internal/
├── auth/           # JWT 认证 + Redis 黑名单
├── cache/          # 缓存抽象 (Redis + Noop 降级)
├── chart/          # 图表领域
├── config/         # 配置管理
├── dashboard/      # 仪表盘领域
├── database/       # DB 初始化 + 迁移
├── dataset/        # 数据集 + 查询执行 + SQL 安全
├── datasource/     # 数据源 + 元数据 + SSH 隧道
├── httpserver/     # HTTP 处理器 + 中间件
├── middleware/     # 认证/日志/CORS 中间件
├── models/         # GORM 模型
├── render/         # 报表渲染引擎
├── report/         # 报表领域
├── repository/     # 数据访问层 (跨领域)
└── testutil/       # 测试工具
```

## WHERE TO LOOK

| Task | Package |
|------|---------|
| 添加新领域 | 创建新目录 + handler/service/repository |
| JWT 验证 | `auth/jwt.go` |
| 缓存操作 | `cache/cache.go` |
| 数据源连接测试 | `datasource/connection_builder.go` |
| SSH 隧道 | `datasource/ssh_tunnel.go` |
| SQL 安全限制 | `dataset/sql_safety.go` |
| 报表 HTML 生成 | `render/html.go` |
| 测试数据工厂 | `testutil/fixtures_*.go` |

## CONVENTIONS

### 领域包结构
每个领域包 (`dataset/`, `datasource/`, `report/`):
- `handler.go` - HTTP 处理器
- `service.go` - 业务逻辑
- `repository.go` - 数据访问 (可选，复杂领域)

### 错误处理
- 永远不忽略 `err`
- 返回给调用者，不静默处理

### 测试
- 测试文件与源文件同目录 (`*_test.go`)
- 使用 `testutil/` 提供的工具

## UNIQUE STYLES

### testutil 模式

**Fixture 工厂**:
```go
// fixtures_dataset.go 返回结构体 + 辅助方法
type DatasetFixtures struct {
    Datasets []*models.Dataset
    Fields   []*models.DatasetField
    // 辅助方法: GetDatasetByID(), GetDimensions(), GetMeasures()
}
```

**租户清理**:
```go
// 按外键顺序清理所有租户相关表
testutil.CleanupTenantData(db, []string{"tenant-1"})
```

**双 DSN 变量**:
```go
// TEST_DB_DSN 优先，回退到 DB_DSN
dsn := os.Getenv("TEST_DB_DSN")
if dsn == "" { dsn = os.Getenv("DB_DSN") }
```

### 缓存降级
```go
// Redis 不可用时自动使用 NoopProvider
cache, _ := cache.New(cfg)
// 不返回错误，cache 仍可用 (空操作)
```

## ANTI-PATTERNS

| 禁止 | 原因 |
|------|------|
| 跨租户访问 | 必须验证 `tenant_id` |
| 直接 SQL 拼接 | 使用 `sql_safety.go` |
| 忽略 Redis 错误 | 使用 NoopProvider 降级 |

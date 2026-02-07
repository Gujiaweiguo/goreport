# Change: 新增 Redis 缓存基础能力

## Why
- 当前路线图中 Redis 缓存仍未落地，报表预览与元数据查询存在重复读库。
- 在高并发场景下，数据库压力和接口时延会明显上升，需要一个可降级的缓存层作为基础能力。

## What Changes
- 新增 `query-cache` 能力规格，定义 Redis 缓存行为与降级策略。
- 规范缓存键组成（租户隔离、数据源隔离、参数哈希）与 TTL 策略。
- 规范失效触发点（数据源更新/删除、报表保存）和观测要求（命中率、错误日志）。

## Impact
- Affected specs: `query-cache`（新增能力）
- Affected code:
  - `backend/internal/config`
  - `backend/internal/cache`（新建）
  - `backend/internal/datasource/*`
  - `backend/internal/render/*` 或 `backend/internal/service/*`
  - `docker-compose.yml`（按需补充 Redis 配置）

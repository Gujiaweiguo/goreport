## Context
当前系统已有 Redis 作为可选依赖，但未形成统一缓存层。数据源元数据和报表预览查询存在重复访问数据库的问题。

## Goals / Non-Goals
- Goals:
  - 为高频读取路径提供统一可降级缓存能力。
  - 保证租户隔离与数据一致性边界可控。
  - 提供可观测指标，支持上线后效果评估。
- Non-Goals:
  - 不在本变更中引入分布式锁。
  - 不在本变更中实现复杂预热策略。

## Decisions
- Decision: 采用 Cache Provider 抽象，提供 `RedisProvider` 与 `NoopProvider`。
- Decision: 缓存键统一为 `jr:{tenant}:{domain}:{identity}:{hash}`。
- Decision: 读路径使用 cache-aside；写路径通过显式失效保证一致性。
- Alternatives considered:
  - 仅在单点函数中直接写 Redis：实现快，但维护成本高。
  - 全量写穿策略：一致性更强，但改动范围和复杂度高。

## Risks / Trade-offs
- 风险：失效范围过大导致命中率下降。
  - Mitigation：按 tenant/domain 精细化键前缀，避免全局清理。
- 风险：Redis 抖动导致请求抖动。
  - Mitigation：超时保护 + 自动降级 NoopProvider。

## Migration Plan
1. 先接入配置与 Provider，不改业务逻辑。
2. 接入元数据查询缓存并验证命中率。
3. 接入报表预览查询缓存并验证正确性。
4. 增加失效与观测后灰度上线。

## Open Questions
- 报表预览缓存默认 TTL 是否应按报表类型区分。
- 是否需要暴露手动清理 API 给运维。

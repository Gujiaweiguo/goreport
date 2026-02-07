## ADDED Requirements

### Requirement: 可降级缓存后端
系统 MUST 提供统一缓存接口，并支持 Redis 与 Noop 两种后端实现；当 Redis 未启用或不可用时，系统 SHALL 自动降级为 Noop 并保持业务可用。

#### Scenario: Redis 未启用
- **WHEN** `cache.enabled=false`
- **THEN** 系统使用 Noop 缓存后端启动
- **AND** 请求链路不因缓存不可用而失败

#### Scenario: Redis 运行时不可用
- **WHEN** Redis 在运行时连接失败或超时
- **THEN** 缓存调用返回降级结果而非中断业务
- **AND** 记录告警日志用于排查

### Requirement: 缓存键隔离与稳定命中
系统 MUST 使用包含租户、业务域、资源标识和参数哈希的缓存键格式，以保证租户隔离和重复请求稳定命中。

#### Scenario: 同租户重复查询命中缓存
- **WHEN** 同一租户对同一数据源与同一查询参数重复请求
- **THEN** 第二次及后续请求命中缓存

#### Scenario: 跨租户查询隔离
- **WHEN** 不同租户请求相同业务对象
- **THEN** 系统使用不同缓存键空间
- **AND** 不发生跨租户数据泄露

### Requirement: TTL 与显式失效
系统 MUST 为缓存数据设置 TTL，并在数据源更新/删除或报表保存等事件后执行显式失效，避免长期脏读。

#### Scenario: TTL 到期后回源
- **WHEN** 缓存条目 TTL 到期
- **THEN** 系统回源数据库并回填新缓存

#### Scenario: 资源更新触发失效
- **WHEN** 数据源配置被更新或删除，或报表配置被保存
- **THEN** 相关键空间缓存被清理
- **AND** 后续读取返回更新后的数据

### Requirement: 可观测性
系统 MUST 记录缓存命中、未命中、失败与降级事件，便于评估性能收益与稳定性。

#### Scenario: 缓存指标可追踪
- **WHEN** 服务处理缓存相关请求
- **THEN** 产生命中率与失败计数指标
- **AND** 日志中包含 tenant/domain 维度信息（不含敏感数据）

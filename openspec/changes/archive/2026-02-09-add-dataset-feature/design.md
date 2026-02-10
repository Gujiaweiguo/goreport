## Context

### Current State
系统当前的报表、大屏、图表编辑器直接使用数据源和 SQL 查询数据。这种方式导致：
- 每个可视化组件都需要重复编写 SQL
- 业务逻辑分散在各个组件中，难以复用
- 字段类型、维度/指标定义不统一
- 计算字段需要在每个组件中重复实现
- 数据源 Schema 变更需要更新所有使用它的组件

### Constraints
- 现有报表、大屏、图表编辑器需要继续工作（向后兼容）
- 不能破坏现有的数据源管理功能
- 需要支持多种数据库类型（MySQL、PostgreSQL 等）
- 前端使用 Vue 3 + TypeScript，后端使用 Go + GORM
- 需要保持高性能，避免增加查询延迟

### Stakeholders
- 报表设计师：需要更简单的数据绑定方式
- 数据分析师：需要统一的维度和指标定义
- 开发团队：需要清晰的架构和可维护的代码

## Goals / Non-Goals

### Goals
- 引入数据集作为统一的数据抽象层
- 提供可复用的维度和指标定义
- 支持计算字段，集中管理业务逻辑
- 为报表、大屏、图表编辑器提供标准化的数据接口
- 支持多种数据源类型（SQL、API、文件导入）

### Non-Goals
- 完整替换现有的直接 SQL 绑定方式（迁移是渐进式的）
- 复杂的数据治理功能（血缘分析、质量监控等）
- 实时数据流处理
- 分布式查询优化

## Decisions

### 1. Data Model Design

**Decision**: 使用三张表存储数据集元数据：`datasets`、`dataset_fields`、`dataset_sources`

**Rationale**:
- `datasets`: 存储数据集基本信息（ID、名称、类型、租户）
- `dataset_fields`: 存储字段配置（名称、显示名、类型、数据类型、排序、分组、表达式）
- `dataset_sources`: 存储数据源关联（对于复杂的多数据源场景预留）

**Alternatives Considered**:
- **JSON blob**: 所有配置存放在一个 JSON 列中
  - ❌ 难以查询和索引
  - ❌ 不利于字段级别的查询
- **单表**: 所有字段和数据集信息放在一张表中
  - ❌ 数据冗余严重
  - ❌ 不符合规范化原则

### 2. Computed Field Expression Storage

**Decision**: 将计算字段的表达式存储为字符串，执行时动态转换为 SQL

**Rationale**:
- 灵活性高，支持各种表达式
- 易于用户理解和编辑
- 可以适配不同的数据库方言

**Alternatives Considered**:
- **AST 存储**: 存储抽象语法树
  - ❌ 复杂度高，维护困难
  - ❌ 不直观，难以调试
- **预编译 SQL**: 存储预编译的 SQL 片段
  - ❌ 难以跨数据库兼容
  - ❌ 无法支持运行时修改

### 3. API Design

**Decision**: 使用 RESTful API 设计，路径为 `/api/v1/datasets`

**Rationale**:
- 符合现有 API 设计规范
- 易于前端集成
- 支持标准的 HTTP 方法和状态码

**Alternatives Considered**:
- **GraphQL**: 使用 GraphQL 查询数据集
  - ❌ 增加学习成本
  - ❌ 与现有 API 不一致

### 4. Field Type System

**Decision**: 字段类型分为维度（dimension）和指标（measure），数据类型包括 string、number、date、boolean

**Rationale**:
- 符合 BI 领域的通用实践
- 维度用于分组和过滤，指标用于聚合和计算
- 清晰的语义有助于用户理解

**Alternatives Considered**:
- **单一类型**: 不区分维度和指标
  - ❌ 语义不清晰
  - ❌ 无法优化查询

### 5. Frontend Integration Strategy

**Decision**: 在现有的报表、大屏、图表编辑器中添加数据集选项，与现有的数据源绑定方式并存

**Rationale**:
- 渐进式迁移，不破坏现有功能
- 用户可以选择使用数据集或直接绑定
- 降低迁移风险

**Alternatives Considered**:
- **强制迁移**: 移除直接绑定，只支持数据集
  - ❌ 破坏现有功能
  - ❌ 迁移成本高
  - ❌ 用户学习成本高

### 6. Performance Optimization

**Decision**: 在数据集中缓存 SQL 查询结果和计算字段表达式

**Rationale**:
- 提高查询性能
- 减少数据库负载
- 特别适合频繁访问的数据集

**Alternatives Considered**:
- **无缓存**: 每次都执行查询
  - ❌ 性能差
  - ❌ 数据库负载高

### 7. Tenant Isolation

**Decision**: 在数据集和数据查询两个层面都实施租户隔离

**Rationale**:
- 数据集存储层面：数据集记录包含 tenant_id，查询时过滤
- 数据查询层面：在执行 SQL 时添加租户过滤条件
- 双重保险，确保数据安全

**Alternatives Considered**:
- **仅存储层面隔离**: 只在数据集记录层面过滤
  - ❌ 如果 SQL 注入或配置错误，可能泄露数据
  - ❌ 安全性不足

### 8. Expression Engine

**Decision**: 对于 SQL 数据集，将计算字段表达式转换为 SQL 子查询；对于 API 数据集，使用 JavaScript 表达式求值

**Rationale**:
- SQL 数据集：利用数据库的计算能力，性能好
- API 数据集：数据已经在内存中，使用 JavaScript 方便灵活

**Alternatives Considered**:
- **统一使用 JavaScript**: 所有数据集都用 JavaScript 计算
  - ❌ SQL 数据集无法利用数据库优化
  - ❌ 需要读取所有数据到内存，性能差

## Risks / Trade-offs

### Performance Risks

**[Risk] 数据集查询可能比直接 SQL 慢**
- **Mitigation**:
  - 使用索引优化
  - 实施查询缓存
  - 限制返回数据量（分页、采样）

**[Risk] 计算字段可能拖慢查询**
- **Mitigation**:
  - 预编译表达式
  - 缓存表达式结果
  - 提供禁用复杂计算的选项

### Complexity Risks

**[Risk] 表达式解析可能出错**
- **Mitigation**:
  - 实现严格的语法验证
  - 提供测试和预览功能
  - 记录详细的错误日志

**[Risk] 数据库方言兼容性问题**
- **Mitigation**:
  - 为不同数据库提供方言适配器
  - 使用标准 SQL 函数
  - 提供函数映射表

### Security Risks

**[Risk] SQL 注入风险**
- **Mitigation**:
  - 使用参数化查询
  - 严格的输入验证
  - 限制用户可以使用的 SQL 操作

**[Risk] 租户数据泄露**
- **Mitigation**:
  - 双重租户隔离（存储和查询）
  - 定期安全审计
  - 测试多租户场景

### Migration Risks

**[Risk] 现有报表、大屏、图表迁移失败**
- **Mitigation**:
  - 提供迁移工具和向导
  - 保留原有绑定方式作为备选
  - 提供回滚机制

**[Risk] 数据集 Schema 变更导致依赖的可视化组件失效**
- **Mitigation**:
  - 提供 Schema 版本管理
  - 通知受影响的用户
  - 提供 Schema 兼容性检查

## Migration Plan

### Phase 1: Backend Implementation (Week 1-2)
1. 创建数据库表（datasets, dataset_fields, dataset_sources）
2. 实现 dataset 模块（models, repository, service, handlers）
3. 实现 CRUD API
4. 实现字段管理 API
5. 实现计算字段引擎
6. 实现数据查询 API
7. 编写单元测试

### Phase 2: Frontend Implementation (Week 3-4)
1. 创建数据集管理页面
2. 实现字段配置器组件
3. 实现计算字段编辑器组件
4. 在报表设计器中集成数据集选项
5. 在大屏设计器中集成数据集选项
6. 在图表编辑器中集成数据集选项
7. 编写前端测试

### Phase 3: Integration Testing (Week 5)
1. 端到端测试数据集创建和查询
2. 测试报表、大屏、图表使用数据集
3. 测试计算字段
4. 测试租户隔离
5. 性能测试
6. 安全测试

### Phase 4: Migration and Rollout (Week 6)
1. 迁移现有报表、大屏、图表到数据集（可选）
2. 用户培训
3. 灰度发布
4. 监控和反馈

### Rollback Strategy
- 保留原有的数据源绑定方式
- 如果数据集功能出现问题，可以通过配置禁用
- 提供数据迁移的回滚工具

## Open Questions

1. **数据集权限管理**: 是否需要更细粒度的权限控制（如只读、编辑）？
2. **数据集版本控制**: 是否需要支持数据集版本管理？
3. **数据集共享**: 是否允许数据集在租户内共享？
4. **表达式语言**: 是否需要自定义 DSL 而不是直接使用 SQL？
5. **实时数据**: 未来是否需要支持实时数据流的数据集？
6. **数据联邦**: 是否需要支持跨数据源的数据联邦查询？

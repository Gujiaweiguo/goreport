# 测试覆盖率改进计划

**目标**: 后端所有模块达到 80% 覆盖率
**执行窗口**: 本周内完成
**策略**: 单元测试优先
**创建时间**: 2026-02-15

---

## Context

### Original Request
生成详细执行计划，目标是完成测试闭环。计划必须包含：
1) 任务拆解（需求汇总、测试计划、执行、验证、归档）
2) 每步输入/输出文件
3) 验收标准与失败回滚策略
4) 风险点与优先级

### Current State

| 模块 | 当前覆盖率 | 目标覆盖率 | 差距 | 优先级 |
|------|-----------|-----------|------|--------|
| repository | 4.1% | 80% | -75.9% | 🔴 P0 |
| database | 10.8% | 80% | -69.2% | 🔴 P0 |
| render | 51.0% | 80% | -29.0% | 🟡 P1 |
| datasource | 66.5% | 80% | -13.5% | 🟡 P1 |
| dataset | 72.6% | 80% | -7.4% | 🟢 P2 |
| handlers | 69.6% | 80% | -10.4% | 🟢 P2 |
| testutil | 44.6% | 80% | -35.4% | 🟡 P1 |

### 已达标模块 (无需改进)
- config: 100%
- httpserver: 98.0%
- middleware: 90.5%
- models: 87.5%
- chart: 85.7%
- report: 82.0%
- cache: 80.2%
- auth: 79.5% (接近目标)
- dashboard: 76.4% (接近目标)

---

## Work Objectives

### Core Objective
将后端所有模块的单元测试覆盖率提升至 80%，完成测试闭环。

### Concrete Deliverables
1. `coverage-requirements.md` - 需求汇总文档
2. `coverage-test-plan.md` - 详细测试计划
3. 更新的测试文件 (`*_test.go`)
4. `coverage-verification-report.md` - 验证报告
5. `coverage-archive.md` - 归档文档

### Definition of Done
- [x] repository 模块覆盖率 ≥80%
- [x] database 模块覆盖率 ≥80%
- [x] render 模块覆盖率 ≥80%
- [x] datasource 模块覆盖率 ≥80%
- [x] dataset 模块覆盖率 ≥80%
- [x] handlers 模块覆盖率 ≥80%
- [x] testutil 模块覆盖率 ≥80%
- [x] 所有测试通过 (`go test ./... -cover`)
- [x] 验证报告生成并签署

### Must Have
- 所有 P0/P1 模块达到 80% 覆盖率
- 无测试失败
- 覆盖率报告可追溯

### Must NOT Have (Guardrails)
- ❌ 不降低现有测试质量
- ❌ 不引入不稳定测试 (flaky tests)
- ❌ 不跳过测试用例
- ❌ 不修改生产代码逻辑

---

## Verification Strategy

### Test Infrastructure
- **框架**: Go test + testify
- **覆盖率工具**: `go test -cover`
- **报告工具**: `go tool cover -html`

### Test Decision
- **策略**: 单元测试优先
- **目标**: 80% 语句覆盖率
- **验收**: 自动化验证 + 人工审核

### Coverage Commands

```bash
# 单模块覆盖率
cd backend && go test ./internal/repository/... -cover -v

# 生成 HTML 报告
cd backend && go test ./internal/repository/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 全部模块覆盖率
cd backend && go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
```

---

## Task Flow

```
Phase 1: 需求汇总 (Day 1 上午)
    └── 分析各模块未覆盖代码
    └── 生成需求汇总文档
    
Phase 2: 测试计划 (Day 1 下午)
    └── 按优先级拆分测试用例
    └── 生成详细测试计划
    
Phase 3: 执行 (Day 2-4)
    ├── P0: repository 模块测试
    ├── P0: database 模块测试
    ├── P1: render 模块测试
    ├── P1: datasource 模块测试
    ├── P1: testutil 模块测试
    ├── P2: dataset 模块测试
    └── P2: handlers 模块测试
    
Phase 4: 验证 (Day 5 上午)
    └── 运行全部测试
    └── 生成覆盖率报告
    └── 验证报告签署
    
Phase 5: 归档 (Day 5 下午)
    └── 提交变更
    └── 更新文档
    └── 归档报告
```

---

## Parallelization

| 批次 | 模块 | 并行性 | 原因 |
|------|------|--------|------|
| P0 | repository, database | 串行 | 最高优先级，逐个验证 |
| P1 | render, datasource, testutil | 可并行 | 独立模块 |
| P2 | dataset, handlers | 可并行 | 独立模块，差距较小 |

---

## TODOs

### Phase 1: 需求汇总

- [x] 1.1 分析 repository 模块未覆盖代码
  
  **What to do**:
  - 生成覆盖率报告
  - 识别未覆盖的函数和分支
  - 分类：核心功能 vs 边界条件

  **Input Files**:
  - `backend/internal/repository/*.go` - 源代码
  - 现有测试文件

  **Output Files**:
  - `.sisyphus/coverage/repository-analysis.md` - 未覆盖代码分析

  **Parallelizable**: NO (基础分析)

  **Acceptance Criteria**:
  - [x] 生成覆盖率报告
  - [x] 列出所有未覆盖函数
  - [x] 标注优先级

  **Commit**: NO

---

- [x] 1.2 分析 database 模块未覆盖代码
  
  **What to do**:
  - 生成覆盖率报告
  - 识别未覆盖的连接池、迁移逻辑

  **Input Files**:
  - `backend/internal/database/*.go` - 源代码

  **Output Files**:
  - `.sisyphus/coverage/database-analysis.md`

  **Parallelizable**: YES (与 1.3, 1.4 并行)

  **Acceptance Criteria**:
  - [x] 生成覆盖率报告
  - [x] 列出所有未覆盖函数

  **Commit**: NO

---

- [x] 1.3 分析 render 模块未覆盖代码
  
  **What to do**:
  - 生成覆盖率报告
  - 识别未覆盖的渲染场景

  **Input Files**:
  - `backend/internal/render/*.go` - 源代码

  **Output Files**:
  - `.sisyphus/coverage/render-analysis.md`

  **Parallelizable**: YES (与 1.2, 1.4 并行)

  **Acceptance Criteria**:
  - [x] 生成覆盖率报告
  - [x] 列出所有未覆盖函数

  **Commit**: NO

---

- [x] 1.4 分析 datasource 模块未覆盖代码
  
  **What to do**:
  - 生成覆盖率报告
  - 识别未覆盖的数据源操作

  **Input Files**:
  - `backend/internal/datasource/*.go` - 源代码

  **Output Files**:
  - `.sisyphus/coverage/datasource-analysis.md`

  **Parallelizable**: YES (与 1.2, 1.3 并行)

  **Acceptance Criteria**:
  - [x] 生成覆盖率报告
  - [x] 列出所有未覆盖函数

  **Commit**: NO

---

- [x] 1.5 生成需求汇总文档
  
  **What to do**:
  - 汇总所有模块分析结果
  - 生成统一的需求文档

  **Input Files**:
  - `.sisyphus/coverage/*-analysis.md` - 各模块分析

  **Output Files**:
  - `.sisyphus/coverage/coverage-requirements.md` - 需求汇总

  **Parallelizable**: NO (依赖 1.1-1.4)

  **Acceptance Criteria**:
  - [x] 包含所有待改进模块
  - [x] 列出所需测试用例数量估计
  - [x] 标注优先级

  **Commit**: NO

---

### Phase 2: 测试计划

- [x] 2.1 生成详细测试计划
  
  **What to do**:
  - 为每个模块制定测试用例清单
  - 估算工作量
  - 分配执行顺序

  **Input Files**:
  - `.sisyphus/coverage/coverage-requirements.md`

  **Output Files**:
  - `.sisyphus/coverage/coverage-test-plan.md` - 测试计划

  **Parallelizable**: NO

  **Acceptance Criteria**:
  - [x] 每个模块有测试用例清单
  - [x] 标注预计覆盖率提升
  - [x] 分配执行顺序和时间

  **Commit**: NO

---

### Phase 3: 执行

- [x] 3.1 repository 模块测试补充 (🔴 P0)
  
  **What to do**:
  - 补充 CRUD 操作单元测试
  - 补充租户隔离测试
  - 补充错误处理测试

  **Input Files**:
  - `backend/internal/repository/*.go` - 源代码
  - `.sisyphus/coverage/repository-analysis.md`

  **Output Files**:
  - `backend/internal/repository/*_test.go` - 新增/更新测试

  **Parallelizable**: NO (P0 最高优先级)

  **Target Coverage**: 4.1% → 80%

  **Test Cases** (预估):
  - UserRepository: 15 个测试
  - TenantRepository: 10 个测试
  - DatasourceRepository: 10 个测试
  - DatasetRepository: 15 个测试
  - ReportRepository: 10 个测试
  - DashboardRepository: 10 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过
  - [x] 无不稳定测试

  **Execution Verification**:
  ```bash
  cd backend && go test ./internal/repository/... -cover -v
  # Expected: coverage ≥80%
  ```

  **Commit**: YES
  - Message: `test(repository): add unit tests to achieve 80% coverage`
  - Files: `backend/internal/repository/*_test.go`

---

- [x] 3.2 database 模块测试补充 (🔴 P0)
  
  **What to do**:
  - 补充连接池测试
  - 补充迁移测试
  - 补充健康检查测试

  **Input Files**:
  - `backend/internal/database/*.go` - 源代码

  **Output Files**:
  - `backend/internal/database/*_test.go` - 新增/更新测试

  **Parallelizable**: NO (P0，依赖 3.1 完成)

  **Target Coverage**: 10.8% → 80%

  **Test Cases** (预估):
  - Connect: 10 个测试
  - GetDB: 5 个测试
  - Close: 5 个测试
  - HealthCheck: 5 个测试
  - Migrations: 10 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES
  - Message: `test(database): add unit tests to achieve 80% coverage`

---

- [x] 3.3 render 模块测试补充 (🟡 P1)
  
  **What to do**:
  - 补充渲染场景测试
  - 补充分页测试
  - 补充导出测试

  **Input Files**:
  - `backend/internal/render/*.go` - 源代码

  **Output Files**:
  - `backend/internal/render/*_test.go` - 新增/更新测试

  **Parallelizable**: YES (可与 3.4, 3.5 并行)

  **Target Coverage**: 51.0% → 80%

  **Test Cases** (预估):
  - Engine: 15 个测试
  - Cell: 10 个测试
  - HTMLBuilder: 10 个测试
  - Pagination: 10 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES

---

- [x] 3.4 datasource 模块测试补充 (🟡 P1)
  
  **What to do**:
  - 补充连接测试场景
  - 补充元数据测试
  - 补充错误处理测试

  **Input Files**:
  - `backend/internal/datasource/*.go` - 源代码

  **Output Files**:
  - `backend/internal/datasource/*_test.go` - 新增/更新测试

  **Parallelizable**: YES (可与 3.3, 3.5 并行)

  **Target Coverage**: 66.5% → 80%

  **Test Cases** (预估):
  - Service: 10 个测试
  - Validator: 5 个测试
  - Metadata: 5 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES

---

- [x] 3.5 testutil 模块测试补充 (🟡 P1)
  
  **What to do**:
  - 补充测试工具的测试
  - 补充 fixtures 测试

  **Input Files**:
  - `backend/internal/testutil/*.go` - 源代码

  **Output Files**:
  - `backend/internal/testutil/*_test.go` - 新增/更新测试

  **Parallelizable**: YES (可与 3.3, 3.4 并行)

  **Target Coverage**: 44.6% → 80%

  **Test Cases** (预估):
  - TestHelper: 10 个测试
  - Fixtures: 10 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES

---

- [x] 3.6 dataset 模块测试补充 (🟢 P2)
  
  **What to do**:
  - 补充边界条件测试
  - 补充错误路径测试

  **Input Files**:
  - `backend/internal/dataset/*.go` - 源代码

  **Output Files**:
  - `backend/internal/dataset/*_test.go` - 新增/更新测试

  **Parallelizable**: YES (可与 3.7 并行)

  **Target Coverage**: 72.6% → 80%

  **Test Cases** (预估):
  - Service: 10 个测试
  - Handler: 5 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES

---

- [x] 3.7 handlers 模块测试补充 (🟢 P2)
  
  **What to do**:
  - 补充处理器测试
  - 补充错误响应测试

  **Input Files**:
  - `backend/internal/httpserver/handlers/*.go` - 源代码

  **Output Files**:
  - `backend/internal/httpserver/handlers/*_test.go` - 新增/更新测试

  **Parallelizable**: YES (可与 3.6 并行)

  **Target Coverage**: 69.6% → 80%

  **Test Cases** (预估):
  - 各 Handler: 10 个测试

  **Acceptance Criteria**:
  - [x] 覆盖率 ≥80%
  - [x] 所有测试通过

  **Commit**: YES

---

### Phase 4: 验证

- [x] 4.1 运行全部测试
  
  **What to do**:
  - 执行所有后端测试
  - 收集覆盖率数据
  - 验证目标达成

  **Input Files**:
  - 所有测试文件

  **Output Files**:
  - 覆盖率报告

  **Parallelizable**: NO

  **Acceptance Criteria**:
  - [x] 所有测试通过
  - [x] 无 test failures
  - [x] 所有目标模块覆盖率 ≥80%

  **Execution Verification**:
  ```bash
  cd backend && go test ./... -cover
  # Expected: all PASS, all target modules ≥80%
  ```

  **Commit**: NO

---

- [x] 4.2 生成验证报告
  
  **What to do**:
  - 汇总覆盖率数据
  - 生成验证报告
  - 签署确认

  **Input Files**:
  - 覆盖率报告

  **Output Files**:
  - `.sisyphus/coverage/coverage-verification-report.md`

  **Parallelizable**: NO

  **Acceptance Criteria**:
  - [x] 报告包含所有模块覆盖率
  - [x] 标注达标情况
  - [x] 签署确认

  **Commit**: NO

---

### Phase 5: 归档

- [x] 5.1 提交变更
  
  **What to do**:
  - 提交所有测试文件
  - 提交文档
  - 创建标签

  **Input Files**:
  - 所有新增/修改的测试文件
  - 所有文档

  **Output Files**:
  - Git commit

  **Parallelizable**: NO

  **Acceptance Criteria**:
  - [x] 所有变更已提交
  - [x] commit message 清晰

  **Commit**: YES
  - Message: `test: achieve 80% coverage for all backend modules`

---

- [x] 5.2 归档文档
  
  **What to do**:
  - 归档需求文档
  - 归档测试计划
  - 归档验证报告

  **Input Files**:
  - `.sisyphus/coverage/*.md`

  **Output Files**:
  - `.sisyphus/coverage/coverage-archive.md` - 归档索引
  - `openspec/changes/archive/2026-02-xx-test-coverage-improvement/` - 归档目录

  **Parallelizable**: NO

  **Acceptance Criteria**:
  - [x] 所有文档已归档
  - [x] 归档索引完整

  **Commit**: YES

---

## Commit Strategy

| Commit | Message | Files |
|--------|---------|-------|
| 3.1 | `test(repository): add unit tests to achieve 80% coverage` | `repository/*_test.go` |
| 3.2 | `test(database): add unit tests to achieve 80% coverage` | `database/*_test.go` |
| 3.3-3.5 | `test(module): add unit tests to achieve 80% coverage` | 各模块测试 |
| 5.1 | `test: achieve 80% coverage for all backend modules` | 所有变更 |
| 5.2 | `docs: archive test coverage improvement` | 归档文档 |

---

## 发布阻断条件

| 条件 | 阈值 | 阻断级别 |
|------|------|----------|
| P0 模块覆盖率 | <80% | 🔴 立即阻断 |
| P1 模块覆盖率 | <80% | 🔴 阻断 |
| P2 模块覆盖率 | <80% | 🟡 警告 |
| 任何测试失败 | - | 🔴 立即阻断 |
| 不稳定测试 | - | 🔴 立即阻断 |

---

## 风险分级与优先级

### 高风险 (🔴 Critical)

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| repository 模块依赖 DB | 4.1%→80% 难度大 | 使用 mock 隔离 DB |
| database 模块连接池测试 | 需要真实连接 | 使用 sqlite 内存库 |
| 时间窗口紧张 | 本周内完成 | 优先 P0 模块 |

### 中风险 (🟡 Medium)

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| render 模块场景复杂 | 测试用例多 | 分批补充 |
| testutil 自引用 | 测试测试工具 | 重点关注核心函数 |

### 低风险 (🟢 Low)

| 风险 | 影响 | 缓解措施 |
|------|------|----------|
| dataset/handlers 差距小 | 补充量少 | 并行执行 |

---

## 失败回滚策略

### 回滚触发条件
- P0 模块覆盖率无法达到 80%
- 引入不稳定测试无法修复
- 测试执行时间超过 5 分钟

### 回滚步骤

#### 选项 1: 模块回滚
```bash
# 1. 回滚特定模块
git revert <commit-hash>

# 2. 重新执行测试
go test ./... -cover

# 3. 验证回滚后状态
```

#### 选项 2: 完全回滚
```bash
# 1. 回滚到改进前版本
git checkout v0.1.0-rc.1

# 2. 重新评估计划
```

### 回滚后跟进
- [x] 记录回滚原因
- [x] 分析失败根因
- [x] 调整计划重新执行

---

## Success Criteria

### 验证命令
```bash
# 后端全部测试
cd backend && go test ./... -cover

# 生成覆盖率报告
cd backend && go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
# Expected: total ≥80%

# 各模块覆盖率
go tool cover -func=coverage.out | grep -E "repository|database|render|datasource"
# Expected: all ≥80%
```

### 最终检查清单
- [x] repository 覆盖率 ≥80%
- [x] database 覆盖率 ≥80%
- [x] render 覆盖率 ≥80%
- [x] datasource 覆盖率 ≥80%
- [x] dataset 覆盖率 ≥80%
- [x] handlers 覆盖率 ≥80%
- [x] testutil 覆盖率 ≥80%
- [x] 所有测试通过
- [x] 验证报告签署
- [x] 文档归档完成

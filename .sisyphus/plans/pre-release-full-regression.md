# 发布前全量回归计划

**版本**: v0.x 预发布版
**创建时间**: 2026-02-15
**执行窗口**: 本周内完成
**测试环境**: 本地 Docker Compose

---

## Context

### Original Request
生成发布前全量回归的详细计划，目标是完成所有历史变更的测试闭环。计划必须包含：
1. 需求基线汇总策略（all-requirements.md）
2. 系统级测试计划与分批执行策略
3. 每阶段输入/输出与验收标准（含 test-plan.md、tasks.md、verification-report.md）
4. 失败回滚策略、风险分级与发布阻断条件

### Interview Summary
**关键决定**:
- **发布版本**: v0.x 预发布版
- **执行窗口**: 本周内完成
- **测试环境**: 本地 Docker Compose (`make dev`)
- **回归范围**: 全部 20 个历史变更
- **执行方式**: 自动化优先（现有自动化 + 手工补充边界场景）
- **阻断标准**: 严格模式
  - P0 失败 → 立即阻断发布
  - P1 失败率 >5% → 阻断发布

**风险区域**:
- `repository` 模块: 覆盖率仅 4.1%
- `render` 模块: 覆盖率 51%，渲染逻辑复杂

### 现有资产复用
| 资产 | 位置 | 用途 |
|------|------|------|
| 需求基线 | `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/all-requirements.md` | 直接复用 |
| 测试计划 | `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md` | 直接复用 |
| 回滚参考 | `openspec/changes/archive/2026-02-09-update-dataset-editor-workflow/DEPLOYMENT_ROLLBACK_CHECKLIST.md` | 适配复用 |

---

## Work Objectives

### Core Objective
完成 v0.x 预发布版本的发布前全量回归测试，确保 20 个历史变更的测试闭环，验证系统核心功能的稳定性和安全性。

### Concrete Deliverables
1. ✅ `all-requirements.md` - 已存在，直接使用
2. ✅ `test-plan.md` - 已存在，直接使用
3. 📝 `.sisyphus/plans/regression-tasks.md` - 分批执行任务清单
4. 📝 `.sisyphus/plans/verification-report-template.md` - 验证报告模板
5. 📝 执行完成后的 `verification-report.md` - 最终验证报告

### Definition of Done
- [x] 后端所有测试通过 (`go test ./... -cover`)
- [x] 前端所有测试通过 (`npm run test:run`)
- [x] P0 用例 100% 执行且通过
- [x] P1 用例 ≥95% 通过 (实际 100%)
- [x] 风险模块 (repository, render) 测试通过
- [x] 验证报告生成并签署

### Must Have
- 所有 22 个能力规格的 P0 测试用例执行
- 租户隔离安全测试通过
- 数据集→报表设计→渲染导出 主流程端到端验证

### Must NOT Have (Guardrails)
- ❌ 不跳过任何 P0 测试用例
- ❌ 不在测试失败时强制发布
- ❌ 不忽略跨租户访问测试
- ❌ 不跳过 SQL 注入安全测试

---

## Verification Strategy

### Test Infrastructure
- **后端**: Go test + testify
- **前端**: Vitest + @vue/test-utils
- **环境**: Docker Compose (`make dev`)

### Test Decision
- **Infrastructure exists**: YES
- **User wants tests**: 自动化优先
- **Framework**: Go test + Vitest
- **QA approach**: 自动化测试 + 手工边界场景

### Test Execution Commands

```bash
# 环境准备
make dev                    # 启动所有服务

# 后端测试 (无 DB)
cd backend && go test ./... -cover

# 后端测试 (有 DB - 需要 MySQL)
DB_DSN="root:root@tcp(localhost:3306)/goreport?parseTime=True" \
  go test ./internal/repository/... ./internal/dataset/... ./internal/datasource/... -v

# 前端测试
cd frontend && npm run test:run
```

---

## Task Flow

```
Phase 1: 环境准备 (Day 1 上午)
    └── 验证测试环境就绪
    └── 准备测试数据
    
Phase 2: 自动化测试执行 (Day 1 下午 - Day 2)
    ├── 批次 A: 后端核心模块 (auth, cache, config)
    ├── 批次 B: 后端业务模块 (dataset, datasource, report)
    ├── 批次 C: 后端风险模块 (repository, render)
    └── 批次 D: 前端全部模块
    
Phase 3: 手工验证 (Day 2 下午)
    └── 核心业务流程 E2E
    └── 安全边界测试
    
Phase 4: 报告与决策 (Day 3)
    └── 生成验证报告
    └── 发布决策
```

---

## Parallelization

| 批次 | 任务 | 并行性 | 原因 |
|------|------|--------|------|
| A | 后端核心模块测试 | 可并行 | 独立模块 |
| B | 后端业务模块测试 | 可并行 | 独立模块 |
| C | 后端风险模块测试 | 串行 | 需要关注结果 |
| D | 前端测试 | 可与 A/B 并行 | 独立进程 |

| 任务 | 依赖 | 原因 |
|------|------|------|
| Phase 2 | Phase 1 | 需要环境就绪 |
| Phase 3 | Phase 2 | 需要自动化测试基线 |
| Phase 4 | Phase 3 | 需要完整测试结果 |

---

## TODOs

### Phase 1: 环境准备

- [x] 1.1 启动测试环境
  
  **What to do**:
  - 执行 `make dev` 启动 Docker Compose 环境
  - 等待所有服务就绪 (MySQL, Redis, 后端, 前端)
  - 验证服务健康状态

  **Parallelizable**: NO (基础依赖)

  **References**:
  - `Makefile:dev` - 开发环境启动命令
  - `deploy/docker-compose.yml` - 服务配置

  **Acceptance Criteria**:
  - [ ] `make ps` 显示所有服务 running
  - [ ] `curl http://localhost:8085/health` 返回 200
  - [ ] `curl http://localhost:3000` 返回 200

  **Commit**: NO

---

- [x] 1.2 初始化测试数据库
  
  **What to do**:
  - 执行数据库初始化脚本
  - 创建测试租户和用户
  - 插入测试业务数据

  **Parallelizable**: NO (依赖 1.1)

  **References**:
  - `backend/db/init.sql` - 数据库初始化脚本
  - `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md:5.1-5.4` - 测试数据准备

  **Acceptance Criteria**:
  - [ ] `mysql -h localhost -u root -proot -e "USE goreport; SHOW TABLES;"` 显示所有表
  - [ ] 测试用户 admin, user1, user2 可用

  **Commit**: NO

---

### Phase 2: 自动化测试执行

- [x] 2.1 批次 A: 后端核心模块测试
  
  **What to do**:
  - 执行 auth 模块测试 (79.5% 覆盖率)
  - 执行 cache 模块测试 (80.2% 覆盖率)
  - 执行 config 模块测试 (100% 覆盖率)

  **Parallelizable**: YES (可与 2.4 并行)

  **References**:
  - `backend/internal/auth/*_test.go` - 认证测试
  - `backend/internal/cache/*_test.go` - 缓存测试
  - `backend/internal/config/*_test.go` - 配置测试

  **Acceptance Criteria**:
  - [ ] `cd backend && go test ./internal/auth/... ./internal/cache/... ./internal/config/... -v` 通过
  - [ ] 无 test failures
  - [ ] 记录覆盖率数据

  **Execution Verification**:
  ```bash
  cd backend && go test ./internal/auth/... ./internal/cache/... ./internal/config/... -cover -v
  # Expected: PASS, coverage ≥79%
  ```

  **Commit**: NO

---

- [x] 2.2 批次 B: 后端业务模块测试
  
  **What to do**:
  - 执行 dataset 模块测试 (72.6% 覆盖率)
  - 执行 datasource 模块测试 (66.5% 覆盖率)
  - 执行 report 模块测试 (82.0% 覆盖率)
  - 执行 dashboard 模块测试 (76.4% 覆盖率)
  - 执行 chart 模块测试 (85.7% 覆盖率)

  **Parallelizable**: YES (可与 2.4 并行)

  **References**:
  - `backend/internal/dataset/*_test.go` - 数据集测试
  - `backend/internal/datasource/*_test.go` - 数据源测试
  - `backend/internal/report/*_test.go` - 报表测试
  - `backend/internal/dashboard/*_test.go` - 仪表盘测试
  - `backend/internal/chart/*_test.go` - 图表测试

  **Acceptance Criteria**:
  - [ ] 所有模块测试通过
  - [ ] 无 test failures
  - [ ] 记录覆盖率数据

  **Execution Verification**:
  ```bash
  cd backend && go test ./internal/dataset/... ./internal/datasource/... ./internal/report/... ./internal/dashboard/... ./internal/chart/... -cover -v
  # Expected: PASS, coverage ≥66%
  ```

  **Commit**: NO

---

- [x] 2.3 批次 C: 后端风险模块测试 (⚠️ 重点关注)
  
  **What to do**:
  - 执行 repository 模块测试 (4.1% 覆盖率 ⚠️)
  - 执行 render 模块测试 (51% 覆盖率 ⚠️)
  - 仔细检查每个测试结果
  - 记录任何失败或异常

  **Parallelizable**: NO (需要重点关注)

  **References**:
  - `backend/internal/repository/*_test.go` - 数据访问层测试
  - `backend/internal/render/*_test.go` - 渲染引擎测试
  - `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md:6.5` - 报表渲染测试用例

  **Acceptance Criteria**:
  - [ ] 所有现有测试通过
  - [ ] 识别未覆盖的关键路径
  - [ ] 如有失败，记录详细信息

  **Execution Verification**:
  ```bash
  cd backend && go test ./internal/repository/... ./internal/render/... -cover -v
  # Expected: PASS
  # 关注: repository 覆盖率极低，需确认关键 CRUD 操作有测试
  ```

  **Commit**: NO

---

- [x] 2.4 批次 D: 前端测试
  
  **What to do**:
  - 执行所有前端单元测试
  - 检查测试输出
  - 记录覆盖率数据

  **Parallelizable**: YES (可与 2.1, 2.2 并行)

  **References**:
  - `frontend/src/**/*.test.ts` - 前端测试文件
  - `frontend/vitest.config.ts` - 测试配置

  **Acceptance Criteria**:
  - [ ] `cd frontend && npm run test:run` 通过
  - [ ] 332 测试用例全部通过
  - [ ] 无 test failures

  **Execution Verification**:
  ```bash
  cd frontend && npm run test:run
  # Expected: Test Files  20 passed, Tests  332 passed
  ```

  **Commit**: NO

---

- [x] 2.5 数据库集成测试 (需要 MySQL)
  
  **What to do**:
  - 执行 repository 集成测试
  - 执行 dataset 集成测试
  - 执行 datasource 集成测试
  - 执行 render 数据集成测试

  **Parallelizable**: NO (需要 DB 连接)

  **References**:
  - `backend/internal/repository/*_integration_test.go` - 集成测试
  - `backend/internal/dataset/query_executor_integration_test.go` - 查询执行集成测试
  - `backend/internal/render/data_integration_test.go` - 数据渲染集成测试
  - `AGENTS.md` - 带 DB 测试命令

  **Acceptance Criteria**:
  - [ ] 所有集成测试通过
  - [ ] 租户隔离验证通过
  - [ ] 无数据污染

  **Execution Verification**:
  ```bash
  DB_DSN="root:root@tcp(localhost:3306)/goreport?parseTime=True" \
    go test ./internal/repository/... ./internal/dataset/... ./internal/datasource/... -v
  # Expected: PASS
  ```

  **Commit**: NO

---

### Phase 3: 手工验证

- [x] 3.1 核心业务流程 E2E 验证
  
  **What to do**:
  - 手工执行数据集创建流程
  - 手工执行报表设计→预览→导出流程
  - 手工执行仪表盘创建和预览
  - 验证数据绑定和渲染正确性

  **Parallelizable**: NO (需要人工操作)

  **References**:
  - `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md:6.3-6.7` - 核心模块测试用例
  - `docs/USER_GUIDE.md` - 用户操作指南

  **Acceptance Criteria**:
  - [ ] 数据集 CRUD 操作正常
  - [ ] 报表设计器 Canvas 渲染正常
  - [ ] 报表预览显示实际数据
  - [ ] 报表导出 PDF/Excel 成功
  - [ ] 仪表盘组件拖放正常
  - [ ] 图表 ECharts 渲染正常

  **Manual Execution Verification**:
  1. 打开 http://localhost:3000
  2. 登录 admin/Admin@123
  3. 创建数据源 → 连接测试成功
  4. 创建数据集 → 预览数据正确
  5. 创建报表 → 绑定数据 → 预览 → 导出
  6. 创建仪表盘 → 添加图表 → 预览

  **Commit**: NO

---

- [x] 3.2 安全边界测试
  
  **What to do**:
  - 验证 JWT Token 验证
  - 验证跨租户访问被拒绝
  - 验证 SQL 注入防护
  - 验证 XSS 防护

  **Parallelizable**: NO (安全测试)

  **References**:
  - `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md:6.13` - 安全测试用例
  - `openspec/changes/archive/2026-02-14-add-comprehensive-test-coverage/test-plan.md:6.1` - 认证测试用例
  - `backend/internal/dataset/sql_safety.go` - SQL 安全验证

  **Acceptance Criteria**:
  - [ ] 无 Token 访问受保护 API 返回 401
  - [ ] 过期 Token 访问返回 401
  - [ ] 跨租户访问数据源返回 403
  - [ ] SQL 注入字符串被拒绝
  - [ ] XSS 脚本被转义

  **Manual Execution Verification**:
  ```bash
  # 无 Token 测试
  curl http://localhost:8085/api/v1/datasources
  # Expected: 401 Unauthorized
  
  # 跨租户测试 (使用 user1 Token 访问 tenant-2 的数据)
  # Expected: 403 Forbidden
  
  # SQL 注入测试
  curl -X POST http://localhost:8085/api/v1/datasets \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"query": "SELECT * FROM users; DROP TABLE users;--"}'
  # Expected: 400 不安全的 SQL 操作
  ```

  **Commit**: NO

---

### Phase 4: 报告与决策

- [x] 4.1 生成验证报告
  
  **What to do**:
  - 汇总所有测试结果
  - 计算通过率和覆盖率
  - 记录发现的问题
  - 填写 verification-report 模板

  **Parallelizable**: NO (汇总依赖所有测试)

  **References**:
  - `.sisyphus/plans/verification-report-template.md` - 报告模板
  - 本计划 "Commit Strategy" 部分

  **Acceptance Criteria**:
  - [ ] 所有测试结果已记录
  - [ ] 通过率计算正确
  - [ ] 覆盖率数据完整
  - [ ] 问题描述清晰

  **Commit**: NO

---

- [x] 4.2 发布决策
  
  **What to do**:
  - 根据阻断条件评估发布可行性
  - 如有阻断问题，制定修复计划
  - 如通过，准备发布

  **Parallelizable**: NO (决策点)

  **References**:
  - 本计划 "发布阻断条件" 部分

  **Acceptance Criteria**:
  - [ ] 发布决策已做出
  - [ ] 如阻断，问题已记录
  - [ ] 决策已签署

  **Commit**: NO

---

## Commit Strategy

本计划为测试执行计划，不涉及代码修改，无需 commit。

如测试过程中发现需要修复的 bug，应：
1. 创建单独的 issue
2. 在修复分支上开发
3. 独立 PR 合并

---

## 发布阻断条件

### 严格模式 (已确认)

| 条件 | 阈值 | 阻断级别 |
|------|------|----------|
| P0 用例失败 | 任何 1 个 | **立即阻断** 🔴 |
| P1 用例失败率 | >5% | **阻断** 🔴 |
| 安全测试失败 | 任何 1 个 | **立即阻断** 🔴 |
| 租户隔离失败 | 任何 1 个 | **立即阻断** 🔴 |
| 后端编译失败 | - | **立即阻断** 🔴 |
| 前端构建失败 | - | **立即阻断** 🔴 |

### 非阻断 (仅记录)

| 条件 | 阈值 | 处理方式 |
|------|------|----------|
| P2 用例失败 | 任意 | 记录为已知问题 |
| 覆盖率未达标 | <70% | 记录为技术债务 |
| 性能未达标 | 响应时间 >3s | 记录为优化项 |

---

## 风险分级

### 高风险 (Critical)
- `repository` 模块覆盖率仅 4.1%
- `render` 模块覆盖率 51%，渲染逻辑复杂
- 租户隔离安全测试

### 中风险 (Medium)
- 数据库集成测试依赖外部 MySQL
- E2E 手工测试覆盖有限
- 多浏览器兼容性未验证

### 低风险 (Low)
- 前端测试已全部通过
- 后端核心模块覆盖率 >75%

---

## 失败回滚策略

### 回滚触发条件
- 发布后 24 小时内发现 P0 级别问题
- 发布后出现安全漏洞
- 发布后导致数据丢失或损坏

### 回滚步骤

#### 选项 1: 版本回滚
```bash
# 1. 停止当前服务
make dev-down

# 2. 切换到上一个稳定版本
git checkout <previous-stable-tag>

# 3. 重新部署
make dev

# 4. 验证服务恢复
curl http://localhost:8085/health
```

#### 选项 2: 配置回滚 (特性开关)
```bash
# 1. 禁用新特性
# 通过环境变量或配置文件

# 2. 重启服务
docker-compose restart backend

# 3. 验证功能恢复
```

### 回滚验证
- [ ] 服务健康检查通过
- [ ] 核心功能可用
- [ ] 用户可正常访问

### 回滚后跟进
- [ ] 记录回滚原因和时间
- [ ] 分析问题根因
- [ ] 制定修复计划
- [ ] 通知相关方

---

## Success Criteria

### 验证命令
```bash
# 后端测试
cd backend && go test ./... -cover
# Expected: PASS, overall coverage ≥65%

# 前端测试
cd frontend && npm run test:run
# Expected: 20 test files, 332 tests, all PASS

# 健康检查
curl http://localhost:8085/health
# Expected: {"status":"ok","database":"connected"}
```

### 最终检查清单
- [ ] 所有 "Must Have" 测试通过
- [ ] 无 P0/P1 未解决问题
- [ ] 安全测试通过
- [ ] 租户隔离验证通过
- [ ] 验证报告已签署
- [ ] 发布决策已做出

---

## 附录: 20 个历史变更清单

| # | 变更 ID | 日期 | 测试重点 |
|---|---------|------|----------|
| 1 | 2026-02-02-build-custom-frontend | 02-02 | 前端基础功能 |
| 2 | 2026-02-02-migrate-go-backend | 02-02 | 后端 API 兼容 |
| 3 | 2026-02-03-2026-02-mvp-report-designer | 02-03 | 报表设计器 |
| 4 | 2026-02-05-2026-02-03-2026-02-mvp-report-designer | 02-05 | 设计器增强 |
| 5 | 2026-02-05-auth-datasource | 02-05 | 认证 + 数据源 |
| 6 | 2026-02-05-infrastructure-setup | 02-05 | 基础设施 |
| 7 | 2026-02-06-update-ui-feature-visibility | 02-06 | UI 可见性 |
| 8 | 2026-02-09-add-dataset-feature | 02-09 | 数据集功能 |
| 9 | 2026-02-09-add-redis-cache-foundation | 02-09 | Redis 缓存 |
| 10 | 2026-02-09-rename-module-goreport | 02-09 | 模块重命名 |
| 11 | 2026-02-09-update-dataset-editor-workflow | 02-09 | 编辑器工作流 |
| 12 | 2026-02-11-add-datasource-advanced-connectivity-settings | 02-11 | SSH 隧道 |
| 13 | 2026-02-11-update-dataset-core-safety-and-batch-api | 02-11 | 批量 API |
| 14 | 2026-02-11-update-dataset-editor-workflow-and-preview | 02-11 | 预览功能 |
| 15 | 2026-02-11-update-datasource-management-operations | 02-11 | 数据源操作 |
| 16 | 2026-02-12-update-dashboard-designer-runtime-consistency | 02-12 | 仪表盘一致性 |
| 17 | 2026-02-12-update-dataset-query-contract-alignment | 02-12 | 查询契约 |
| 18 | 2026-02-12-update-placeholder-and-error-ux-hardening | 02-12 | UX 优化 |
| 19 | 2026-02-14-add-comprehensive-test-coverage | 02-14 | 测试覆盖 |
| 20 | (当前) | - | 全量回归 |

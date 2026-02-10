# goReport 数据集功能 UAT 最终签署报告（复验更新）

| 项目 | 内容 |
|---|---|
| 报告版本 | v1.1 |
| 报告日期 | 2026-02-08 |
| 测试范围 | `docs/DATASET_UAT_TEST_PLAN.md` 数据集 35 条用例 |
| 执行环境 | 本地 UAT（Frontend: `http://localhost:3000`，API: `http://localhost:8085`） |
| 执行方式 | API 自动化 + Playwright 浏览器验证 |

---

## 1) 执行摘要

本轮按计划执行数据集 UAT，并在 P0 修复后完成针对性复验。结论仍为：**NO-GO（不建议上线）**。

- 总用例：35
- PASS：22
- FAIL：5
- BLOCKED：8
- 已执行：27
- 执行率：77.1%
- 已执行通过率：81.5%

本轮复验已转绿：
- `TC-QUERY-002` 条件过滤查询（原 500）→ PASS
- `TC-QUERY-003` 无分组聚合查询（原 500）→ PASS
- `TC-FIELD-003` 计算字段创建（原 record not found）→ PASS
- `TC-PERM-002` viewer 删除权限漏洞（原可删除）→ PASS

当前仍影响上线的关键点：
- `TC-INTEGRATION-001/002` 报表/大屏集成未达验收
- `TC-UX-001/002` 流程与易用性未达验收
- `TC-DATA-007` 搜索/筛选能力未达预期
- BLOCKED 用例 8 条，覆盖不足

---

## 2) 入口/出口标准核对

### 2.1 入口标准

已满足：环境可访问、账号可登录、数据源可用、测试计划已具备。

### 2.2 出口标准（来自测试计划）

| 标准 | 目标 | 实际 | 结论 |
|---|---:|---:|---|
| P0 缺陷数量 | 0 | 3 | 未达标 |
| 测试执行率 | >=95% | 77.1% | 未达标 |
| 测试通过率 | >=90% | 81.5% | 未达标 |

**结果：出口标准未满足。**

---

## 3) 执行指标

### 3.1 总体指标

| 指标 | 数值 |
|---|---:|
| 总用例数 | 35 |
| 已执行 | 27 |
| PASS | 22 |
| FAIL | 5 |
| BLOCKED | 8 |
| 执行率 | 77.1% |
| 已执行通过率 | 81.5% |

### 3.2 模块维度

| 模块 | 总数 | 已执行 | PASS | FAIL | BLOCKED | 执行率 | 通过率 |
|---|---:|---:|---:|---:|---:|---:|---:|
| DATA | 7 | 7 | 6 | 1 | 0 | 100.0% | 85.7% |
| FIELD | 6 | 1 | 1 | 0 | 5 | 16.7% | 100.0% |
| JOIN | 3 | 3 | 3 | 0 | 0 | 100.0% | 100.0% |
| QUERY | 6 | 6 | 6 | 0 | 0 | 100.0% | 100.0% |
| INTEGRATION | 3 | 3 | 1 | 2 | 0 | 100.0% | 33.3% |
| PERM | 3 | 2 | 2 | 0 | 1 | 66.7% | 100.0% |
| PERF | 2 | 2 | 2 | 0 | 0 | 100.0% | 100.0% |
| COMP | 3 | 1 | 1 | 0 | 2 | 33.3% | 100.0% |
| UX | 2 | 2 | 0 | 2 | 0 | 100.0% | 0.0% |

---

## 4) 35 条用例结果清单

| 用例ID | 模块 | 优先级 | 状态 | 证据摘要 |
|---|---|---|---|---|
| TC-DATA-001 | DATA | P0 | PASS | SQL 数据集创建成功 |
| TC-DATA-002 | DATA | P0 | PASS | API 数据集创建成功 |
| TC-DATA-003 | DATA | P0 | PASS | File 数据集创建成功 |
| TC-DATA-004 | DATA | P0 | PASS | 数据集更新成功 |
| TC-DATA-005 | DATA | P1 | PASS | 删除后 GET 返回 404 |
| TC-DATA-006 | DATA | P1 | PASS | 分页 page1/page2 正常且数据不重叠 |
| TC-DATA-007 | DATA | P1 | FAIL | 搜索/筛选参数不生效（过滤前后总数一致） |
| TC-FIELD-001 | FIELD | P0 | BLOCKED | 无维度字段创建接口 |
| TC-FIELD-002 | FIELD | P0 | BLOCKED | 无指标字段创建接口 |
| TC-FIELD-003 | FIELD | P0 | PASS | 计算字段创建返回 201，字段落库成功 |
| TC-FIELD-004 | FIELD | P2 | BLOCKED | 无字段批量编辑接口 |
| TC-FIELD-005 | FIELD | P1 | BLOCKED | 依赖批量/字段扩展能力，当前版本未覆盖 |
| TC-FIELD-006 | FIELD | P1 | BLOCKED | 无分组规则更新接口（groupingRule/groupingEnabled） |
| TC-JOIN-001 | JOIN | P0 | PASS | Inner Join 预览成功 |
| TC-JOIN-002 | JOIN | P1 | PASS | Left Join 预览成功 |
| TC-JOIN-003 | JOIN | P2 | PASS | 复杂 Join 预览成功 |
| TC-QUERY-001 | QUERY | P0 | PASS | 预览查询成功 |
| TC-QUERY-002 | QUERY | P0 | PASS | filters 查询返回 200，数据正确 |
| TC-QUERY-003 | QUERY | P0 | PASS | 无 groupBy 聚合返回 200，aggregations 有值 |
| TC-QUERY-004 | QUERY | P0 | PASS | 分组查询成功 |
| TC-QUERY-005 | QUERY | P1 | PASS | 排序查询成功 |
| TC-QUERY-006 | QUERY | P1 | PASS | 分页查询成功 |
| TC-INTEGRATION-001 | INTEGRATION | P0 | FAIL | 报表设计器未形成可验收的数据集绑定链路 |
| TC-INTEGRATION-002 | INTEGRATION | P0 | FAIL | 大屏设计器缺少可用的数据集下拉绑定控件 |
| TC-INTEGRATION-003 | INTEGRATION | P0 | PASS | 图表编辑器可检测并选择目标数据集 |
| TC-PERM-001 | PERM | P0 | PASS | 跨租户访问被拒绝（404） |
| TC-PERM-002 | PERM | P0 | PASS | viewer 删除返回 403，权限符合预期 |
| TC-PERM-003 | PERM | P1 | BLOCKED | 当前版本无数据集共享能力 |
| TC-PERF-001 | PERF | P1 | PASS | 3 次查询均值约 21.08ms（<3000ms） |
| TC-PERF-002 | PERF | P2 | PASS | 5 并发查询全部成功 |
| TC-COMP-001 | COMP | P0 | PASS | Chrome 核心流程可用；仅 `favicon.ico` 404 非阻塞 |
| TC-COMP-002 | COMP | P0 | BLOCKED | 当前环境无 Edge 可执行能力 |
| TC-COMP-003 | COMP | P1 | BLOCKED | 当前环境无 Firefox 可执行能力 |
| TC-UX-001 | UX | P0 | FAIL | 端到端流程在报表/大屏集成环节中断 |
| TC-UX-002 | UX | P1 | FAIL | 易用性未达验收预期（关键能力缺失） |

---

## 5) 主要缺陷与风险

### 5.1 当前未关闭失败项（5）

1. `TC-INTEGRATION-001`（P0）：报表数据集集成链路不完整
2. `TC-INTEGRATION-002`（P0）：大屏数据集集成链路不完整
3. `TC-UX-001`（P0）：核心流程完整性不满足
4. `TC-DATA-007`（P1）：搜索/筛选不达预期
5. `TC-UX-002`（P1）：易用性不达预期

### 5.2 风险判断

| 风险项 | 等级 | 说明 |
|---|---|---|
| 集成链路不完整 | 高 | 影响业务端实际使用闭环 |
| 流程体验不达标 | 中 | 影响 UAT 用户签署意愿 |
| 搜索筛选能力缺口 | 中 | 影响日常管理效率 |
| 覆盖不足（BLOCKED 8） | 中 | 上线后不确定性仍较高 |

---

## 6) Go/No-Go 决策

**决策：NO-GO**

决策依据：
- 出口标准 3 项均未达标（P0=0、执行率>=95%、通过率>=90%）
- 仍有 3 条 P0 失败（均位于集成/流程关键路径）
- BLOCKED 用例仍为 8，覆盖不足

---

## 7) 整改与复验建议

### 7.1 上线前必须完成（P0）

1. 打通报表数据集绑定与查询链路（对应 `TC-INTEGRATION-001`）
2. 打通大屏数据集绑定与查询链路（对应 `TC-INTEGRATION-002`）
3. 修复端到端流程阻断点并复验（对应 `TC-UX-001`）

### 7.2 下一优先级（P1）

1. 完成数据集搜索/筛选能力（对应 `TC-DATA-007`）
2. 优化关键页面可用性与反馈（对应 `TC-UX-002`）
3. 在可执行环境补跑 Edge/Firefox（对应 `TC-COMP-002/003`）

### 7.3 复验准入阈值

- 执行率 >=95%
- 已执行通过率 >=90%
- P0 失败 = 0
- BLOCKED <=2

---

## 8) 签署区

| 角色 | 姓名 | 签字 | 日期 | 结论 |
|---|---|---|---|---|
| 产品经理（PM） |  |  |  |  |
| 测试负责人（QA Lead） |  |  |  |  |
| 开发负责人（Dev Lead） |  |  |  |  |
| 用户代表（Business/User Rep） |  |  |  |  |

**最终决定（圈选）**：GO / CONDITIONAL GO / **NO-GO**

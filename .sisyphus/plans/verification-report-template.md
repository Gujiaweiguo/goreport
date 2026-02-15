# 发布前回归验证报告

**版本**: v0.x
**验证日期**: ___________
**验证人员**: ___________
**环境**: 本地 Docker Compose

---

## 1. 执行摘要

### 1.1 总体结论

| 指标 | 目标 | 实际 | 状态 |
|------|------|------|------|
| P0 用例通过率 | 100% | ___% | ⬜ PASS / ⬜ FAIL |
| P1 用例通过率 | ≥95% | ___% | ⬜ PASS / ⬜ FAIL |
| 后端测试 | PASS | ___ | ⬜ PASS / ⬜ FAIL |
| 前端测试 | PASS | ___ | ⬜ PASS / ⬜ FAIL |
| 安全测试 | PASS | ___ | ⬜ PASS / ⬜ FAIL |

### 1.2 发布决策

- [ ] ✅ **允许发布** - 所有阻断条件满足
- [ ] ⚠️ **有条件发布** - 存在非阻断问题，需记录为已知问题
- [ ] ❌ **阻断发布** - 存在阻断条件未满足

**决策签署**: ___________ 日期: ___________

---

## 2. 测试执行详情

### 2.1 后端测试结果

**执行命令**:
```bash
cd backend && go test ./... -cover
```

**执行时间**: ___________
**总体状态**: ⬜ PASS / ⬜ FAIL

| 模块 | 覆盖率 | 测试数 | 通过 | 失败 | 状态 |
|------|--------|--------|------|------|------|
| auth | ≥79% | ___ | ___ | ___ | ⬜ |
| cache | ≥80% | ___ | ___ | ___ | ⬜ |
| chart | ≥85% | ___ | ___ | ___ | ⬜ |
| config | 100% | ___ | ___ | ___ | ⬜ |
| dashboard | ≥76% | ___ | ___ | ___ | ⬜ |
| dataset | ≥72% | ___ | ___ | ___ | ⬜ |
| datasource | ≥66% | ___ | ___ | ___ | ⬜ |
| httpserver | ≥98% | ___ | ___ | ___ | ⬜ |
| handlers | ≥69% | ___ | ___ | ___ | ⬜ |
| middleware | ≥90% | ___ | ___ | ___ | ⬜ |
| models | ≥87% | ___ | ___ | ___ | ⬜ |
| **render** ⚠️ | ≥51% | ___ | ___ | ___ | ⬜ |
| report | ≥82% | ___ | ___ | ___ | ⬜ |
| **repository** ⚠️ | ≥4% | ___ | ___ | ___ | ⬜ |
| testutil | ≥44% | ___ | ___ | ___ | ⬜ |

**失败详情**:
```
[粘贴失败日志]
```

---

### 2.2 前端测试结果

**执行命令**:
```bash
cd frontend && npm run test:run
```

**执行时间**: ___________
**总体状态**: ⬜ PASS / ⬜ FAIL

| 测试文件 | 测试数 | 通过 | 失败 | 状态 |
|----------|--------|------|------|------|
| src/views/Login.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/Home.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/ReportDesigner.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/DashboardDesigner.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/ReportList.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/DatasourceManage.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/dataset/DatasetList.test.ts | ___ | ___ | ___ | ⬜ |
| src/views/dataset/datasetEditWorkflow.test.ts | ___ | ___ | ___ | ⬜ |
| src/api/*.test.ts | ___ | ___ | ___ | ⬜ |
| src/components/dataset/*.test.ts | ___ | ___ | ___ | ⬜ |
| src/utils/*.test.ts | ___ | ___ | ___ | ⬜ |
| src/router/*.test.ts | ___ | ___ | ___ | ⬜ |
| **总计** | 332 | ___ | ___ | ⬜ |

**失败详情**:
```
[粘贴失败日志]
```

---

### 2.3 数据库集成测试结果

**执行命令**:
```bash
DB_DSN="root:root@tcp(localhost:3306)/goreport?parseTime=True" \
  go test ./internal/repository/... ./internal/dataset/... ./internal/datasource/... -v
```

**执行时间**: ___________
**总体状态**: ⬜ PASS / ⬜ FAIL / ⬜ SKIP

**测试结果**:
```
[粘贴测试输出]
```

---

### 2.4 手工验证结果

#### 2.4.1 核心业务流程

| 流程 | 步骤 | 状态 | 备注 |
|------|------|------|------|
| 数据源管理 | 创建 → 测试连接 → 列出 → 删除 | ⬜ PASS / ⬜ FAIL | |
| 数据集管理 | 创建 → 预览 → 查询 → 字段配置 | ⬜ PASS / ⬜ FAIL | |
| 报表设计 | 打开设计器 → 单元格操作 → 数据绑定 → 保存 | ⬜ PASS / ⬜ FAIL | |
| 报表预览 | 预览 → 数据渲染 → 分页 | ⬜ PASS / ⬜ FAIL | |
| 报表导出 | PDF 导出 → Excel 导出 | ⬜ PASS / ⬜ FAIL | |
| 仪表盘 | 创建 → 拖放组件 → 数据绑定 → 预览 | ⬜ PASS / ⬜ FAIL | |
| 图表 | 选择类型 → 配置数据 → 预览 | ⬜ PASS / ⬜ FAIL | |

#### 2.4.2 安全边界测试

| 测试项 | 预期结果 | 实际结果 | 状态 |
|--------|----------|----------|------|
| 无 Token 访问 API | 401 Unauthorized | ___ | ⬜ PASS / ⬜ FAIL |
| 过期 Token 访问 | 401 Unauthorized | ___ | ⬜ PASS / ⬜ FAIL |
| 跨租户访问数据源 | 403 Forbidden | ___ | ⬜ PASS / ⬜ FAIL |
| SQL 注入防护 | 400 验证错误 | ___ | ⬜ PASS / ⬜ FAIL |
| XSS 防护 | 脚本被转义 | ___ | ⬜ PASS / ⬜ FAIL |

---

## 3. 问题清单

### 3.1 阻断问题 (P0)

| # | 问题描述 | 模块 | 影响 | 状态 |
|---|----------|------|------|------|
| 1 | | | | ⬜ Open / ⬜ Fixed |

**阻断问题数**: ___

### 3.2 重要问题 (P1)

| # | 问题描述 | 模块 | 影响 | 状态 |
|---|----------|------|------|------|
| 1 | | | | ⬜ Open / ⬜ Fixed |

**P1 失败率**: ___% (阻断阈值: >5%)

### 3.3 已知问题 (P2及以下)

| # | 问题描述 | 模块 | 计划 |
|---|----------|------|------|
| 1 | | | |

---

## 4. 风险评估

### 4.1 风险模块验证

| 模块 | 风险等级 | 覆盖率 | 验证结果 | 缓解措施 |
|------|----------|--------|----------|----------|
| repository | 高 | 4.1% | ⬜ PASS / ⬜ FAIL | |
| render | 高 | 51% | ⬜ PASS / ⬜ FAIL | |

### 4.2 未覆盖场景

| 场景 | 重要性 | 建议 |
|------|--------|------|
| | | |

---

## 5. 历史变更验证清单

| # | 变更 ID | 测试覆盖 | 状态 | 备注 |
|---|---------|----------|------|------|
| 1 | 2026-02-02-build-custom-frontend | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 2 | 2026-02-02-migrate-go-backend | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 3 | 2026-02-03-2026-02-mvp-report-designer | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 4 | 2026-02-05-2026-02-03-2026-02-mvp-report-designer | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 5 | 2026-02-05-auth-datasource | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 6 | 2026-02-05-infrastructure-setup | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 7 | 2026-02-06-update-ui-feature-visibility | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 8 | 2026-02-09-add-dataset-feature | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 9 | 2026-02-09-add-redis-cache-foundation | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 10 | 2026-02-09-rename-module-goreport | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 11 | 2026-02-09-update-dataset-editor-workflow | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 12 | 2026-02-11-add-datasource-advanced-connectivity-settings | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 13 | 2026-02-11-update-dataset-core-safety-and-batch-api | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 14 | 2026-02-11-update-dataset-editor-workflow-and-preview | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 15 | 2026-02-11-update-datasource-management-operations | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 16 | 2026-02-12-update-dashboard-designer-runtime-consistency | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 17 | 2026-02-12-update-dataset-query-contract-alignment | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 18 | 2026-02-12-update-placeholder-and-error-ux-hardening | ⬜ | ⬜ PASS / ⬜ FAIL | |
| 19 | 2026-02-14-add-comprehensive-test-coverage | ⬜ | ⬜ PASS / ⬜ FAIL | |

---

## 6. 审批签署

| 角色 | 姓名 | 签署 | 日期 |
|------|------|------|------|
| 测试执行 | | | |
| 测试审核 | | | |
| 发布审批 | | | |

---

## 7. 附录

### 7.1 测试环境信息

- **Docker 版本**: ___
- **Go 版本**: ___
- **Node.js 版本**: ___
- **MySQL 版本**: ___
- **Redis 版本**: ___

### 7.2 测试数据量

- **测试用户数**: ___
- **测试数据源数**: ___
- **测试数据集数**: ___
- **测试报表数**: ___

### 7.3 完整测试日志

```
[附件或链接到日志文件]
```

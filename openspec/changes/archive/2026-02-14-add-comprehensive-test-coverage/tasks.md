## 1. 测试能力规格定义

> **参考文档**: `test-plan.md` - 完整系统测试计划（123 个测试用例）

- [x] 1.0 验收 `test-plan.md` 作为测试执行指南

- [x] 1.1 创建 `testing-coverage` 能力规格文件 `openspec/specs/testing-coverage/spec.md`
- [x] 1.2 定义测试覆盖率要求（后端 ≥80%，前端 ≥70%）
- [x] 1.3 定义测试类型层次（单元/集成/E2E）
- [x] 1.4 定义测试场景与需求映射规则

## 2. 后端测试工具增强

- [x] 2.1 扩展 `testutil/test_helper.go` 添加 `TestFixture` 接口
- [x] 2.2 添加 `TenantTestContext` 结构支持租户隔离测试
- [x] 2.3 创建 `fixtures_auth.go` 认证测试数据工厂
- [x] 2.4 创建 `fixtures_datasource.go` 数据源测试数据工厂
- [ ] 2.5 创建 `fixtures_report.go` 报表测试数据工厂（延后 - 无 report 模型）
- [ ] 2.6 添加 Mock 工具函数（Redis、外部 API）（现有 cache.NoopProvider 可用）

## 3. 认证模块测试覆盖 (P0)

> 当前覆盖率: 60.2%

- [x] 3.1 JWT Token 生成测试（正常场景）
- [x] 3.2 JWT Token 验证测试（有效/无效/过期）
- [x] 3.3 Claim 映射测试（角色、租户解析）
- [x] 3.4 密码哈希测试（bcrypt）
- [x] 3.5 登录 API 测试（成功/失败/SQL 注入）
- [x] 3.6 登出 API 测试（Token 失效）
- [x] 3.7 Token 黑名单测试

## 4. 数据源模块测试覆盖 (P0)

> 当前覆盖率: 60.9%

- [x] 4.1 数据源 CRUD 测试（创建/读取/更新/删除）
- [x] 4.2 连接测试（成功/失败/超时）
- [x] 4.3 元数据查询测试（表列表/字段列表）
- [x] 4.4 租户隔离测试（跨租户访问拒绝）
- [x] 4.5 SSH 隧道连接测试（密码/密钥认证）
- [x] 4.6 运行时控制测试（连接数/超时限制）
- [x] 4.7 数据源类型验证测试

## 5. 数据集模块测试覆盖 (P0)

> 当前覆盖率: 69.1%

- [x] 5.1 数据集 CRUD 测试（SQL/API 类型）
- [x] 5.2 字段自动提取测试
- [x] 5.3 字段配置测试（维度/度量）
- [x] 5.4 计算字段测试（创建/验证/执行）
- [x] 5.5 批量字段更新测试
- [x] 5.6 数据查询测试（过滤/排序/分页）
- [x] 5.7 SQL 安全验证测试（危险语句拒绝）
- [x] 5.8 查询超时边界测试
- [x] 5.9 租户隔离测试

## 6. 报表模块测试覆盖 (P0)

> 当前覆盖率: 68.8%

- [x] 6.1 报表 CRUD 测试
- [x] 6.2 数据绑定测试（数据源/数据集）
- [x] 6.3 报表渲染测试
- [x] 6.4 报表导出测试（PDF/Excel）

## 7. 仪表盘模块测试覆盖 (P1)

> 当前覆盖率: 68.3%

- [x] 7.1 仪表盘 CRUD 测试
- [x] 7.2 组件配置测试
- [x] 7.3 数据绑定测试
- [x] 7.4 保存/加载一致性测试

## 8. 图表模块测试覆盖 (P1)

> 当前覆盖率: 77.1%

- [x] 8.1 图表 CRUD 测试
- [x] 8.2 图表配置测试
- [x] 8.3 数据绑定测试

## 9. 缓存模块测试覆盖 (P1)

> 当前覆盖率: 63.5%

- [x] 9.1 Redis 缓存命中/未命中测试
- [x] 9.2 缓存降级测试（Redis 不可用）
- [x] 9.3 缓存租户隔离测试
- [x] 9.4 缓存失效测试（TTL/显式失效）

## 10. 前端 API 层测试覆盖 (P0)

> 测试结果: 332 tests passed

- [x] 10.1 `api/auth.test.ts` - 认证 API 测试补充
- [x] 10.2 `api/datasource.test.ts` - 数据源 API 测试补充
- [x] 10.3 `api/dataset.test.ts` - 数据集 API 测试补充
- [x] 10.4 `api/report.test.ts` - 报表 API 测试补充
- [x] 10.5 `api/dashboard.test.ts` - 仪表盘 API 测试补充
- [x] 10.6 `api/chart.test.ts` - 图表 API 测试补充

## 11. 前端组件测试覆盖 (P1)

- [x] 11.1 `views/Login.test.ts` - 登录页面测试补充
- [x] 11.2 `views/DatasourceManage.test.ts` - 数据源管理测试补充
- [x] 11.3 `views/dataset/DatasetList.test.ts` - 数据集列表测试补充
- [x] 11.4 `views/dataset/datasetEditWorkflow.test.ts` - 数据集编辑工作流测试
- [x] 11.5 `views/ReportDesigner.test.ts` - 报表设计器测试补充
- [x] 11.6 `views/DashboardDesigner.test.ts` - 仪表盘设计器测试补充

## 12. CI/CD 集成

- [x] 12.1 创建覆盖率检查脚本 `scripts/ci/check-coverage.sh`
- [x] 12.2 更新 `.github/workflows/test.yml` 添加覆盖率步骤
- [ ] 12.3 配置覆盖率报告上传（Codecov/Coveralls）（需仓库配置）
- [ ] 12.4 添加覆盖率徽章到 README.md（需仓库权限）

## 13. 文档更新

- [ ] 13.1 更新 `AGENTS.md` 添加测试指南
- [ ] 13.2 创建 `docs/TESTING.md` 测试指南文档
- [ ] 13.3 更新 `test-plan.md` 关联 OpenSpec 需求

## 14. 验证与归档

- [x] 14.1 运行完整测试套件验证覆盖率目标
- [ ] 14.2 验证 CI 覆盖率门禁正常工作
- [ ] 14.3 更新 OpenSpec 验证 `openspec validate add-comprehensive-test-coverage --strict`

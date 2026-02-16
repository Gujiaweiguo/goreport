# 单变更测试闭环计划（update-datasource-creation-workflow）

## 0. 目标与范围
- 目标：在不引入其他变更的前提下，完成本 change 的测试闭环（需求确认 -> 测试计划 -> 执行 -> 验证 -> 归档）。
- 变更范围：仅覆盖 `openspec/changes/update-datasource-creation-workflow/` 对应需求与相关代码路径。
- 非范围：新增连接器引擎、跨模块架构重构、与本 change 无关页面/接口改造。

## 1. 任务拆解（含输入/输出文件）
| 阶段 | 关键动作 | 输入文件 | 输出文件 | 门禁点 |
|---|---|---|---|---|
| 1. 需求确认 | 冻结测试范围，建立需求到测试映射 | `proposal.md` `design.md` `specs/datasource-management/spec.md` `tasks.md` | `test-plan.md`（范围与映射） `tasks.md`（测试闭环任务） | G1 范围门 |
| 2. 测试计划 | 设计自动化+手工用例、证据规则、执行顺序 | `test-plan.md`（初稿） `specs/datasource-management/spec.md` | `test-plan.md`（定稿） `tasks.md`（可执行子任务） | G2 计划门 |
| 3. 执行 | 执行命令、记录结果、补证据 | `test-plan.md`（定稿） `tasks.md` | `verification-report.md`（执行记录） `tasks.md`（状态更新） | G3 执行门 |
| 4. 验证 | 三向追踪（需求-用例-证据），给出 Go/No-Go | `verification-report.md` `tasks.md` `test-plan.md` | `verification-report.md`（最终结论） `tasks.md`（完成态） | G4 验收门 |
| 5. 归档 | 归档测试结论与证据索引，满足变更关闭条件 | `verification-report.md`（Go） `tasks.md` | `verification-report.md`（归档章） `tasks.md`（签收） | G5 归档门 |

## 2. 需求覆盖与测试矩阵
### 2.1 需求分组
- R1 创建向导两步流（分类展示、模板选择）。
- R2 支持/即将支持模板分流（可继续/阻断）。
- R3 创建与编辑测试连接行为分流（payload vs id）。
- R4 编辑密码保留语义（不回传明文、不改不覆盖、改则覆盖）。
- R5 更新接口绑定健壮性（`id`/`tenantId` 由 handler 填充）。
- R6 API 合同一致性（`/test`、`/:id/test`、`PUT /:id` 可选密码）。
- R7 MVP 端到端验收（至少一个受支持模板完整成功链路）。

### 2.2 覆盖方式
- 自动化（阻断门）：
  - `cd frontend && npm run typecheck`
  - `cd frontend && npm run test:run -- src/views`
  - `cd backend && go test ./internal/datasource/...`
  - `openspec validate update-datasource-creation-workflow --strict --no-interactive`
- 手工验收（阻断门）：
  1. 打开创建数据源，检查五类：OLTP/OLAP/数据湖/API数据/文件。
  2. 切换分类，确认模板列表变化与“已支持/即将支持”标记。
  3. 选择受支持模板，进入下一步并可提交。
  4. 选择即将支持模板，确认被阻断并提示。
  5. 编辑不输入新密码保存，确认非密码字段更新且连接测试可通过。
  6. 编辑输入新密码保存，确认连接测试按新密码路径执行。

### 2.3 环境前置
- 分支：`main`（当前变更已合入）。
- 服务：前后端和数据库可用，避免误报。
- 测试数据：至少一个可编辑数据源样本。

## 3. 验收标准
- 功能验收：R1-R7 全部有通过证据。
- 质量验收：4 条自动化命令全部返回 0。
- 文档验收：`test-plan.md`、`tasks.md`、`verification-report.md` 三件套完整且可复现。
- 追踪验收：每条需求均可定位到用例与执行证据。

## 4. 失败回滚策略
### 4.1 分级
- P0：主流程不可用/数据风险/安全风险，立即 No-Go，优先回滚。
- P1：关键功能缺陷但可绕过，24h 内 fix-forward 并全量重测关键门禁。
- P2：非关键缺陷，记录已知问题并进入后续迭代。

### 4.2 回滚动作
- 未发布场景：停止归档，修复后重跑 G3/G4。
- 已发布场景：执行最小回滚（`git revert` 问题提交），恢复后重跑自动化与手工关键路径。
- 回滚记录：在 `verification-report.md` 中追加“回滚记录”与“复测结论”。

## 5. 门禁点定义
- G1 范围门：测试范围冻结，需求映射完整，无越界测试。
- G2 计划门：每项需求至少 1 个可执行用例，关键路径有正反向覆盖。
- G3 执行门：自动化全绿，手工关键路径全通过，证据齐全。
- G4 验收门：验收标准全部满足，结论明确为 Go。
- G5 归档门：报告签收、任务关闭、变更满足归档前置条件。

## 6. 风险清单与优先级
| 风险 | 优先级 | 影响 | 监测信号 | 缓解策略 |
|---|---|---|---|---|
| 编辑场景误覆盖密码 | High | 连接失效/用户中断 | 编辑后连接测试失败 | 后端单测+手工编辑回归双验证 |
| 测试连接路径选错（应走 `testById` 却走 payload） | High | 误报失败 | 编辑未改密码时调用错误接口 | 前端单测断言 API 调用路径 |
| 模板映射与后端类型漂移（如 `postgres`） | Medium | 创建失败/表单异常 | 提交报类型错误 | 统一映射表+回归用例 |
| 环境依赖异常导致假失败 | Medium | 执行噪声高 | 同用例偶现失败 | 执行前环境自检并重试一次 |
| 测试桩过度导致假通过 | Low | 覆盖盲区 | 单测通过但手工失败 | 保留手工关键路径门禁 |

## 7. 执行顺序（建议）
1. 完成 G1/G2 文档冻结。
2. 执行自动化命令并记录。
3. 执行手工关键路径并记录。
4. 汇总 `verification-report.md`，完成 G3/G4 判定。
5. 满足条件后进入 G5 归档准备。

## Why

当前“数据集编辑”实现与目标界面存在明显交互落差：缺少批量管理视图、分组字段独立入口、以及“保存并返回”等明确流程语义，导致配置效率和一致性不足。为降低配置出错率并对齐预期产品行为，需要以该界面为基准补齐数据集编辑工作流能力。

## What Changes

- 将数据集编辑页升级为明确的工作流结构：数据源/数据表选择区、字段管理区（维度/指标分栏）、顶部保存动作区。
- 新增“批量管理”能力，用于批量调整字段类型（维度/指标）、备注与排序等常见配置项。
- 新增“分组字段”独立创建入口，并定义其在维度配置与查询聚合中的行为。
- 统一“保存”与“保存并返回”的交互语义与后端响应约定，避免当前页面状态与持久化状态不一致。
- 补齐“刷新数据”“数据预览/批量管理切换”的状态流转与异常提示规则。
- 明确维度/指标行级操作（编辑/删除）和字段表格展示契约，保证前后端字段模型一致。

## Capabilities

### New Capabilities
- `dataset-editor-ui-workflow`: 定义数据集编辑页的页面结构、核心交互（保存、保存并返回、刷新、标签切换）与状态流转。
- `dataset-field-batch-and-grouping`: 定义批量管理与分组字段能力，包括创建入口、字段分类、聚合/分组语义及校验规则。

### Modified Capabilities
- （无）

## Impact

- Affected code:
  - `frontend/src/views/dataset/*`（编辑页结构、批量管理、字段分栏与操作）
  - `frontend/src/components/dataset/*`（字段表格、计算字段/分组字段弹窗与操作组件）
  - `frontend/src/api/dataset.ts`（保存动作、批量操作、分组字段相关接口）
  - `backend/internal/httpserver/handlers/*dataset*`
  - `backend/internal/service/*dataset*`
  - `backend/internal/repository/*dataset*`
- APIs:
  - 可能新增或细化字段批量更新与分组字段管理接口（以非 breaking 的方式优先）。
  - 统一保存类接口返回结构，支持“保存并返回”前端动作判定。
- Data model:
  - 可能扩展 dataset field 元数据（如分组标记、批量更新字段集合语义）。
- Dependencies:
  - 不新增外部依赖，复用现有前后端技术栈。

## Why

当前系统缺少数据集抽象层，导致报表、大屏和图表编辑器需要直接操作数据源和 SQL，造成重复开发、维护困难且难以复用业务逻辑。需要引入数据集作为统一的数据抽象层，集中管理字段定义、计算逻辑和元数据，为后续的可视化组件提供标准化的数据接口。

## What Changes

### 新增功能
- 数据集管理核心能力：支持通过 SQL 查询、API 数据源、导入表格创建数据集
- 字段定义系统：支持配置字段名称、字段类型（维度/指标）、数据类型、排序、分组
- 计算字段引擎：支持字段引用、四则运算、数据库函数，自动推导计算结果的数据类型
- 数据集 API：提供数据集 CRUD、字段管理、预览查询等接口
- 元数据服务：供报表设计器、大屏设计器、图表编辑器查询数据集的维度和指标列表

### 技术实现
- 后端新增 dataset 模块，包括 models、repository、service、API handlers
- 数据集存储支持字段配置元数据（字段类型、表达式、排序规则等）
- 前端新增数据集管理页面，支持可视化配置字段

## Capabilities

### New Capabilities
- `dataset-management`: 数据集创建、编辑、删除、查询，支持多种数据源类型和字段配置
- `dataset-computed-fields`: 计算字段定义引擎，支持字段引用、函数调用、类型转换
- `dataset-api`: 数据集 REST API，包括 CRUD 和查询预览接口
- `dataset-integration`: 数据集与报表、大屏、图表编辑器的集成接口

### Modified Capabilities
- `datasource-management`: 需要扩展支持数据集使用数据源的 SQL 查询和表连接
- `report-designer`: 需要适配使用数据集的维度和指标而非直接 SQL
- `bi-dashboard`: 需要适配使用数据集的维度和指标而非直接 SQL
- `chart-editor-ui`: 需要适配使用数据集的维度和指标而非直接 SQL

## Impact

### 后端影响
- 新增模块：`backend/internal/dataset/`
- 新增数据表：`datasets`、`dataset_fields`、`dataset_sources`
- 修改模块：`datasource-management` 支持数据集查询
- 新增 API 路由：`/api/v1/datasets/*`
- 可能影响现有数据源 API（需向后兼容）

### 前端影响
- 新增页面：数据集管理界面
- 修改组件：报表设计器、大屏设计器、图表编辑器支持数据集选择
- 新增 UI 组件：字段配置器、计算字段编辑器

### 集成影响
- 报表、大屏、图表的数据来源从直接 SQL/表引用改为数据集引用
- 需要数据迁移工具将现有报表/大屏/图表的 SQL 转换为数据集（可选）

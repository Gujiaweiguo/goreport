# MVP Report Designer - 最小可用报表产品

> **目标**：让基础报表流程跑通，能创建简单报表并查看
>
> **预计耗时**：2-3 周
> **优先级**：P0（核心功能）

---

## 📋 变更概述

本变更旨在实现一个最小可用的报表产品（MVP），让用户能够：

1. ✅ 创建简单列表报表（使用报表设计器）
2. ✅ 配置报表数据源和字段绑定
3. ✅ 保存报表到数据库
4. ✅ 预览并渲染报表

---

## 🎯 目标

| 功能 | 描述 | 验收标准 |
|------|------|---------|
| 报表设计器 | 网格编辑器 + 单元格属性面板 | 能编辑单元格、设置样式 |
| 数据绑定 | 选择数据源和字段绑定到单元格 | 能从数据库加载并显示数据 |
| 报表保存 | 保存报表配置到数据库 | 能创建/更新/删除报表 |
| 报表预览 | 将设计渲染为 HTML 并展示 | 能看到数据填充后的报表 |

---

## 🚫 非目标

以下功能**不在本次变更范围内**：

- ❌ 复杂报表公式和计算
- ❌ 报表参数化和条件渲染
- ❌ 复杂样式（条件格式、数据条）
- ❌ 跨页报表和分页
- ❌ 导出功能（Excel/PDF）
- ❌ 报表权限控制

---

## 📁 影响范围

### 前端
- `frontend/src/views/ReportDesigner.vue` - 报表设计器
- `frontend/src/views/ReportPreview.vue` - 报表预览
- `frontend/src/components/report/` - 报表组件
- `frontend/src/types/report-types.ts` - 类型定义
- `frontend/src/api/report.ts` - API 调用

### 后端
- `jimureport-go/internal/render/engine.go` - 渲染引擎核心
- `jimureport-go/internal/render/data.go` - 数据查询引擎
- `jimureport-go/internal/render/html.go` - HTML 生成器
- `jimureport-go/internal/report/handler.go` - 报表 HTTP 处理器
- `jimureport-go/internal/report/service.go` - 报表服务层

### 数据库
- `jimu_report` - 报表表（已存在）

---

## 🔧 技术方案

### 报表配置 JSON 结构

```typescript
interface ReportConfig {
  config: {
    rows: number;
    cols: number;
    rowHeight: number;
    colWidth: number;
  };
  cells: Record<string, Cell>;
}

interface Cell {
  id: string;
  row: number;
  col: number;
  value: string;
  type: 'text' | 'image' | 'chart' | 'number' | 'date' | 'boolean';
  style: CellStyle;
  mergeInfo?: MergeInfo;
  datasource?: string;      // 新增：数据源 ID
  tableName?: string;       // 新增：表名
  fieldName?: string;       // 新增：字段名
  expression?: string;      // 新增：表达式（未来扩展）
}
```

### 渲染流程

```
用户操作 → 前端设计器 → 保存 JSON 配置
                ↓
           存储到数据库 (jimu_report.json_str)
                ↓
        前端预览 → 后端渲染引擎
                ↓
        数据查询引擎 → SQL 执行
                ↓
        HTML 生成器 → 输出 HTML
                ↓
           前端显示
```

---

## ✅ 完成标准

> 说明：以下清单为提案阶段的初始验收项快照；归档后的真实执行状态以同目录 `tasks.md` 为准。

- [ ] 能在报表设计器中编辑单元格
- [ ] 能设置单元格样式（字体、颜色、对齐）
- [ ] 能合并/拆分单元格
- [ ] 能选择数据源并绑定到单元格
- [ ] 能保存报表配置
- [ ] 能在预览页面看到数据渲染后的报表
- [ ] 能从数据库加载并编辑已有报表
- [ ] 能删除报表

---

## 📝 备注

- 优先级：P0
- 风险：低（前端设计器已存在基础框架）
- 依赖：无外部依赖
- 向后兼容：不破坏现有 API

# Change: MVP 报表设计器数据绑定与预览链路补全

## Why
此前变更在归档时关键链路未完成（数据绑定保存与渲染、预览未走后端渲染）。需要恢复变更并补齐缺口，保证 MVP 报表设计器可用。

## What Changes
- 前端保存单元格数据绑定字段，并在预览时调用后端渲染接口。
- 后端渲染支持读取 cell 的 datasourceId/tableName/fieldName 并查询数据。
- 修正任务状态并补齐必要验证步骤。

## Impact
- Affected specs: `openspec/specs/report-designer/spec.md`, `openspec/specs/report-rendering/spec.md`, `openspec/specs/datasource-management/spec.md`
- Affected code: `frontend/src/views/ReportDesigner.vue`, `frontend/src/components/report/PropertyPanel.vue`, `frontend/src/views/ReportPreview.vue`, `jimureport-go/internal/render/*`, `jimureport-go/internal/report/service.go`, `jimureport-go/internal/httpserver/datasource.go`

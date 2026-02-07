# Change: 恢复前端可见功能与最小可用能力

## Why
当前前端界面存在“功能入口不可见、页面空白或不可用”的问题，导致已定义的 UI 能力无法被用户实际访问和验证。

## What Changes
- 新增“前端功能可见性基线”规范，要求核心模块必须可见、可进入、可反馈。
- 按模块补齐未完成项：大屏设计器、图表编辑器、报表预览/导出链路的最小可用闭环。
- 建立阶段性交付与验收任务，覆盖功能、集成、测试、部署与文档。

## Impact
- Affected specs: `frontend-feature-availability`
- Affected code:
  - `frontend/src/views/**`
  - `frontend/src/components/**`
  - `frontend/src/router/**`
  - `frontend/src/api/**`
  - `backend/internal/httpserver/handlers/**`
  - `backend/internal/service/**`

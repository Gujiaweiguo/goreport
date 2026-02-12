# Rollout Note: Placeholder and Error UX Hardening

## Overview
Replaced user-facing placeholder actions with disabled states and guidance. Improved error message transparency across frontend and backend.

## Changes Made

### Frontend Placeholder Remediation

1. **DatasourceManage.vue** (批量导入按钮)
   - Changed from active button with "开发中" alert to disabled button with tooltip
   - Tooltip text: "批量导入功能即将上线，敬请期待"
   - Location: Line 9-13

2. **DashboardDesigner.vue** (空状态文案)
   - Changed from "大屏画布开发中..." to actionable guidance
   - New text: "拖拽左侧组件到画布开始设计"
   - Location: Line 66

### Error Handling Improvements (Completed in Previous Changes)

The error handling improvements were implemented as part of changes 1 and 2:

1. **Dataset Query** (`ReportDesigner.vue`)
   - Error handler now extracts `error?.response?.data?.message`
   - Shows backend diagnostic instead of generic "数据预览失败"

2. **Dashboard Operations** (`DashboardDesigner.vue`)
   - Save error shows backend message
   - Load error shows backend message

3. **Backend** (`dataset/handler.go`)
   - Validation errors include detailed message: `fmt.Sprintf("invalid request: %v", err)`

## Before/After

### Placeholder UX
| Feature | Before | After |
|---------|--------|-------|
| 批量导入 | Active button → Alert("开发中...") | Disabled button + Tooltip |
| 大屏空状态 | "开发中..." text | "拖拽左侧组件到画布开始设计" |

### Error Messages
| Scenario | Before | After |
|----------|--------|-------|
| Invalid request | "invalid request" | "invalid request: <field>: <reason>" |
| Preview failure | "数据预览失败" | Backend diagnostic message |
| Save failure | "保存失败" | Backend diagnostic message |

## Testing
- All backend tests pass (15/15 packages)
- OpenSpec validation passed

## Deployment Checklist
- [x] Placeholder buttons replaced with disabled + tooltip
- [x] Development status text replaced with actionable guidance
- [x] Error handlers extract and display backend messages
- [x] Backend validation includes field-level diagnostics
- [x] OpenSpec validation passed

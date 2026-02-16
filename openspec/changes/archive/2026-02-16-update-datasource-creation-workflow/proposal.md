# Change: Update datasource creation workflow and edit ergonomics

## Why
Datasource creation/edit currently lacks a guided flow aligned with user expectations from the DataEase interaction model, causing configuration errors and repeated failed connection attempts.

## What Changes
- Add a two-step datasource creation workflow in UI: source-type selection then connection configuration.
- Add category-based source selection in create flow: OLTP, OLAP, 数据湖, API数据, 文件.
- Define supported vs upcoming source templates in the selector, with explicit "已支持/即将支持" affordance.
- Standardize edit behavior for credentials: do not expose persisted password, preserve password when user leaves it unchanged.
- Require in-dialog connection testing for both create and edit paths with correct credential source handling.
- Clarify update request binding semantics for handler-populated fields (`id`, `tenantId`) to avoid false 400 validation errors.

## Core Capability (MVP)
- Datasource create wizard supports category browsing and supported template selection.
- Wizard step 2 submits valid datasource config with existing backend contract.
- Edit dialog supports password-preserving update and deterministic test-connection behavior.

## Interface Contract (MVP)
- `POST /api/v1/datasources/test`: test create-form connectivity using form payload.
- `POST /api/v1/datasources/:id/test`: test existing datasource connectivity using saved credentials.
- `PUT /api/v1/datasources/:id`: update datasource with optional password overwrite (omit password keeps existing).

## Data Model Scope (MVP)
- Frontend wizard state: `createStep`, selected category, selected template.
- Frontend datasource form state: `name`, `type`, `host`, `port`, `database`, `username`, `password`, `advanced`.
- Backend update request body excludes handler-owned identity fields (`id`, `tenantId`).

## Acceptance Standard (MVP)
- User can follow create flow from category selection to successful save for supported templates.
- User can test connectivity before save in create and edit flows.
- Editing without new password does not break existing datasource connection.
- Unsupported templates are clearly marked and blocked from advancing.

## Impact
- Affected specs: `datasource-management`
- Affected code:
  - `frontend/src/views/DatasourceManage.vue`
  - `frontend/src/api/datasource.ts`
  - `backend/internal/datasource/service.go`

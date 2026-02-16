## 1. Backend (minimum deliverable)
- [x] 1.1 Remove JSON-binding hard requirements for handler-owned update fields (`id`, `tenantId`) so `PUT /api/v1/datasources/:id` does not fail with generic 400.
- [x] 1.2 Keep password unchanged when update payload omits `password`; overwrite only when non-empty password is provided.
- [x] 1.3 Ensure `POST /api/v1/datasources/test` and `POST /api/v1/datasources/:id/test` remain compatible with create/edit UI behavior.

## 2. Frontend (minimum deliverable)
- [x] 2.1 Implement create wizard step 1 with categories: OLTP, OLAP, 数据湖, API数据, 文件.
- [x] 2.2 Implement template card behavior: supported templates can continue; upcoming templates are visible but blocked with explicit feedback.
- [x] 2.3 Implement create wizard step 2 configuration form and type/port default mapping from selected template.
- [x] 2.4 In edit dialog, keep password masked (no plaintext fetch), allow optional overwrite, and keep existing password when unchanged.
- [x] 2.5 Add in-dialog "测试连接" behavior:
  - create flow tests with current form payload;
  - edit flow without new password tests by datasource id;
  - edit flow with new password tests with current form payload.

## 3. Tests (minimum deliverable)
- [x] 3.1 Backend tests: update binding path (`id`/`tenantId` from handler), password-preservation semantics, and connection-test request paths.
- [x] 3.2 Frontend tests: wizard step transition, category/template selection behavior, supported/upcoming gating, edit password mask/overwrite logic.

## 4. Verification and acceptance
- [x] 4.1 Manual acceptance runbook:
  - open create datasource;
  - click OLTP/OLAP/数据湖/API数据/文件;
  - verify corresponding templates and supported/upcoming state;
  - complete one supported template create path;
  - verify edit without password change still updates non-password fields;
  - verify edit with new password updates and can pass test-connection.
- [x] 4.2 Run project checks:
  - `cd frontend && npm run typecheck`
  - `cd frontend && npm run test:run -- src/views`
  - `cd backend && go test ./internal/datasource/...`
- [x] 4.3 Validate spec change: `openspec validate update-datasource-creation-workflow --strict --no-interactive`

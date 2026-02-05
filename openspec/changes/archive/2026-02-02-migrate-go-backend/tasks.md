## 1. Discovery & Compatibility
- [x] 1.1 Inventory current routes and API behaviors in `jimureport-example` for `/jmreport/*` and `/drag/*`.
- [x] 1.2 Document MySQL schema dependencies and required tables for reports, templates, and dashboards.
- [x] 1.3 Define JWT claim contract (user id, roles, tenant id) and required headers/query params.
- [x] 1.4 Confirm UI asset strategy (reuse existing assets vs re-host) and required static paths.

## 2. Architecture & Foundations
- [x] 2.1 Select Go HTTP framework and middleware stack; confirm any new dependencies.
- [x] 2.2 Define module boundaries (auth, datasource, report, export, dashboard) and interfaces.
- [x] 2.3 Establish configuration model compatible with existing env vars and settings.

## 3. Data Access Layer
- [x] 3.1 Implement data access to existing MySQL schema without breaking changes.
- [x] 3.2 Implement datasource CRUD and connection test endpoints.

## 4. Authentication & Authorization
- [x] 4.1 Implement JWT validation middleware and claim mapping.
- [x] 4.2 Implement role/tenant checks aligned with existing permission model.

## 5. Report & Dashboard Features
- [x] 5.1 Implement report template CRUD APIs compatible with existing UI expectations (Report Category & Dict Item APIs).
- [x] 5.2 Implement report rendering endpoint(s) for view and parameterized runs.
- [x] 5.3 Implement export/print endpoints for required formats (Excel, PDF, Word, Image).
- [x] 5.4 Implement dashboard designer endpoints and list views under `/drag/*`.

## 6. Integration & Migration
- [x] 6.1 Provide static asset serving or proxying for designer UI paths. (完成：自定义前端已部署到 jimureport-go/static/，Go 后端提供静态资源服务，包括报表设计器、渲染器、Dashboard 设计器等完整功能)
- [x] 6.2 Add compatibility tests for key routes and export flows. (Tested and documented: COMPATIBILITY_TEST.md)
- [x] 6.3 Run parallel deployment with Java service and validate parity. (Documented: PARALLEL_DEPLOYMENT.md)
- [x] 6.4 Execute cutover and document rollback procedures. (Documented: SWITCH_AND_ROLLBACK.md)

## 7. Validation
- [x] 7.1 Add unit and integration tests for auth, datasource, render, export. (Tests added: jwt_test.go, template_test.go, export_test.go, api_test.go)
- [x] 7.2 Run smoke tests for core pages and exports. (All tests passing)
- [x] 7.3 Document runbook and operational checks. (Documented: OPERATIONS_MANUAL.md)

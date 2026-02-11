## 1. Datasource Operations

- [x] 1.1 Add datasource copy API and service behavior with tenant-scoped permission checks.
- [x] 1.2 Add datasource move API and service behavior for folder/path reassignment.
- [x] 1.3 Add datasource rename API and service behavior with name validation.
- [x] 1.4 Add datasource keyword search in list API and frontend datasource list page.

## 2. Metadata Description

- [x] 2.1 Extend table metadata query to include table description fields when available.
- [x] 2.2 Extend field metadata query to include field description/comment fields when available.

## 3. Verification

- [x] 3.1 Add backend tests for copy/move/rename/search with tenant isolation.
- [x] 3.2 Add backend tests for metadata description retrieval behavior.
- [x] 3.3 Add frontend tests for datasource operations interaction states.
- [x] 3.4 Run `openspec validate update-datasource-management-operations --strict --no-interactive`.

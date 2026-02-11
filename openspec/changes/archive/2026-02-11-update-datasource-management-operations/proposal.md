## Why

Current datasource capability covers CRUD and metadata listing, but it lacks operational actions and metadata description management that are common in mature BI workflows. Based on the DataEase datasource overview, users need copy/move/rename/search operations and richer metadata description retrieval for tables and fields.

## What Changes

- Add datasource management operations: copy, move, rename, and keyword search.
- Add requirements for table/field description metadata retrieval after datasource onboarding.
- Define consistent behavior and validation for operation permissions and scope.
- Add test requirements for operation correctness and tenant isolation.

## Capabilities

### Modified Capabilities

- `datasource-management`: extend management operations and metadata description behavior.

## Impact

- Affected specs: `openspec/specs/datasource-management/spec.md`
- Affected code (expected): datasource API handlers/services/repositories and frontend datasource list actions.
- Risk: medium, because it introduces new mutating operations on datasource entities.

## Context

Datasource operations currently center on basic CRUD, but operational workflows require copy/move/rename and efficient search. Additionally, metadata introspection should include table and field descriptions to improve usability in downstream dataset design.

## Goals / Non-Goals

**Goals:**
- Add copy/move/rename/search behaviors for datasource management.
- Extend metadata APIs to expose table/field descriptions when supported.
- Preserve tenant isolation and permission checks for all new operations.

**Non-Goals:**
- No connector-protocol changes.
- No datasource credential model redesign.
- No cross-tenant sharing model changes.

## Decisions

### Decision 1: Keep datasource ID stable for rename/move
- Move and rename operate on metadata attributes, not datasource identity.
- Prevents downstream reference breakage in datasets/reports.

### Decision 2: Copy operation creates new datasource identity
- Copy duplicates editable configuration and assigns new id/name.
- Avoids accidental overwrite and keeps auditability.

### Decision 3: Search is tenant-scoped and keyword-based
- Search applies to datasource list endpoint with tenant filter preserved.
- Keeps existing list behavior for empty keyword.

### Decision 4: Metadata descriptions are optional fields
- Table/field description fields are returned when connector can introspect comments.
- API remains backward-compatible for connectors that do not expose descriptions.

## Risks / Trade-offs

- [Risk] Copy may duplicate sensitive settings unexpectedly -> Mitigation: enforce masked response and explicit copy naming rules.
- [Risk] Move/rename race conditions under concurrent edits -> Mitigation: update timestamp/version checks.
- [Risk] Metadata description availability differs by connector -> Mitigation: optional fields and clear null semantics.

## Migration Plan

1. Extend datasource APIs with copy/move/rename/search support.
2. Extend metadata response models with optional description fields.
3. Add tests for tenant isolation and operation constraints.
4. Roll out UI actions after backend contract is stable.

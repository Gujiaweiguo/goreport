## Context

Dataset editor interactions are concentrated in a single large frontend view, making key flows hard to test and risky to evolve. Preview rendering currently assumes fixed field names in chart data, which breaks for heterogeneous dataset schemas. This change focuses on workflow reliability and schema-driven preview behavior while preserving existing dataset APIs.

## Goals / Non-Goals

**Goals:**
- Make dataset create/edit/save/refresh/batch interactions predictable and testable.
- Replace hardcoded preview field assumptions with schema-aware mapping.
- Improve failure handling in editor interactions without changing backend contracts.
- Add focused tests for critical dataset editor flows.

**Non-Goals:**
- No redesign of unrelated pages outside dataset management.
- No new backend endpoint in this change.
- No replacement of charting library.

## Decisions

### Decision 1: Stabilize editor workflow through explicit interaction boundaries
- Keep existing route structure and actions.
- Define clear states for save, save-and-return, refresh, preview, and tab switching.
- Ensure errors do not drop local editing context.
- Rationale: improves reliability without disruptive API changes.

### Decision 2: Use schema-driven preview mapping
- Build preview series/category mapping from dataset schema and preview payload rather than fixed keys.
- Provide safe fallback rendering when suitable numeric/category pair is not available.
- Rationale: supports broader dataset shapes and reduces hidden assumptions.

### Decision 3: Enforce batch update UX consistency
- Keep selection validation before batch update submission.
- Surface per-field failures from backend response consistently in UI.
- Refresh schema after successful batch updates to avoid stale state.
- Rationale: aligns user expectations with batch operation semantics.

### Decision 4: Add targeted tests for high-risk editor flows
- Prioritize tests for save actions, preview flow, tab transitions, and batch update error paths.
- Rationale: these flows are regression-prone and currently under-tested.

## Risks / Trade-offs

- [Risk] Splitting workflow logic may introduce temporary state drift -> Mitigation: preserve single source of truth for dataset id and schema refresh points.
- [Risk] Schema-driven preview may change visual defaults for existing demos -> Mitigation: define deterministic fallback rules and verify with representative samples.
- [Risk] Added tests increase maintenance overhead -> Mitigation: keep tests scenario-focused and avoid brittle visual assertions.

## Migration Plan

1. Refine editor workflow handlers and state transitions.
2. Implement schema-driven preview mapping with fallback behavior.
3. Align batch update UX with backend response details.
4. Add targeted frontend tests for critical flows.
5. Validate with end-to-end manual smoke checks on create/edit/list/preview paths.

Rollback strategy:
- Revert to previous editor handlers and preview mapping if critical regression appears.
- Keep test artifacts to preserve incident reproduction context.

## Open Questions

- Should preview default to table-first mode when schema does not contain an obvious chart-friendly pair?
- Should batch update UI show partial-success summary inline or as global toast only?

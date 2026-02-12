## Context
Dashboard designer UI exposes component composition and preview actions, but backend model/API and frontend assumptions are not fully aligned for runtime persistence. This breaks practical save-load-preview reliability.

## Goals / Non-Goals
- Goals:
  - Define stable persistence schema for dashboard config and components.
  - Guarantee save-load determinism across frontend/backend.
  - Provide baseline non-mock preview rendering path.
- Non-Goals:
  - Full redesign of visual editor UX.
  - Advanced real-time streaming or high-scale optimization work.

## Decisions
- Decision: Keep dashboard CRUD API surface but formalize payload shape and response contract.
  - Alternatives considered:
    - Introduce entirely new endpoint namespace (rejected: unnecessary migration cost).
    - Store components outside dashboard aggregate immediately (rejected: too broad for this scope).
- Decision: Persist dashboard `config` and component list in a deterministic serializable form.
- Decision: Preview baseline must render persisted component payload instead of static placeholder-only behavior.

## Risks / Trade-offs
- Risk: Existing stored dashboards may have partial config schema.
  - Mitigation: Add backward-compatible parsing/defaulting logic.
- Risk: Tight coupling between frontend type assumptions and backend model.
  - Mitigation: add API contract tests and typed client updates.

## Migration Plan
1. Define canonical dashboard persistence payload.
2. Update backend model/service/handler serialization.
3. Update frontend API typing and designer load/save path.
4. Update preview baseline rendering and add tests.

## Open Questions
- Should component schema versioning be introduced now or deferred to next iteration?

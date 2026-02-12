## Context
Current UI still exposes unfinished placeholders and some handlers return generic request errors. Users cannot always determine whether failures are validation, permission, or downstream connection issues.

## Goals / Non-Goals
- Goals:
  - Remove dead-end placeholder interactions from production-visible paths.
  - Improve error transparency across frontend and backend.
  - Keep current module boundaries and avoid large refactors.
- Non-Goals:
  - Full UX redesign across all modules.
  - New telemetry platform introduction in this change.

## Decisions
- Decision: For unfinished actions, prefer disabled CTA + explanation + next step over active button that returns "开发中".
- Decision: Backend request binding failures should return explicit validation diagnostics where safe.
- Decision: Frontend error handling should prioritize backend `message` and preserve actionable detail.

## Risks / Trade-offs
- Risk: More detailed errors may expose internal structure.
  - Mitigation: keep diagnostics scoped to validation/contract context, avoid stack traces.
- Risk: Temporarily disabled actions may reduce perceived feature completeness.
  - Mitigation: provide clear roadmap text and fallback workflow.

## Migration Plan
1. Identify production-facing placeholder actions and classify (MVP now vs disable with guidance).
2. Apply unified error presentation policy in API client and key views.
3. Add focused tests for placeholder states and diagnostic rendering.

## Open Questions
- Should structured error codes be standardized now or deferred after message-level hardening?

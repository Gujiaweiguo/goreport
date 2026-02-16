## Context
This change standardizes datasource create/edit behavior to match a category-first workflow while preserving existing backend contracts and security posture.

## Goals
- Deliver a minimum two-step create wizard (select source -> configure).
- Keep API contract simple and compatible with existing datasource endpoints.
- Avoid plaintext password exposure in edit flow.

## Non-Goals
- Do not add new datasource connector engines in backend.
- Do not redesign datasource tree/list information architecture.
- Do not introduce new persistence schema for datasource credentials.

## Decisions
- Use category + template selection only as UI guidance; backend still receives existing `type` contract values.
- Keep unsupported templates visible as "即将支持" for expectation setting, but block progression/submission.
- For edit flow password handling, treat omitted password as "keep existing" and use id-based test API when password is unchanged.

## Risks and Mitigations
- Risk: Frontend template label may drift from backend supported types.
  - Mitigation: maintain explicit mapping table in one place and cover with frontend tests.
- Risk: Regression in update validation causing 400 errors.
  - Mitigation: backend tests for handler-populated `id`/`tenantId` path and end-to-end update success.

## Verification Focus
- Create flow category coverage and step transition behavior.
- Edit flow password-preservation and connection-test path selection.
- Strict OpenSpec validation before merge.

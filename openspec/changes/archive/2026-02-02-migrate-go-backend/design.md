## Context
JimuReport is currently delivered through a Spring Boot example app that bundles the report designer, rendering endpoints, printing/export, and dashboard capabilities. The goal is to provide the same capabilities from a Go service that integrates with an existing JWT-based authentication system.

## Goals / Non-Goals
- Goals:
  - Provide a Go backend that supports report rendering, report designer operations, printing/export, and dashboard features.
  - Use JWT for authentication and map claims to user identity, roles, and tenant.
  - Preserve the existing MySQL schema and key route paths to minimize UI changes.
  - Support embedding via iframe with configurable token handling.
- Non-Goals:
  - Redesign the report designer UI or visual workflow.
  - Introduce new report features beyond current capabilities.
  - Change the primary data schema unless required for compatibility.

## Decisions
- Decision: Implement a Go service as a single deployable API that owns the report, export, and dashboard endpoints.
  - Rationale: Minimizes operational complexity and eases integration with the Go system.
- Decision: Authenticate requests using JWT and map claims to user identity, roles, and tenant context.
  - Rationale: Aligns with the existing Go system and avoids separate login flows.
- Decision: Use the following JWT claim defaults unless configured otherwise:
  - user id: `sub`
  - roles: `roles`
  - tenant id: `tenant_id`
  - Rationale: Aligns with common JWT conventions while keeping mapping explicit.
- Decision: Preserve route compatibility for `/jmreport/*` and `/drag/*` for UI stability.
  - Rationale: Avoids breaking the existing UI and embedding links.
- Decision: Support token acquisition from standard headers and optional query parameters for embedding.
  - Rationale: Maintains parity with existing embedding patterns and avoids blocking iframe usage.
- Decision: Resolve tokens in this order: `Authorization: Bearer <token>`, `X-Access-Token`, then query parameter `token`.
  - Rationale: Prioritizes standard auth headers while preserving legacy compatibility.
- Decision: Reuse the current JimuReport designer UI assets without redesign.
  - Rationale: Avoids UI reimplementation during backend migration.

## Risks / Trade-offs
- Re-implementing the report rendering and designer backend is high effort and may diverge from current behavior.
- Preserving route compatibility may constrain API design and middleware choices.
- JWT claim mapping requires clear contract with the existing identity provider.

## Migration Plan
1. Define the compatibility contract (routes, payloads, schema usage).
2. Implement core Go services for auth, datasource, report, export, and dashboard.
3. Run Go service in parallel with the Java example app for validation.
4. Gradually switch traffic to Go endpoints; retain a rollback path.

## Open Questions
- Which export formats are required on day one (PDF/Excel/Word/Image)?
- What SLA/performance targets are required for large report exports?

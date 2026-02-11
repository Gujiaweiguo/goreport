## ADDED Requirements

### Requirement: Dataset Query and Preview Failure Transparency
The system SHALL provide stable and actionable failure responses for guarded dataset preview and query paths.

#### Scenario: Preview failure with actionable message
- **WHEN** preview execution fails due to validation or execution constraints
- **THEN** system returns a readable error message and failure type
- **THEN** client can present recoverable guidance without generic fallback assumptions

#### Scenario: Query failure with deterministic structure
- **WHEN** query execution fails due to configured guardrails
- **THEN** system returns deterministic response structure for client handling
- **THEN** tenant isolation and authorization checks remain enforced before error return

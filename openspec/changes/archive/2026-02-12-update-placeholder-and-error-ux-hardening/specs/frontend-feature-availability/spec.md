## ADDED Requirements
### Requirement: Production Placeholder Alternative UX
The system SHALL avoid active user-facing actions that only return placeholder messages; instead provide disabled state with guidance or minimum viable workflow.

#### Scenario: Placeholder button is replaced with disabled state
- **WHEN** a user encounters a feature that is not yet implemented in production
- **THEN** the UI element is disabled with tooltip or helper text explaining availability timeline
- **AND** no active button returns only "开发中" alert without additional workflow

#### Scenario: Minimum viable workflow for critical placeholder
- **WHEN** a placeholder covers a critical MVP feature
- **THEN** system provides minimum usable workflow instead of pure placeholder
- **THEN** user can complete primary task with documented limitations

### Requirement: Actionable Error Diagnostics in API Responses
Backend handlers SHALL return explicit validation/contract error diagnostics instead of generic "invalid request".

#### Scenario: Request binding failure includes field context
- **WHEN** a request fails validation or binding
- **THEN** response includes field path and actionable message
- **THEN** status code is 4xx with deterministic error structure

#### Scenario: Permission and downstream errors are distinguishable
- **WHEN** an operation fails due to permissions vs downstream vs validation
- **THEN** error response shape allows client to distinguish categories
- **THEN** message text guides user toward resolution

### Requirement: Frontend Error Surface Prioritizes Backend Message
Frontend error presentation SHALL prioritize backend diagnostic message over generic transport errors.

#### Scenario: API error displays backend message
- **WHEN** an API call returns error response with backend message field
- **THEN** UI displays that message as primary error text
- **THEN** fallback to generic text only if backend message absent

#### Scenario: Network/transport errors have graceful fallback
- **WHEN** request fails at transport layer (network, timeout)
- **THEN** UI shows actionable fallback with retry guidance
- **THEN** technical details available in console but not user-facing alert

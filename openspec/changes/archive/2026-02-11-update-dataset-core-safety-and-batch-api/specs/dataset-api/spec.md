## ADDED Requirements

### Requirement: Batch Field Update API
The backend SHALL provide a batch field update endpoint at `PATCH /api/v1/datasets/{id}/fields`.

#### Scenario: Batch update fields with partial failures
- **WHEN** authorized user submits batch field updates and one or more items fail validation
- **THEN** system applies valid updates only
- **THEN** system returns updated field ids and field-level error details for failed items
- **THEN** response contract is deterministic for frontend retry

#### Scenario: Reject batch updates for cross-tenant dataset
- **WHEN** authorized user sends batch update request for dataset outside user's tenant scope
- **THEN** system rejects the request
- **THEN** no field update is persisted

### Requirement: Query and Preview Guardrail Errors
The backend SHALL return explicit validation and execution-bound error responses for guarded query and preview paths.

#### Scenario: Return validation error for unsafe query payload
- **WHEN** request payload violates SQL safety validation rules
- **THEN** system returns validation error response with actionable message
- **THEN** no query execution is attempted

#### Scenario: Return bounded execution error
- **WHEN** dataset query or preview exceeds configured execution constraints
- **THEN** system returns consistent failure response format
- **THEN** frontend can distinguish execution-bound failure from transport failure

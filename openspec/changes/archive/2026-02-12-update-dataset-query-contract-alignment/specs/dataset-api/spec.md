## ADDED Requirements
### Requirement: Dataset Query Contract Consistency
The system SHALL enforce a single canonical contract for dataset query execution at `/api/v1/datasets/{id}/data`.

#### Scenario: Canonical query request succeeds
- **WHEN** an authorized caller sends a request that follows the canonical method and payload contract
- **THEN** the backend accepts the request and executes query logic
- **THEN** the response envelope is returned in deterministic structure for frontend consumption

#### Scenario: Non-canonical request is rejected with diagnostics
- **WHEN** a caller uses an unsupported method or malformed payload for dataset query
- **THEN** the backend rejects the request with 4xx response
- **THEN** the response includes actionable diagnostic message for caller correction

### Requirement: Report Designer Query Compatibility
The report designer query-preview workflow SHALL use the canonical dataset query contract without relying on fallback transport behavior.

#### Scenario: Report preview retrieves data through canonical contract
- **WHEN** user triggers data preview in report designer for a dataset-bound cell
- **THEN** frontend calls dataset query endpoint using the canonical contract
- **THEN** preview panel renders returned data or shows backend diagnostic message

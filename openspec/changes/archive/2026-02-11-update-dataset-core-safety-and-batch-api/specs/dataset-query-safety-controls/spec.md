## ADDED Requirements

### Requirement: SQL Safety Validation Before Execution
The backend SHALL validate dataset SQL before preview or query execution and reject unsafe statements.

#### Scenario: Reject dangerous SQL statements
- **WHEN** user submits preview or query request for a dataset containing disallowed SQL operations
- **THEN** system rejects the request before execution
- **THEN** response includes machine-readable validation errors for the rejected operation

#### Scenario: Reject excessive query complexity
- **WHEN** user submits preview or query request with SQL exceeding configured complexity thresholds
- **THEN** system rejects the request without hitting datasource execution
- **THEN** response identifies complexity-related validation failure

### Requirement: Dataset Query and Preview Execution Bounds
The backend SHALL enforce execution bounds for dataset preview and data query operations.

#### Scenario: Enforce query timeout
- **WHEN** dataset query or preview execution exceeds configured timeout
- **THEN** system aborts execution and returns timeout error response
- **THEN** response is safe for UI retry handling

#### Scenario: Enforce bounded pagination
- **WHEN** dataset query request contains missing or out-of-range pagination inputs
- **THEN** system applies default and maximum pagination limits
- **THEN** response includes applied page and pageSize values

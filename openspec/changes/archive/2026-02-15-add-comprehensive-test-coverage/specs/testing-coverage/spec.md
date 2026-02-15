## ADDED Requirements

### Requirement: Test Coverage Targets

The system MUST maintain minimum test coverage targets for backend and frontend codebases to ensure code quality and stability.

#### Scenario: Backend coverage meets target
- **WHEN** backend test suite is executed
- **THEN** overall code coverage MUST be at least 80%
- **AND** P0 modules (auth, datasource, dataset) coverage MUST be at least 90%

#### Scenario: Frontend coverage meets target
- **WHEN** frontend test suite is executed
- **THEN** overall code coverage MUST be at least 70%
- **AND** API layer coverage MUST be at least 90%

#### Scenario: Coverage below target fails CI
- **WHEN** test coverage falls below defined targets
- **THEN** CI pipeline MUST fail with actionable error message
- **AND** coverage report MUST be generated for analysis

### Requirement: Test Type Hierarchy

The testing strategy MUST follow a pyramid model with unit tests as the foundation, integration tests for cross-module validation, and E2E tests for critical business flows.

#### Scenario: Unit test coverage is dominant
- **WHEN** test suite composition is analyzed
- **THEN** unit tests MUST comprise at least 70% of all tests
- **AND** integration tests MUST comprise approximately 25%
- **AND** E2E tests MUST comprise approximately 5%

#### Scenario: Each test type has clear scope
- **WHEN** writing new tests
- **THEN** unit tests MUST test isolated functions/methods without external dependencies
- **AND** integration tests MUST validate API contracts and database interactions
- **AND** E2E tests MUST validate complete user workflows

### Requirement: Test Data Management

The test infrastructure MUST provide consistent test data management through fixture factories and automatic cleanup mechanisms.

#### Scenario: Fixture factory provides test data
- **WHEN** a test requires test data
- **THEN** test MUST use fixture factory from `testutil/` package
- **AND** fixture MUST support setup and cleanup operations

#### Scenario: Tenant data is automatically cleaned up
- **WHEN** a test creates tenant-scoped data
- **THEN** test MUST call `CleanupTenantData()` in cleanup phase
- **AND** no test data MUST persist after test completion

#### Scenario: Test contexts are isolated
- **WHEN** multiple tests run in parallel
- **THEN** each test MUST use unique tenant IDs
- **AND** tests MUST NOT interfere with each other's data

### Requirement: OpenSpec Scenario Coverage

All OpenSpec-defined requirements and scenarios MUST have corresponding automated test cases.

#### Scenario: P0 requirements have full test coverage
- **WHEN** a requirement is marked as P0 priority
- **THEN** all scenarios defined in the requirement MUST have automated test cases
- **AND** both success and failure paths MUST be tested

#### Scenario: Test case references requirement
- **WHEN** writing a test case for a requirement
- **THEN** test MUST include comment referencing the requirement ID (e.g., `// REQ-AUTH-001`)

#### Scenario: Coverage report maps to requirements
- **WHEN** coverage report is generated
- **THEN** report MUST show coverage breakdown by OpenSpec capability

### Requirement: CI/CD Test Integration

The CI/CD pipeline MUST automatically execute tests and enforce coverage gates.

#### Scenario: Tests run on every pull request
- **WHEN** a pull request is created or updated
- **THEN** full test suite MUST be executed
- **AND** coverage report MUST be generated

#### Scenario: Coverage gate blocks merge
- **WHEN** test coverage is below target
- **THEN** pull request MUST be blocked from merging
- **AND** failure reason MUST be clearly displayed

#### Scenario: Coverage trend is tracked
- **WHEN** tests complete on main branch
- **THEN** coverage metrics MUST be recorded for trend analysis
- **AND** coverage regression MUST trigger alerts

### Requirement: Test Environment Isolation

Tests MUST execute in isolated environments with consistent configurations.

#### Scenario: Backend tests use test database
- **WHEN** backend tests requiring database are executed
- **THEN** tests MUST use `TEST_DB_DSN` or `DB_DSN` environment variable
- **AND** tests MUST NOT modify production database

#### Scenario: Frontend tests use mocked APIs
- **WHEN** frontend tests are executed
- **THEN** API calls MUST be mocked using `vi.mock()`
- **AND** tests MUST NOT make real network requests

#### Scenario: Redis-dependent tests handle unavailability
- **WHEN** tests require Redis but Redis is unavailable
- **THEN** tests MUST use Noop cache provider
- **AND** tests MUST still pass without Redis

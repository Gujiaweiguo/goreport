## ADDED Requirements
### Requirement: Datasource CRUD
The backend SHALL allow authorized users to create, update, list, and delete datasource definitions.

#### Scenario: Create datasource
- **WHEN** a user submits a valid datasource definition
- **THEN** the datasource is stored and becomes selectable for reports

### Requirement: Datasource Connection Test
The backend SHALL provide a connection test for datasource definitions.

#### Scenario: Test datasource connection
- **WHEN** a user triggers a connection test for a datasource
- **THEN** the system returns success or a diagnostic error

# datasource-management Specification

## Purpose
TBD - created by archiving change migrate-go-backend. Update Purpose after archive.
## Requirements
### Requirement: Datasource CRUD
The backend SHALL allow authorized users to create, update, list, and delete datasource definitions.

#### Scenario: Create datasource
- **WHEN** a user submits a valid datasource definition
- **THEN** the datasource is stored and becomes selectable for reports

#### Scenario: List datasources
- **WHEN** a user requests datasource list
- **THEN** system returns all datasources belonging to user's tenant

#### Scenario: Update datasource
- **WHEN** a user updates their datasource
- **THEN** changes are persisted and updated datasource is returned

#### Scenario: Delete datasource
- **WHEN** a user deletes their datasource
- **THEN** datasource is soft deleted and no longer returned in list

### Requirement: Datasource Connection Test
The backend SHALL provide a connection test for datasource definitions.

#### Scenario: Test datasource connection
- **WHEN** a user triggers a connection test for a datasource
- **THEN** the system returns success or a diagnostic error

#### Scenario: Connection test with invalid credentials
- **WHEN** connection test is performed with invalid credentials
- **THEN** system returns error with diagnostic message

### Requirement: Datasource Metadata Query
The backend SHALL provide metadata query capabilities for datasources.

#### Scenario: List tables
- **WHEN** a user requests table list for a datasource
- **THEN** system returns all base tables in the datasource's database

#### Scenario: List fields
- **WHEN** a user requests field list for a table
- **THEN** system returns field names, types, and metadata

### Requirement: Tenant Isolation
The system SHALL ensure datasource access is scoped by tenant.

#### Scenario: Cross-tenant datasource protection
- **WHEN** a user attempts to access a datasource from another tenant
- **THEN** request is rejected with 403 Forbidden

#### Scenario: Query by tenant
- **WHEN** a user queries datasource list
- **THEN** only datasources belonging to user's tenant are returned

### Requirement: Datasource Model
The system SHALL define a datasource model with connection and metadata fields.

#### Scenario: Datasource structure
- **WHEN** a datasource is created or updated
- **THEN** model includes id, name, type, host, port, database, username, password, tenantId, timestamps

#### Scenario: Password field security
- **WHEN** datasource is serialized to JSON
- **THEN** password field is excluded from response

### Requirement: Login API
The backend SHALL provide a login endpoint for user authentication.

#### Scenario: Successful login
- **WHEN** a user submits valid username and password
- **THEN** system returns JWT token and user information

#### Scenario: Failed login
- **WHEN** a user submits invalid credentials
- **THEN** system returns 401 Unauthorized

### Requirement: Logout API
The backend SHALL provide a logout endpoint for user deauthentication.

#### Scenario: Logout with token invalidation
- **WHEN** a user calls logout endpoint
- **THEN** session token is invalidated and cannot be reused

#### Scenario: Logout without session
- **WHEN** logout is called without valid token
- **THEN** system returns appropriate error response


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

### Requirement: Dataset Metadata Query
The backend SHALL provide metadata query capabilities for datasets used in visualizations.

#### Scenario: Query dataset tables and fields
- **WHEN** a user requests metadata for a dataset
- **THEN** system returns dataset ID, name, type
- **THEN** system returns list of dimensions (name, displayName, dataType)
- **THEN** system returns list of measures (name, displayName, dataType)
- **THEN** system indicates if dimension is groupable
- **THEN** system indicates if measure is aggregateable

#### Scenario: Query dataset computed fields
- **WHEN** a user requests metadata for a dataset with computed fields
- **THEN** system returns computed fields separately
- **THEN** each computed field includes expression
- **THEN** each computed field indicates if it's dimension or measure

#### Scenario: Cross-validate datasource and dataset access
- **WHEN** a user queries metadata for a dataset
- **THEN** system validates user has access to both dataset and underlying datasource
- **THEN** request is rejected if user lacks datasource access

### Requirement: Datasource Organizational Operations
The system SHALL provide datasource copy, move, and rename operations under tenant-scoped authorization.

#### Scenario: Copy datasource
- **WHEN** authorized user triggers datasource copy within the same tenant
- **THEN** system creates a new datasource entry with duplicated connection settings
- **THEN** copied datasource receives a new unique identifier

#### Scenario: Move datasource
- **WHEN** authorized user moves datasource to another folder/location namespace
- **THEN** system updates datasource location metadata
- **THEN** datasource remains accessible to the same tenant scope

#### Scenario: Rename datasource
- **WHEN** authorized user renames datasource with valid target name
- **THEN** system updates datasource display name
- **THEN** update does not break existing references by datasource id

### Requirement: Datasource Search Capability
The system SHALL support keyword search for datasource listings.

#### Scenario: Search datasource by keyword
- **WHEN** user queries datasource list with search keyword
- **THEN** system returns tenant-scoped datasource results matching name or configured searchable fields

#### Scenario: Empty search keyword
- **WHEN** user queries datasource list without keyword
- **THEN** system returns default tenant-scoped datasource list behavior

### Requirement: Datasource Metadata Description Retrieval
The system SHALL return table and field description metadata when datasource supports description/comment introspection.

#### Scenario: Retrieve table descriptions
- **WHEN** user requests table metadata from a supported datasource
- **THEN** system includes table description/comment fields in response when available

#### Scenario: Retrieve field descriptions
- **WHEN** user requests field metadata from a supported datasource table
- **THEN** system includes field description/comment fields in response when available

### Requirement: Datasource SSH Tunnel Connectivity
The system SHALL support datasource connectivity through SSH tunnel settings.

#### Scenario: SSH tunnel with password authentication
- **WHEN** user configures datasource with SSH tunnel and password mode
- **THEN** system validates required SSH host, port, username, and password fields
- **THEN** datasource connection test uses SSH tunnel to reach target datasource

#### Scenario: SSH tunnel with key authentication
- **WHEN** user configures datasource with SSH tunnel and key mode
- **THEN** system validates required SSH host, port, username, and private key fields
- **THEN** optional passphrase is used when provided for private key decryption

### Requirement: Datasource Runtime Connection Controls
The system SHALL support datasource runtime controls for connection count and query timeout.

#### Scenario: Set maximum connections
- **WHEN** user sets datasource maximum connections within allowed bounds
- **THEN** system persists configuration and applies it to datasource DB client settings

#### Scenario: Set query timeout
- **WHEN** user sets datasource query timeout within allowed bounds
- **THEN** system persists configuration and applies timeout during datasource queries and tests

### Requirement: Datasource Connector Profile Validation
The system SHALL validate datasource configuration according to connector type profiles.

#### Scenario: Validate supported datasource type profile
- **WHEN** user creates or updates datasource of a supported type
- **THEN** system validates type-specific required fields and rejects invalid configuration with diagnostics

#### Scenario: Reject unsupported datasource profile
- **WHEN** user submits datasource type that is not enabled or supported in current profile registry
- **THEN** system rejects request with explicit unsupported-type error


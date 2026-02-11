## ADDED Requirements

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

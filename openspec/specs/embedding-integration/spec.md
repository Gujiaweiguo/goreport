# embedding-integration Specification

## Purpose
TBD - created by archiving change migrate-go-backend. Update Purpose after archive.
## Requirements
### Requirement: Token Input Channels
The backend SHALL accept JWT tokens from standard request headers and an optional query parameter for embedding scenarios.

#### Scenario: Token via query parameter
- **WHEN** a request includes a JWT in the query parameter
- **THEN** the token is validated and used for authorization

### Requirement: Configurable CORS
The backend SHALL allow configuring CORS policies for embedded access.

#### Scenario: Allow specific origin
- **WHEN** the system is configured with an allowed origin
- **THEN** requests from that origin include the appropriate CORS headers


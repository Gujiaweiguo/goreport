## MODIFIED Requirements

### Requirement: JWT Authentication

The Go backend SHALL validate JWT access tokens for all protected API and UI routes, except explicitly configured public routes.

#### Scenario: Valid token access
- **WHEN** a request includes a valid JWT
- **THEN** request is authorized and user identity is available to handlers

#### Scenario: Health check public access
- **WHEN** a request hits `/health` endpoint
- **THEN** request bypasses JWT validation
- **THEN** system returns database status and service status

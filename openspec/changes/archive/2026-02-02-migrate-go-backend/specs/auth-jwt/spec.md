## ADDED Requirements
### Requirement: JWT Authentication
The Go backend SHALL validate JWT access tokens for all protected API and UI routes, except explicitly configured public routes.

#### Scenario: Valid token access
- **WHEN** a request includes a valid JWT
- **THEN** the request is authorized and user identity is available to handlers

### Requirement: Claim Mapping
The system SHALL map JWT claims to user id, role list, and tenant id for authorization decisions.

#### Scenario: Role and tenant resolution
- **WHEN** a JWT includes role and tenant claims
- **THEN** the backend resolves roles and tenant context for the request

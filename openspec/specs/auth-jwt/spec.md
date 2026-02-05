# auth-jwt Specification

## Purpose
TBD - created by archiving change migrate-go-backend. Update Purpose after archive.
## Requirements
### Requirement: JWT Authentication
The Go backend SHALL validate JWT access tokens for all protected API and UI routes, except explicitly configured public routes.

#### Scenario: Valid token access
- **WHEN** a request includes a valid JWT
- **THEN** request is authorized and user identity is available to handlers

#### Scenario: Health check public access
- **WHEN** a request hits `/health` endpoint
- **THEN** request bypasses JWT validation
- **THEN** system returns database status and service status

#### Scenario: Invalid token rejection
- **WHEN** a request includes an invalid or expired JWT
- **THEN** request is rejected with 401 Unauthorized

### Requirement: Claim Mapping
The system SHALL map JWT claims to user id, role list, and tenant id for authorization decisions.

#### Scenario: Role and tenant resolution
- **WHEN** a JWT includes role and tenant claims
- **THEN** the backend resolves roles and tenant context for the request

#### Scenario: User identity extraction
- **WHEN** a request includes a valid JWT
- **THEN** user identity (userId, username, tenantId) is available to handlers via context

### Requirement: Token Generation
The system SHALL generate JWT access tokens for authenticated users containing user claims.

#### Scenario: Generate token on login
- **WHEN** a user successfully logs in with valid credentials
- **THEN** system returns a JWT token with user id, role list, and tenant id claims

#### Scenario: Token payload structure
- **WHEN** a JWT token is generated
- **THEN** token includes userId, username, roles, tenantId, and standard claims (exp, iat, iss)

### Requirement: Password Hashing
The system SHALL securely hash user passwords using bcrypt before storage.

#### Scenario: Hash password on user creation
- **WHEN** a new user is created or password is updated
- **THEN** password is hashed using bcrypt before saving to database

#### Scenario: Verify password on login
- **WHEN** a user attempts to login
- **THEN** provided password is compared to stored hash using bcrypt


## ADDED Requirements
### Requirement: MySQL Schema Compatibility
The Go backend SHALL reuse the existing goReport MySQL schema without requiring destructive migrations.

#### Scenario: Existing database reuse
- **WHEN** the Go service starts against an existing goReport database
- **THEN** existing reports and dashboards remain available

### Requirement: Route Compatibility
The backend SHALL preserve existing route prefixes for UI access (`/jmreport/*` and `/drag/*`).

#### Scenario: Legacy link access
- **WHEN** a user follows a legacy link under `/jmreport/` or `/drag/`
- **THEN** the request resolves successfully without UI changes

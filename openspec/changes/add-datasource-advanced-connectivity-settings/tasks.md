## 1. Datasource Advanced Connectivity

- [x] 1.1 Extend datasource config schema/model with optional SSH tunnel settings.
- [x] 1.2 Implement SSH tunnel validation for password mode.
- [x] 1.3 Implement SSH tunnel validation for key mode.
- [x] 1.4 Integrate SSH tunnel parameters into datasource connection test flow.

## 2. Runtime Controls

- [x] 2.1 Add `maxConnections` setting with bounds validation.
- [x] 2.2 Add `queryTimeoutSeconds` setting with bounds validation.
- [x] 2.3 Apply runtime controls when creating datasource DB connections.

## 3. Connector Profiles

- [x] 3.1 Define datasource type profiles/capability matrix for supported categories.
- [x] 3.2 Validate datasource type-specific required fields against profile.

## 4. Verification

- [x] 4.1 Add backend tests for SSH password/key modes and validation failures.
- [x] 4.2 Add backend tests for timeout/connection bounds behavior.
- [x] 4.3 Add integration tests for datasource connection test with advanced settings.
- [x] 4.4 Run `openspec validate add-datasource-advanced-connectivity-settings --strict --no-interactive`.

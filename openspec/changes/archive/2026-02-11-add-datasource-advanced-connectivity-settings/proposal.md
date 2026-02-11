## Why

Datasource setup currently focuses on basic connection definitions. DataEase datasource guidance highlights production-ready connectivity features, including SSH tunnel access through bastion hosts and runtime connection controls such as connection pool size and query timeout. These are necessary for secure cross-network access and stable query execution.

## What Changes

- Add datasource advanced connectivity settings for SSH tunnel access.
- Support SSH tunnel authentication modes: password and key.
- Add datasource runtime controls: maximum connections and query timeout.
- Define datasource type profile requirements and validation behavior for supported connector categories.

## Capabilities

### Modified Capabilities

- `datasource-management`: add advanced connectivity settings, runtime controls, and connector profile requirements.

## Impact

- Affected specs: `openspec/specs/datasource-management/spec.md`
- Affected code (expected): datasource model/config schema, connection builder, datasource validation and test connection APIs.
- Risk: medium-high because connectivity settings affect security and runtime behavior.

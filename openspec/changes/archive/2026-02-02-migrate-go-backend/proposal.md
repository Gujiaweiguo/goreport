# Change: Migrate goReport backend to Go

## Why
The current integration relies on a Java Spring Boot backend. We need a Go-native backend to align with the existing Go system, reduce cross-stack operational overhead, and enable consistent JWT-based SSO.

## What Changes
- **BREAKING** Replace the Spring Boot example backend runtime with a Go service that provides report, print/export, and dashboard capabilities.
- Add JWT-based authentication and claim mapping compatible with the existing Go system.
- Preserve the current MySQL schema and core data model for reports and dashboards.
- Maintain route compatibility for existing UI paths (e.g., `/jmreport/*`, `/drag/*`).

## Impact
- Affected specs: auth-jwt, report-rendering, report-designer, report-export, datasource-management, bi-dashboard, embedding-integration, migration-compatibility.
- Affected code: new Go service (to be introduced), existing `jimureport-example` used only for reference and data model alignment.

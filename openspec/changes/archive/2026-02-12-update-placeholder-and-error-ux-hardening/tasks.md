## 1. Placeholder UX Remediation
- [x] 1.1 Audit and catalog all user-facing placeholder actions ("开发中", dead buttons, mock-only behavior).
- [x] 1.2 Replace placeholder-only buttons with disabled state + explanatory tooltip/tooltip guidance.
- [x] 1.3 For critical MVP features, implement minimum viable workflow instead of placeholder.
- [x] 1.4 Verify no active production button triggers "开发中" alert only.

## 2. Backend Error Message Hardening
- [x] 2.1 Identify handlers returning generic "invalid request" without field-level context.
- [x] 2.2 Update binding error responses to include actionable validation diagnostics (field path, expected format).
- [x] 2.3 Standardize error response envelope across all modules (success/message optional details).
- [x] 2.4 Ensure no stack traces or internal paths leak to client responses.

## 3. Frontend Error Presentation Alignment
- [x] 3.1 Update API client error interceptors to extract backend `message` field as primary display text.
- [x] 3.2 Replace generic "请求失败" fallbacks with backend diagnostic when available.
- [x] 3.3 Ensure error toast/modal preserves actionable detail without technical noise.
- [x] 3.4 Add frontend tests verifying error message extraction from various response shapes.

## 4. Regression and Validation
- [x] 4.1 Add backend tests for request binding error scenarios with diagnostic assertions.
- [x] 4.2 Add frontend component tests for placeholder state visibility and error rendering.
- [x] 4.3 Run `openspec validate update-placeholder-and-error-ux-hardening --strict --no-interactive`.

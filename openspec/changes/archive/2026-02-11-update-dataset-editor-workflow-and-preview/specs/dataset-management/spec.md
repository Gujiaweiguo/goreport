## ADDED Requirements

### Requirement: Schema-Driven Dataset Preview Rendering
The frontend SHALL render dataset preview using schema-aware mapping rather than fixed field keys.

#### Scenario: Render preview for non-region/sales schema
- **WHEN** preview payload fields differ from fixed demo keys
- **THEN** system derives display mapping from available schema and data
- **THEN** preview still renders meaningful output without hardcoded assumptions

#### Scenario: Fallback when chart mapping is not available
- **WHEN** preview data does not provide a suitable chart dimension/measure pair
- **THEN** system falls back to safe preview presentation
- **THEN** user receives clear feedback about current preview mode

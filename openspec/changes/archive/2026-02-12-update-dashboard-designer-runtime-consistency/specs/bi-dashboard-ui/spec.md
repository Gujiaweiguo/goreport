## ADDED Requirements
### Requirement: Designer Save-Load Runtime Consistency
The dashboard designer SHALL restore persisted component state consistently after save and reload.

#### Scenario: User reloads saved dashboard in designer
- **WHEN** user saves dashboard and later loads it in designer
- **THEN** component list, layer order, style settings, and data bindings are restored from persisted payload
- **THEN** designer does not rely on hardcoded mock component defaults for restoration

### Requirement: Preview Baseline Renders Persisted Components
Dashboard preview SHALL render a minimum runtime baseline based on persisted component payload rather than placeholder-only view.

#### Scenario: Preview renders baseline component output
- **WHEN** user enters preview mode after loading a saved dashboard
- **THEN** preview renders baseline output for supported component types using persisted data
- **THEN** preview refresh action updates the rendered state without switching to placeholder-only content

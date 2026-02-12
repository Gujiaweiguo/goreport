# Rollout Note: Dashboard Designer Runtime Consistency

## Overview
Dashboard persistence and preview runtime have been aligned to support deterministic save/load cycles with actual component data.

## Changes Made

### Backend
1. **Model** (`backend/internal/models/dashboard.go`)
   - Added `DashboardConfig` struct with width, height, backgroundColor
   - Added `DashboardComponent` struct with full component schema
   - Added `Components` and `ComponentsJSON` fields to Dashboard
   - Added GORM hooks for automatic serialization/deserialization
   - Added compatibility defaults for legacy records

2. **Service** (`backend/internal/dashboard/service.go`)
   - Updated `CreateRequest` and `UpdateRequest` to use typed Config and Components
   - Added default config values (1920x1080, backgroundColor: #0a0e27)
   - Proper handling of Config and Components in Create and Update

### Frontend
1. **API** (`frontend/src/api/dashboard.ts`)
   - Types already aligned with new backend contract

2. **Designer** (`frontend/src/views/DashboardDesigner.vue`)
   - Improved error handling to show backend diagnostic messages
   - Save/load flow properly handles components array

3. **Preview** (`frontend/src/components/dashboard/DashboardPreview.vue`)
   - Replaced mock-only rendering with component-type-specific display
   - Text components show actual title
   - Chart/Table/Image components show type label, title, and data source hint
   - Components render with their persisted styles

## API Contract

### Dashboard Structure
```typescript
interface Dashboard {
  id: string
  name: string
  code: string
  config: {
    width: number
    height: number
    backgroundColor: string
  }
  components: DashboardComponent[]
  // ... other fields
}

interface DashboardComponent {
  id: string
  title: string
  type: 'text' | 'chart' | 'table' | 'image' | 'decorative'
  width: number
  height: number
  x: number
  y: number
  visible: boolean
  locked: boolean
  style: object
  data: object
  interaction: object
}
```

### Endpoints
- `POST /api/v1/dashboard/create` - Create with components
- `PUT /api/v1/dashboard/:id` - Update with components
- `GET /api/v1/dashboard/list` - List with components
- `GET /api/v1/dashboard/:id` - Get with components

## Backward Compatibility
- Legacy records with empty config get defaults (1920x1080, dark background)
- Legacy records with empty components get empty array
- Config and Components are serialized to JSON in database

## Testing
- Backend dashboard tests pass (1/1)
- Manual verification: Create → Save → Load → Preview cycle works end-to-end

## Deployment Checklist
- [x] Backend model updated with serialization hooks
- [x] Backend service handles typed Config and Components
- [x] Frontend designer save/load flow aligned
- [x] Preview component renders from persisted data
- [x] Error handling shows backend diagnostics
- [x] OpenSpec validation passed

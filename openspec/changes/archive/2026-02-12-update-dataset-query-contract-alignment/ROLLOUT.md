# Rollout Note: Dataset Query Contract Alignment

## Breaking Change
The dataset query endpoint method has been changed from `GET` to `POST` to properly support request body for complex query parameters.

## Changes Made
1. **Backend** (`backend/internal/httpserver/server.go:144`)
   - Route changed: `GET /api/v1/datasets/:id/data` → `POST /api/v1/datasets/:id/data`

2. **Backend** (`backend/internal/dataset/handler.go:195`)
   - Error message now includes detailed validation diagnostics

3. **Frontend** (`frontend/src/views/ReportDesigner.vue:485`)
   - Error handling now surfaces backend diagnostic messages

## API Contract
### Endpoint
```
POST /api/v1/datasets/{id}/data
```

### Request Body
```typescript
interface QueryRequest {
  fields?: string[]
  filters?: Filter[]
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  page?: number
  pageSize?: number
  aggregations?: Record<string, Aggregation>
}
```

### Response
```typescript
interface QueryResponse {
  success: boolean
  result: {
    data: Record<string, any>[]
    total: number
    page: number
    pageSize: number
    executionTime: number
    aggregations?: Record<string, any>
  }
  message: string
}
```

## Migration Guide
### For API Consumers
If you have external integrations using the old `GET` endpoint:
1. Change HTTP method from `GET` to `POST`
2. Move query parameters from URL to request body
3. Update error handling to expect detailed diagnostic messages

### Compatibility
- Frontend is already updated to use `POST`
- Report designer preview flow is now working correctly
- Error messages are now actionable with field-level details

## Testing
- All existing backend tests pass
- Frontend dataset API tests pass (4/4)
- Manual verification: Report designer data preview now works end-to-end

## Deployment Checklist
- [x] Backend route updated
- [x] Handler error messages improved
- [x] Frontend error handling updated
- [x] All tests passing
- [x] OpenSpec validation passed

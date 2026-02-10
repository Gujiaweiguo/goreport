<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# AGENTS.md
# Guidance for agentic coding in this repository.

## Scope and layout
- Primary runnable project lives in `backend/` (Go) and `frontend/` (Vue).
- Main Go entrypoint: `backend/cmd/server/main.go`.
- Docs reference running from project root with Docker Compose.

## Build, lint, test, run

### Go Backend
- Build:
  - `cd backend && go build -o bin/server cmd/server/main.go`
  - Or use `make build`
- Run:
  - `cd backend && go run cmd/server/main.go`
  - Or use `make dev` to start all services
- Test:
  - `cd backend && go test ./... -v`
  - Or use `make test`

### Frontend
- Install dependencies:
  - `cd frontend && npm install`
- Run dev server:
  - `cd frontend && npm run dev`
- Build:
  - `cd frontend && npm run build`

### Docker
- Start all services:
  - `make dev`
- View logs:
  - `make logs`
- Stop services:
  - `make dev-down`

### Database bootstrap
- Initialize DB schema:
  - `backend/db/init.sql` is automatically loaded by Docker Compose.
  - Or manually: `mysql -uroot -p < backend/db/init.sql`

### Lint/format
- Go: Uses `gofmt` (standard)
  - `gofmt -w backend/`
- No explicit frontend linter configured.

## Code style guidelines

### Language and frameworks
- Go 1.22+
- Gin web framework
- GORM for database
- JWT for authentication
- Vue 3 + TypeScript
- Vite build tool
- Element Plus UI library

### Project structure (Go)
- `backend/cmd/server/` - Application entry point
- `backend/internal/` - Internal packages
  - `auth/` - JWT authentication
  - `config/` - Configuration management
  - `models/` - Data models
  - `repository/` - Data access layer
  - `service/` - Business logic layer
  - `httpserver/` - HTTP handlers
  - `middleware/` - HTTP middleware

### Go code style
- Standard Go formatting (`gofmt`)
- 4-space indentation (project convention)
- Package comments for public packages
- Exported functions/structs have comments
- Error handling: check errors explicitly, never ignore
- Use `gorm` tags for model fields
- Use `json` tags for API structs

### Naming conventions
- Go: Follow standard Go conventions
  - `PascalCase` for exported names
  - `camelCase` for unexported names
  - `UPPER_SNAKE_CASE` for constants
- Frontend: Vue/TypeScript conventions
  - Components: `PascalCase.vue`
  - Functions: `camelCase`
  - Constants: `UPPER_SNAKE_CASE`

### Error handling (Go)
- Always check errors
- Return errors to caller
- Log errors at appropriate level
- Use custom error types when needed

### HTTP conventions
- RESTful API design
- Versioned routes: `/api/v1/...`
- JSON request/response bodies
- Standard HTTP status codes
- JWT in `Authorization` header

### Auth conventions
- JWT for stateless authentication
- Token format: `Bearer <token>`
- Claims include: user_id, tenant_id, exp
- Middleware extracts user info to context

### Comments and docs
- Go: Package and function comments
- English preferred for technical terms
- Chinese OK for business logic comments

## Testing guidance

### Go testing
- Standard `go test`
- Table-driven tests preferred
- Mock external dependencies
- Test files: `*_test.go`

### Running tests
```bash
# All tests
cd backend && go test ./... -v

# Specific package
cd backend && go test ./internal/dashboard -v

# Coverage
cd backend && go test ./... -cover
```

## Environment assumptions
- Go 1.22+ required
- Node.js 20+ required
- MySQL 5.7+ required
- Redis optional
- Docker recommended for local development

## Docker notes
- Docker Compose for local development
- Production Dockerfile in `backend/Dockerfile` and `frontend/Dockerfile`

## Localization
- Comments and log messages can be bilingual
- Preserve existing language style

## Repo-specific notes
- `.gitignore` ignores `target/`, `logs/`, `node_modules/`, and IDE files
- No CI configuration currently

## Cursor/Copilot rules
- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` found

## 规则
- 默认中文回复
- 有歧义先提问；先给计划/方案，得到确认再执行
- 最小改动，不做无关重构
- 不新增依赖，需先说明并确认
- 未明确要求不提交；不使用 --amend

## OpenSpec（如使用）
- 新功能/改接口/改表结构：先 Proposal，经批准后 Apply
- 上线后 Archive，并执行 `openspec validate --strict --no-interactive`

## Sources consulted
- `README.md`
- `README.en-US.md`
- `backend/go.mod`
- `backend/cmd/server/main.go`
- `frontend/package.json`
- `docker-compose.yml`
- `.gitignore`

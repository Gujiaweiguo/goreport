# goReport Makefile
# 提供常用开发命令

.PHONY: help dev dev-down dev-logs dev-ps build build-frontend build-backend build-prod test test-frontend test-backend test-coverage clean db-shell redis-cli docs

# 默认目标
help:
	@echo "goReport 开发命令:"
	@echo ""
	@echo "  make dev         - 启动开发环境 (Docker Compose)"
	@echo "  make dev-down    - 停止开发环境"
	@echo "  make build       - 构建生产镜像"
	@echo "  make build-prod  - 启动生产环境"
	@echo "  make test        - 运行测试"
	@echo "  make test-frontend - 运行前端测试"
	@echo "  make test-backend - 运行后端测试"
	@echo "  make test-coverage - 生成测试覆盖率报告"
	@echo "  make clean       - 清理容器和卷"
	@echo "  make logs        - 查看日志"
	@echo "  make ps          - 查看容器状态"
	@echo "  make docs        - 查看文档列表"
	@echo "  make docs-user   - 查看用户指南"
	@echo "  make docs-dev    - 查看开发指南"
	@echo ""

# 开发环境
dev:
	docker-compose up -d
	@echo "开发环境已启动:"
	@echo "  前端: http://localhost:3000"
	@echo "  后端: http://localhost:8085"
	@echo "  MySQL: localhost:3306"
	@echo "  Redis: localhost:6379"

dev-down:
	docker-compose down

dev-logs:
	docker-compose logs -f

dev-ps:
	docker-compose ps

# 构建
build:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# 测试
test:
	cd backend && go test ./... -v

test-frontend:
	cd frontend && npm test -- --passWithNoTests

test-backend:
	cd backend && go test ./... -v

test-coverage:
	cd backend && go test -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html

# 构建
build-frontend:
	cd frontend && npm ci --only=production
	cd frontend && npm run build

build-backend:
	cd backend && go mod download
	cd backend && go build -ldflags="-s -w" -o bin/server cmd/server/main.go

build: build-frontend build-backend

# 生产部署
build-prod:
	docker-compose -f docker-compose.prod.yml build
	docker-compose -f docker-compose.prod.yml up -d
	@echo "生产环境已启动:"
	@echo "  前端: http://localhost"
	@echo "  后端: http://localhost:8085"

# 清理
clean:
	docker-compose down -v
	docker system prune -f

# 数据库
db-shell:
	docker-compose exec mysql mysql -uroot -proot goreport

redis-cli:
	docker-compose exec redis redis-cli

# 文档
docs:
	@echo "可用文档："
	@echo "  - 用户指南: docs/USER_GUIDE.md"
	@echo "  - 开发指南: docs/DEVELOPMENT_GUIDE.md"
	@echo "  - 迁移指南: docs/MIGRATION_GUIDE.md"
	@echo "  - 贡献指南: docs/CONTRIBUTING.md"

# 查看文档
docs-user:
	cat docs/USER_GUIDE.md | less

docs-dev:
	cat docs/DEVELOPMENT_GUIDE.md | less

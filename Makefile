# JimuReport Makefile
# 提供常用开发命令

.PHONY: help dev build test clean

# 默认目标
help:
	@echo "JimuReport 开发命令:"
	@echo ""
	@echo "  make dev         - 启动开发环境 (Docker Compose)"
	@echo "  make dev-down    - 停止开发环境"
	@echo "  make build       - 构建生产镜像"
	@echo "  make test        - 运行测试"
	@echo "  make clean       - 清理容器和卷"
	@echo "  make logs        - 查看日志"
	@echo "  make ps          - 查看容器状态"
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
	cd backend && go test ./...
	test:
	cd frontend && npm test

# 清理
clean:
	docker-compose down -v
	docker system prune -f

# 数据库
db-shell:
	docker-compose exec mysql mysql -uroot -proot jimureport

redis-cli:
	docker-compose exec redis redis-cli

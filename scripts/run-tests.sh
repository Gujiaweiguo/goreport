#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

COMPOSE_FILE="docker-compose.test.yml"

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

cleanup() {
    print_info "清理测试环境..."
    docker compose -f "$COMPOSE_FILE" down -v
    docker volume rm goreport_mysql-test-data goreport_redis-test-data 2>/dev/null || true
    print_info "测试环境已清理"
}

trap cleanup EXIT INT TERM

print_info "启动测试环境..."
docker compose -f "$COMPOSE_FILE" up -d --build

print_info "等待 MySQL 服务就绪..."
MAX_RETRIES=30
RETRY_COUNT=0
until docker compose -f "$COMPOSE_FILE" exec -T mysql mysqladmin ping -h localhost -u root -proot &>/dev/null; do
    RETRY_COUNT=$((RETRY_COUNT + 1))
    if [ $RETRY_COUNT -ge $MAX_RETRIES ]; then
        print_error "MySQL 服务启动超时"
        exit 1
    fi
    echo -n "."
    sleep 2
done
echo ""

print_info "运行后端测试..."
BACKEND_RESULT=0
docker compose -f "$COMPOSE_FILE" exec -T backend go test ./... -v || BACKEND_RESULT=$?

print_info "运行前端测试..."
FRONTEND_RESULT=0
docker compose -f "$COMPOSE_FILE" exec -T frontend npm run test:run -- --passWithNoTests || FRONTEND_RESULT=$?

print_info ""
print_info "======================================"
print_info "测试结果汇总"
print_info "======================================"

if [ $BACKEND_RESULT -eq 0 ]; then
    echo -e "${GREEN}后端测试: PASS${NC}"
else
    echo -e "${RED}后端测试: FAIL${NC}"
fi

if [ $FRONTEND_RESULT -eq 0 ]; then
    echo -e "${GREEN}前端测试: PASS${NC}"
else
    echo -e "${RED}前端测试: FAIL${NC}"
fi

print_info ""

if [ $BACKEND_RESULT -ne 0 ] || [ $FRONTEND_RESULT -ne 0 ]; then
    print_error "部分测试失败，请查看详细日志"
    exit 1
else
    print_info "所有测试通过"
    exit 0
fi

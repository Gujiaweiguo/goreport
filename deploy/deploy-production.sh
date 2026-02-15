#!/bin/bash

# goReport 生产部署脚本
# 用于自动化生产环境部署

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BACKUP_DIR="./backups/$(date +%Y%m%d_%H%M%S)"
COMPOSE_FILE="docker-compose.prod.yml"
ENV_FILE=".env.prod"

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查环境
print_info "检查部署环境..."

if [ ! -f "$COMPOSE_FILE" ]; then
    print_error "未找到 $COMPOSE_FILE 文件"
    exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
    print_error "未找到 $ENV_FILE 文件"
    print_info "请复制 .env.example 并配置生产环境变量"
    exit 1
fi

print_info "环境检查通过"
echo ""

# 读取版本号
VERSION=$(git describe --tags --always 2>/dev/null || echo "v1.0.0")
print_info "部署版本: $VERSION"
echo ""

# 确认部署
print_warn "即将执行生产环境部署!"
print_warn "版本: $VERSION"
print_warn "环境文件: $ENV_FILE"
echo ""
read -p "确认继续? (yes/no): " confirm
echo

if [ "$confirm" != "yes" ]; then
    print_info "部署已取消"
    exit 0
fi

# 创建备份目录
print_info "创建备份目录: $BACKUP_DIR"
mkdir -p "$BACKUP_DIR"

# 备份数据库
print_info "备份数据库..."
if docker compose -f $COMPOSE_FILE ps | grep -q "mysql"; then
    docker compose -f $COMPOSE_FILE exec -T mysql mysqldump -uroot -p"${DB_PASSWORD}" goreport > "$BACKUP_DIR/database_backup.sql"
    print_info "数据库备份完成: $BACKUP_DIR/database_backup.sql"
else
    print_warn "MySQL 服务未运行，跳过数据库备份"
fi
echo ""

# 备份当前运行的容器（如果有）
print_info "备份当前部署状态..."
docker compose -f $COMPOSE_FILE ps > "$BACKUP_DIR/container_status.txt" 2>/dev/null || true
docker images > "$BACKUP_DIR/docker_images.txt" 2>/dev/null || true
print_info "状态备份完成"
echo ""

# 拉取最新代码
print_info "拉取最新代码..."
git fetch origin
git checkout "tags/$VERSION" 2>/dev/null || git checkout main
print_info "代码更新完成"
echo ""

# 构建镜像
print_info "构建 Docker 镜像..."
docker compose -f $COMPOSE_FILE --env-file $ENV_FILE build --no-cache
print_info "镜像构建完成"
echo ""

# 停止旧服务
print_info "停止旧服务..."
docker compose -f $COMPOSE_FILE down
print_info "旧服务已停止"
echo ""

# 清理旧镜像（可选）
read -p "是否清理旧镜像? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    print_info "清理旧镜像..."
    docker image prune -f
    print_info "旧镜像已清理"
fi
echo ""

# 启动新服务
print_info "启动新服务..."
docker compose -f $COMPOSE_FILE --env-file $ENV_FILE up -d
print_info "服务已启动"
echo ""

# 等待服务启动
print_info "等待服务启动..."
sleep 10

# 健康检查
print_info "执行健康检查..."
HEALTH_CHECK_URL="http://localhost:8085/health"
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -s "$HEALTH_CHECK_URL" > /dev/null 2>&1; then
        print_info "健康检查通过!"
        break
    fi
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    print_warn "健康检查失败，重试 $RETRY_COUNT/$MAX_RETRIES..."
    sleep 5
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    print_error "健康检查失败，请检查服务日志"
    docker compose -f $COMPOSE_FILE logs
    exit 1
fi

# 执行数据库迁移（如果有）
print_info "检查数据库迁移..."
if [ -d "db/migrations" ]; then
    for migration in db/migrations/*.sql; do
        if [ -f "$migration" ]; then
            print_info "执行迁移: $migration"
            docker compose -f $COMPOSE_FILE exec -T mysql mysql -uroot -p"${DB_PASSWORD}" goreport < "$migration"
        fi
    done
fi

# 验证部署
print_info "验证部署状态..."
docker compose -f $COMPOSE_FILE ps

echo ""
print_info "================================"
print_info "部署完成!"
print_info "================================"
echo ""
print_info "访问地址:"
print_info "  前端: http://localhost"
print_info "  后端: http://localhost:8085"
print_info "  健康检查: http://localhost:8085/health"
echo ""
print_info "备份位置: $BACKUP_DIR"
echo ""
print_info "查看日志:"
print_info "  docker compose -f $COMPOSE_FILE logs -f"
echo ""

# 生成部署报告
cat > "$BACKUP_DIR/deploy_report.txt" << EOF
部署报告
========

部署时间: $(date '+%Y-%m-%d %H:%M:%S')
部署版本: $VERSION
部署环境: 生产环境
部署状态: 成功

服务状态:
$(docker compose -f $COMPOSE_FILE ps)

镜像信息:
$(docker images | grep goreport)

访问地址:
- 前端: http://localhost
- 后端: http://localhost:8085

备份位置: $BACKUP_DIR

回滚命令:
  # 如果需要回滚，执行:
  docker compose -f $COMPOSE_FILE down
  # 恢复数据库:
  # docker compose -f $COMPOSE_FILE exec -T mysql mysql -uroot -p"${DB_PASSWORD}" goreport < $BACKUP_DIR/database_backup.sql
  # 重新部署上一个版本
EOF

print_info "部署报告已生成: $BACKUP_DIR/deploy_report.txt"

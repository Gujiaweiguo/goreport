# 数据集功能部署指南

## 概述

本文档描述数据集（Dataset）功能的部署流程和监控配置。

## 前置条件

### 1. 基础环境

- Go 1.22+
- Node.js 20+
- MySQL 5.7+ 或 8.0+
- Redis 6.0+（可选，用于缓存）
- Docker 和 Docker Compose（推荐）

### 2. 现有系统

- goReport 系统已部署并运行
- 数据源（DataSource）功能已配置
- 用户认证系统已启用

## 部署步骤

### 1. 数据库迁移

#### 1.1 备份数据库

```bash
# 使用 mysqldump 备份现有数据库
mysqldump -u root -p goreport > backup_$(date +%Y%m%d_%H%M%S).sql

# 或使用 Docker 容器
docker exec goreport-mysql mysqldump -u root -proot goreport > backup.sql
```

#### 1.2 执行迁移脚本

```bash
# 进入数据库迁移目录
cd backend/db/migrations

# 执行迁移脚本
mysql -u root -p goreport < 001_add_datasets.sql

# 或使用 Docker 容器
docker exec -i goreport-mysql mysql -u root -proot goreport < 001_add_datasets.sql
```

#### 1.3 验证迁移

```bash
# 连接到数据库
mysql -u root -p goreport

# 验证表已创建
SHOW TABLES LIKE 'dataset%';

# 验证表结构
DESC datasets;
DESC dataset_fields;
DESC dataset_sources;

# 退出
exit
```

### 2. 后端部署

#### 2.1 更新 Go 依赖

```bash
cd backend
go mod download
```

#### 2.2 编译后端

```bash
# 构建二进制文件
go build -o bin/server cmd/server/main.go

# 或使用 Docker
docker build -t goreport-backend:latest -f Dockerfile .
```

#### 2.3 配置环境变量

创建或更新 `.env` 文件：

```bash
# 数据集相关配置
DATASET_CACHE_ENABLED=true
DATASET_CACHE_TTL=3600
DATASET_EXPRESSION_CACHE_ENABLED=true
DATASET_MAX_EXPRESSIONS=1000

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=root
DB_NAME=goreport

# Redis 配置（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
REDIS_DB=0

# JWT 配置
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
```

#### 2.4 启动后端服务

```bash
# 直接运行
./bin/server

# 或使用 Docker Compose
docker-compose up -d backend

# 或使用 systemd（生产环境）
sudo systemctl start goreport-backend
```

#### 2.5 验证后端

```bash
# 检查健康检查端点
curl http://localhost:8085/health

# 检查数据集 API
curl http://localhost:8085/api/v1/datasets \
  -H "Authorization: Bearer <your-token>"

# 查看后端日志
tail -f logs/backend.log
# 或 Docker logs
docker logs goreport-backend -f
```

### 3. 前端部署

#### 3.1 构建前端

```bash
cd frontend

# 安装依赖
npm install

# 构建生产版本
npm run build

# 或使用 Docker
docker build -t goreport-frontend:latest -f Dockerfile .
```

#### 3.2 部署静态文件

```bash
# 使用 nginx 部署
sudo cp -r dist/* /var/www/goreport/

# 或使用 Docker Compose
docker-compose up -d frontend
```

#### 3.3 验证前端

```bash
# 检查前端服务
curl http://localhost:3000

# 在浏览器中访问
# http://localhost:3000
# 导航到数据集页面
```

### 4. 数据库索引优化

为提升查询性能，创建以下索引：

```sql
-- 数据集表索引
CREATE INDEX idx_datasets_tenant_type_status ON datasets(tenant_id, type, status);
CREATE INDEX idx_datasets_created_at ON datasets(created_at DESC);

-- 数据集字段表索引
CREATE INDEX idx_dataset_fields_dataset_type ON dataset_fields(dataset_id, type);
CREATE INDEX idx_dataset_fields_computed ON dataset_fields(is_computed);

-- 数据集源表索引
CREATE INDEX idx_dataset_sources_dataset_type ON dataset_sources(dataset_id, source_type);
```

### 5. 缓存配置

#### 5.1 Redis 配置（推荐）

```bash
# Redis 缓存配置示例
REDIS_MAXMEMORY=256mb
REDIS_MAXMEMORY_POLICY=allkeys-lru
REDIS_SAVE_INTERVAL=300
REDIS_TIMEOUT=300
```

#### 5.2 验证缓存

```bash
# 检查 Redis 连接
redis-cli ping

# 查看缓存键
redis-cli keys "dataset:*"

# 清空缓存（如需要）
redis-cli flushdb
```

## 监控配置

### 1. 应用监控

#### 1.1 Prometheus 监控

创建 `prometheus.yml` 配置：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'goreport'
    static_configs:
      - targets: ['localhost:8085']
    metrics_path: '/metrics'
```

启动 Prometheus：

```bash
docker run -d \
  --name prometheus \
  -p 9090:9090 \
  -v /path/to/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus
```

#### 1.2 Grafana 仪表盘

创建 Grafana 数据源和仪表盘：

```bash
# 启动 Grafana
docker run -d \
  --name grafana \
  -p 3001:3000 \
  grafana/grafana

# 访问 http://localhost:3001
# 默认用户名/密码：admin/admin
```

推荐监控指标：

| 指标 | 描述 | 警告阈值 |
|------|------|----------|
| dataset_query_duration | 数据集查询延迟 | > 500ms |
| dataset_cache_hit_rate | 缓存命中率 | < 80% |
| dataset_error_rate | 错误率 | > 5% |
| dataset_active_count | 活跃数据集数量 | - |

### 2. 数据库监控

#### 2.1 慢查询日志

```sql
-- 启用慢查询日志
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- 分析慢查询
SELECT * FROM mysql.slow_log
WHERE sql_text LIKE '%dataset%'
ORDER BY query_time DESC
LIMIT 10;
```

#### 2.2 表大小监控

```sql
-- 查看数据集相关表大小
SELECT
    table_name AS 'Table',
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.TABLES
WHERE table_schema = 'goreport'
  AND table_name LIKE 'dataset%'
ORDER BY (data_length + index_length) DESC;
```

### 3. 日志监控

#### 3.1 配置日志级别

在 `backend/config/config.yaml` 中：

```yaml
logging:
  level: info  # debug, info, warn, error
  format: json
  output:
    - stdout
    - file: logs/dataset.log
```

#### 3.2 集中式日志（可选）

使用 ELK Stack 或 Loki：

```yaml
# 示例：Loki 配置
loki:
  url: http://localhost:3100/loki/api/v1/push
  batch_size: 100
  batch_wait: 1s
```

### 4. 健康检查

实现以下健康检查端点：

```bash
# 数据集服务健康检查
curl http://localhost:8085/health/datasets

# 返回示例
{
  "status": "healthy",
  "checks": {
    "database": "ok",
    "cache": "ok",
    "expression_cache": "ok"
  },
  "timestamp": "2024-02-08T10:00:00Z"
}
```

## 性能优化

### 1. 数据库优化

```sql
-- 分区大数据表（如果数据量超过 1000 万）
ALTER TABLE dataset_fields PARTITION BY RANGE (YEAR(created_at));

-- 优化表
OPTIMIZE TABLE datasets;
OPTIMIZE TABLE dataset_fields;
OPTIMIZE TABLE dataset_sources;

-- 分析表
ANALYZE TABLE datasets;
ANALYZE TABLE dataset_fields;
ANALYZE TABLE dataset_sources;
```

### 2. 缓存策略

- **查询缓存**：缓存常用查询结果（TTL: 1小时）
- **表达式缓存**：缓存计算字段表达式（TTL: 24小时）
- **Schema 缓存**：缓存数据集 Schema（TTL: 10分钟）

### 3. 连接池配置

```yaml
# backend/config/config.yaml
database:
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 3600
  conn_max_idle_time: 300
```

## 回滚计划

### 1. 回滚步骤

```bash
# 1. 停止服务
systemctl stop goreport-backend

# 2. 恢复数据库备份
mysql -u root -p goreport < backup_YYYYMMDD_HHMMSS.sql

# 3. 回滚代码版本
git checkout <previous-version>

# 4. 重新构建和部署
docker-compose down
docker-compose up -d

# 5. 验证服务
curl http://localhost:8085/health
```

### 2. 回滚检查清单

- [ ] 数据库迁移已回滚
- [ ] 应用代码已回滚
- [ ] 服务已重启
- [ ] 健康检查通过
- [ ] 关键功能验证通过
- [ ] 日志无错误

## 安全检查

### 1. 权限检查

```sql
-- 验证用户权限
SHOW GRANTS FOR 'goreport_user'@'%';

-- 确保只有必要的权限
GRANT SELECT, INSERT, UPDATE, DELETE ON goreport.* TO 'goreport_user'@'%';
```

### 2. 数据加密

- 使用 TLS 加密数据库连接
- 使用 HTTPS 加密 API 通信
- 敏感配置使用环境变量或加密存储

### 3. SQL 注入防护

系统已内置 SQL 注入防护：
- 参数化查询
- 输入验证
- 使用 GORM ORM

## 故障排查

### 1. 常见问题

| 问题 | 可能原因 | 解决方案 |
|------|----------|----------|
| 数据集创建失败 | 数据源连接错误 | 检查数据源配置 |
| 查询超时 | 数据库性能问题 | 添加索引或优化查询 |
| 缓存不生效 | Redis 未启动 | 检查 Redis 服务状态 |
| API 返回 500 | 服务内部错误 | 查看应用日志 |

### 2. 日志分析

```bash
# 查看错误日志
grep -i error logs/dataset.log

# 查看慢查询日志
grep -i "slow" logs/backend.log

# 统计错误类型
awk '/ERROR/ {print $NF}' logs/dataset.log | sort | uniq -c
```

## 部署后验证

### 1. 功能验证

- [ ] 数据集列表加载正常
- [ ] 创建数据集成功
- [ ] 编辑数据集成功
- [ ] 删除数据集成功
- [ ] 字段管理正常
- [ ] 计算字段正常工作
- [ ] 数据查询返回正确
- [ ] 缓存功能正常

### 2. 性能验证

- [ ] 数据集列表加载时间 < 1s
- [ ] 数据集查询时间 < 500ms
- [ ] 缓存命中率 > 80%
- [ ] API 响应时间 < 200ms

### 3. 集成验证

- [ ] 报表设计器可以使用数据集
- [ ] 仪表盘可以使用数据集
- [ ] 图表编辑器可以使用数据集
- [ ] 权限控制正常

## 维护计划

### 1. 日常维护

- 每日检查错误日志
- 每周检查慢查询
- 每周检查缓存命中率
- 每月分析表大小和碎片

### 2. 定期维护

- 每月优化数据库表
- 每月清理过期缓存
- 每季度审查索引使用情况
- 每半年清理历史数据

## 支持

如遇到问题，请联系：

- 技术支持：[weiguogu@163.com](mailto:weiguogu@163.com)
- GitHub Issues：[https://github.com/gujiaweiguo/goreport/issues](https://github.com/gujiaweiguo/goreport/issues)
- 文档：[项目 Wiki](https://github.com/gujiaweiguo/goreport/wiki)

## 附录

### A. 环境变量参考

| 变量名 | 默认值 | 描述 |
|---------|---------|------|
| DB_HOST | localhost | 数据库主机 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USER | root | 数据库用户 |
| DB_PASS | - | 数据库密码 |
| DB_NAME | goreport | 数据库名称 |
| REDIS_HOST | localhost | Redis 主机 |
| REDIS_PORT | 6379 | Redis 端口 |
| REDIS_PASS | - | Redis 密码 |
| JWT_SECRET | - | JWT 密钥 |
| JWT_EXPIRATION | 24h | Token 过期时间 |
| PORT | 8085 | 应用端口 |

### B. 端口使用

| 服务 | 端口 | 协议 | 描述 |
|------|------|------|------|
| 后端 API | 8085 | HTTP/HTTPS | REST API 服务 |
| 前端 Web | 3000 | HTTP | Web UI 服务 |
| MySQL | 3306 | TCP | 数据库服务 |
| Redis | 6379 | TCP | 缓存服务 |
| Prometheus | 9090 | HTTP | 监控服务 |
| Grafana | 3001 | HTTP | 可视化服务 |

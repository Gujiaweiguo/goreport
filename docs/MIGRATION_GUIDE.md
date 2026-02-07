# JimuReport 迁移指南

本指南帮助您从旧版本或其他报表系统迁移到 JimuReport。

## 目录

- [迁移前准备](#迁移前准备)
- [数据迁移](#数据迁移)
- [配置迁移](#配置迁移)
- [代码迁移](#代码迁移)
- [验证测试](#验证测试)
- [常见问题](#常见问题)

## 迁移前准备

### 环境检查

确保满足以下要求：

- **数据库**：MySQL 5.7+ 或 MySQL 8.0+
- **缓存**：Redis 5.0+ (可选）
- **运行时**：Docker 20.10+
- **网络**：确保网络连接稳定

### 数据备份

在迁移前务必备份数据：

```bash
# 备份数据库
mysqldump -uroot -p jimureport > backup_$(date +%Y%m%d).sql

# 压缩备份
gzip backup_$(date +%Y%m%d).sql
```

### 资源评估

评估需要迁移的资源：

- **报表数量**：统计报表总数
- **数据源数量**：统计数据库连接数
- **用户数量**：统计用户和租户数
- **存储需求**：评估图片、文件等存储大小

## 数据迁移

### 数据库迁移

#### 1. 导出旧数据

从旧系统导出数据：

```sql
-- 导出报表数据
SELECT * FROM reports INTO OUTFILE '/tmp/reports.csv';

-- 导出数据源配置
SELECT * FROM data_sources INTO OUTFILE '/tmp/data_sources.csv';

-- 导出用户数据
SELECT * FROM users INTO OUTFILE '/tmp/users.csv';
```

#### 2. 转换数据格式

JimuReport 使用 JSON 格式存储配置，需要转换：

| 旧系统字段 | JimuReport 字段 | 转换规则 |
|-----------|-----------------|---------|
| report_xml | config | XML → JSON |
| template | config | 模板语法适配 |
| sql_query | config.datasource.query | 提取 SQL 语句 |

#### 3. 导入数据

使用脚本导入转换后的数据：

```bash
# 导入报表
python scripts/import_reports.py --source old_export.json

# 导入数据源
python scripts/import_datasources.py --source old_export.json

# 导入用户
python scripts/import_users.py --source old_export.json
```

### SQL 迁移

如果使用 SQL 模板报表：

```bash
# 使用转换工具
./tools/sql-converter --input reports/ --output converted/

# 验证转换
./tools/sql-converter --validate converted/
```

## 配置迁移

### 系统配置

迁移系统配置到环境变量：

```bash
# 创建 .env 文件
cat > .env << EOF
DB_DSN=mysql://user:pass@localhost:3306/jimureport
JWT_SECRET=$(openssl rand -hex 32)
CACHE_ENABLED=true
CACHE_ADDR=localhost:6379
EOF
```

### 数据源配置

将旧数据源配置映射到 JimuReport 格式：

**旧格式示例：**

```json
{
  "driver": "mysql",
  "host": "192.168.1.100",
  "port": 3306,
  "database": "reportdb",
  "username": "report_user"
}
```

**JimuReport 格式：**

```json
{
  "type": "mysql",
  "host": "192.168.1.100",
  "port": 3306,
  "database_name": "reportdb",
  "username": "report_user",
  "password": "encrypted_password"
}
```

### 用户权限

映射旧系统的角色到 JimuReport 角色：

| 旧角色 | JimuReport 角色 | 权限 |
|-------|--------------|--------|
| Administrator | admin | 全部权限 |
| Developer | developer | 设计、编辑 |
| Viewer | viewer | 仅查看 |

## 代码迁移

### API 调用更新

旧 API 端点 → 新 API 端点：

| 旧端点 | 新端点 | 方法 |
|---------|---------|------|
| /report/list | /api/v1/jmreport/list | GET |
| /report/create | /api/v1/jmreport/create | POST |
| /report/update | /api/v1/jmreport/update | POST |
| /report/delete | /api/v1/jmreport/delete | DELETE |

**示例迁移：**

```javascript
// 旧代码
fetch('/report/list?id=' + reportId)

// 新代码
import { reportApi } from '@/api/report'
reportApi.get(reportId)
```

### 前端组件迁移

JimuReport 使用 Vue 3 + Element Plus，需要适配：

```javascript
// 旧组件 (Vue 2)
<template>
  <el-button @click="handleClick">Click</el-button>
</template>

// 新组件 (Vue 3)
<script setup lang="ts">
const handleClick = () => {
  console.log('Clicked')
}
</script>

<template>
  <el-button @click="handleClick">Click</el-button>
</template>
```

### 表达式语法

JimuReport 表达式语法：

| 类型 | 语法 | 示例 |
|-----|------|------|
| 字段引用 | `${field}` | `${username}` |
| 函数调用 | `@func()` | `@sum(amount)` |
| 条件判断 | `? :` | `${amount} > 100 ? 'High' : 'Low'` |

**迁移示例：**

```javascript
// 旧表达式
{user_name}

// 新表达式
${user_name}

// 旧函数
#DATE_FORMAT(date, '%Y-%m-%d')

// 新函数
@DATE_FORMAT(date, 'yyyy-MM-dd')
```

## 验证测试

### 功能验证

创建验证清单：

- [ ] 用户登录成功
- [ ] 报表列表正确显示
- [ ] 报表预览数据正确
- [ ] 数据源连接正常
- [ ] 导出功能正常
- [ ] 权限控制生效

### 性能验证

对比迁移前后的性能：

| 指标 | 迁移前 | 迁移后 | 目标 |
|-------|--------|--------|------|
| 报表加载时间 | 3s | < 1s | 改善 |
| 数据查询时间 | 2s | < 0.5s | 改善 |
| 并发用户数 | 50 | 100+ | 提升 |

### 数据完整性

验证数据迁移完整性：

```sql
-- 检查报表数量
SELECT COUNT(*) FROM reports;

-- 检查数据源数量
SELECT COUNT(*) FROM data_sources;

-- 验证关键报表数据
SELECT id, name, config FROM reports WHERE id IN (
  'report-1', 'report-2', 'report-3'
);
```

## 回滚方案

如果迁移失败，执行回滚：

```bash
# 停止新系统
make dev-down

# 恢复旧数据库
mysql -uroot -p jimureport < backup_20240206.sql

# 重启旧系统
systemctl start old-report-system
```

## 常见问题

### 编码问题

**问题**：导入数据后出现乱码

**解决方案：**

```bash
# 确保数据库字符集
ALTER DATABASE jimureport CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 导出时指定编码
mysqldump --default-character-set=utf8mb4 -uroot -p jimureport > backup.sql
```

### 数据量过大

**问题**：数据导入超时

**解决方案：**

1. 分批导入数据
2. 调整数据库配置：

```ini
max_allowed_packet = 256M
innodb_buffer_pool_size = 1G
```

3. 使用导入脚本

```bash
./tools/batch-import --chunk-size 1000
```

### 配置不兼容

**问题**：部分配置无法迁移

**解决方案：**

1. 使用 OpenSpec 提交配置转换需求
2. 编写自定义转换脚本
3. 手动配置不兼容部分

### 权限丢失

**问题**：迁移后用户权限不对

**解决方案：**

```sql
-- 重新映射角色
UPDATE user_tenants
SET role = 'admin'
WHERE user_id IN ('user-1', 'user-2');

-- 验证权限
SELECT u.username, ut.role
FROM users u
JOIN user_tenants ut ON u.id = ut.user_id;
```

## 迁移工具

### 数据导出工具

```bash
# 从旧系统导出
./tools/export --source mysql --output export.json
```

### 数据转换工具

```bash
# 转换数据格式
./tools/convert --input export.json --output jimureport.json
```

### 数据导入工具

```bash
# 导入到 JimuReport
./tools/import --input jimureport.json --validate
```

## 获取帮助

如遇到迁移问题：

- 查看迁移文档：`https://docs.jimureport.com/migration`
- 提交问题：GitHub Issues (标签：migration)
- 联系支持：`migration@jimureport.com`

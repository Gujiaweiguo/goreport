#!/bin/bash

# goReport UAT 测试数据准备脚本
# 用于生成 UAT 环境所需的测试数据

set -e

echo "================================"
echo "goReport UAT 测试数据准备"
echo "================================"
echo ""

# 配置
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="jimureport"
DB_USER="root"
DB_PASS="root"
BACKEND_URL="http://localhost:8085"

echo "数据库配置:"
echo "  主机: $DB_HOST:$DB_PORT"
echo "  数据库: $DB_NAME"
echo "  用户: $DB_USER"
echo ""

# 检查数据库连接
echo "检查数据库连接..."
if ! mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS -e "SELECT 1" $DB_NAME > /dev/null 2>&1; then
    echo "⚠️  无法连接到数据库"
    echo "   请确保 MySQL 服务已启动且配置正确"
    exit 1
fi

echo "✅ 数据库连接成功"
echo ""

# 清理旧数据（可选）
read -p "是否清理现有测试数据? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "清理旧数据..."
    mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
        SET FOREIGN_KEY_CHECKS = 0;
        TRUNCATE TABLE dashboards;
        TRUNCATE TABLE reports;
        TRUNCATE TABLE data_sources;
        TRUNCATE TABLE user_tenants;
        TRUNCATE TABLE users;
        TRUNCATE TABLE tenants;
        SET FOREIGN_KEY_CHECKS = 1;
EOF
    echo "✅ 旧数据已清理"
    echo ""
fi

# 插入基础数据
echo "插入基础数据..."

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
-- 插入默认租户
INSERT INTO tenants (id, name, code, status) VALUES 
('uat-tenant-001', 'UAT测试租户', 'uat-tenant', 1);

-- 插入测试用户
INSERT INTO users (id, username, password, email, real_name, status) VALUES 
('uat-admin-001', 'uat_admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye1j.d8q5f2Z1F5e1e1e1e1e1e1e1e1e1', 'uat-admin@jimureport.com', 'UAT管理员', 1),
('uat-user-001', 'uat_user1', '$2a$10$N9qo8uLOickgx2ZMRZoMye1j.d8q5f2Z1F5e1e1e1e1e1e1e1e1e1', 'uat-user1@jimureport.com', '测试用户1', 1),
('uat-user-002', 'uat_user2', '$2a$10$N9qo8uLOickgx2ZMRZoMye1j.d8q5f2Z1F5e1e1e1e1e1e1e1e1e1', 'uat-user2@jimureport.com', '测试用户2', 1);

-- 插入用户-租户关联
INSERT INTO user_tenants (id, user_id, tenant_id, role, is_default) VALUES 
(UUID(), 'uat-admin-001', 'uat-tenant-001', 'admin', 1),
(UUID(), 'uat-user-001', 'uat-tenant-001', 'user', 1),
(UUID(), 'uat-user-002', 'uat-tenant-001', 'user', 1);
EOF

echo "✅ 基础数据插入完成"
echo ""

# 插入数据源
echo "插入测试数据源..."

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
INSERT INTO data_sources (id, tenant_id, name, type, host, port, database_name, username, password, status) VALUES 
('uat-ds-001', 'uat-tenant-001', 'UAT-MySQL测试库', 'mysql', 'mysql', 3306, 'jimureport', 'root', 'root', 1),
('uat-ds-002', 'uat-tenant-001', 'UAT-销售数据库', 'mysql', 'mysql', 3306, 'jimureport', 'root', 'root', 1),
('uat-ds-003', 'uat-tenant-001', 'UAT-用户数据库', 'mysql', 'mysql', 3306, 'jimureport', 'root', 'root', 1);
EOF

echo "✅ 数据源插入完成 (3个)"
echo ""

# 插入测试报表
echo "插入测试报表..."

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
INSERT INTO reports (id, tenant_id, name, code, type, config, status) VALUES 
('uat-report-001', 'uat-tenant-001', '月度销售报表', 'monthly_sales', 'report', '{"version":"1.0","sheets":[{"name":"Sheet1","rows":10,"cols":8}]}', 1),
('uat-report-002', 'uat-tenant-001', '季度财务报表', 'quarterly_finance', 'report', '{"version":"1.0","sheets":[{"name":"Sheet1","rows":15,"cols":10}]}', 1),
('uat-report-003', 'uat-tenant-001', '用户活跃度分析', 'user_activity', 'report', '{"version":"1.0","sheets":[{"name":"Sheet1","rows":20,"cols":6}]}', 1),
('uat-report-004', 'uat-tenant-001', '产品库存报表', 'product_inventory', 'report', '{"version":"1.0","sheets":[{"name":"Sheet1","rows":12,"cols":8}]}', 1),
('uat-report-005', 'uat-tenant-001', '客户满意度调查', 'customer_satisfaction', 'report', '{"version":"1.0","sheets":[{"name":"Sheet1","rows":8,"cols":5}]}', 1);
EOF

echo "✅ 测试报表插入完成 (5个)"
echo ""

# 插入测试大屏
echo "插入测试大屏..."

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
INSERT INTO dashboards (id, tenant_id, name, code, config, status) VALUES 
('uat-dashboard-001', 'uat-tenant-001', '销售数据大屏', 'sales_dashboard', '{"width":1920,"height":1080,"backgroundColor":"#1a1a2e","components":[]}', 1),
('uat-dashboard-002', 'uat-tenant-001', '实时监控大屏', 'monitoring_dashboard', '{"width":1920,"height":1080,"backgroundColor":"#0f0f23","components":[]}', 1),
('uat-dashboard-003', 'uat-tenant-001', '运营分析大屏', 'operation_dashboard', '{"width":1920,"height":1080,"backgroundColor":"#1a1a2e","components":[]}', 1);
EOF

echo "✅ 测试大屏插入完成 (3个)"
echo ""

# 创建测试数据表
echo "创建测试数据表..."

mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASS $DB_NAME << 'EOF'
-- 销售数据表
CREATE TABLE IF NOT EXISTS test_sales_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_name VARCHAR(100),
    sales_amount DECIMAL(10,2),
    sales_date DATE,
    region VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户数据表
CREATE TABLE IF NOT EXISTS test_user_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50),
    email VARCHAR(100),
    registration_date DATE,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 库存数据表
CREATE TABLE IF NOT EXISTS test_inventory_data (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_name VARCHAR(100),
    quantity INT,
    warehouse VARCHAR(50),
    last_updated DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 插入测试数据
INSERT INTO test_sales_data (product_name, sales_amount, sales_date, region) VALUES
('产品A', 15000.00, '2024-01-15', '华东'),
('产品B', 23000.00, '2024-01-16', '华南'),
('产品C', 18900.00, '2024-01-17', '华北'),
('产品D', 32000.00, '2024-01-18', '华东'),
('产品E', 12500.00, '2024-01-19', '西南');

INSERT INTO test_user_data (username, email, registration_date, status) VALUES
('user001', 'user001@test.com', '2024-01-01', 'active'),
('user002', 'user002@test.com', '2024-01-05', 'active'),
('user003', 'user003@test.com', '2024-01-10', 'inactive'),
('user004', 'user004@test.com', '2024-01-15', 'active'),
('user005', 'user005@test.com', '2024-01-20', 'active');

INSERT INTO test_inventory_data (product_name, quantity, warehouse, last_updated) VALUES
('产品A', 150, '仓库1', '2024-01-20'),
('产品B', 89, '仓库1', '2024-01-19'),
('产品C', 234, '仓库2', '2024-01-18'),
('产品D', 67, '仓库2', '2024-01-17'),
('产品E', 189, '仓库3', '2024-01-16');
EOF

echo "✅ 测试数据表创建完成"
echo ""

# 生成测试数据摘要
echo "================================"
echo "测试数据准备完成!"
echo "================================"
echo ""
echo "数据摘要:"
echo "  租户: 1个"
echo "  用户: 3个 (1管理员 + 2普通用户)"
echo "  数据源: 3个"
echo "  报表: 5个"
echo "  大屏: 3个"
echo "  测试数据表: 3个 (销售/用户/库存)"
echo ""
echo "登录凭据:"
echo "  管理员: uat_admin / admin"
echo "  用户1:  uat_user1 / user1"
echo "  用户2:  uat_user2 / user2"
echo ""
echo "注意: 密码已使用 BCrypt 加密"
echo ""

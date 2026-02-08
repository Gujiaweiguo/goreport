-- 数据集功能数据库迁移脚本
-- 添加数据集、数据集字段、数据集源表

USE goreport;

-- 数据集表
CREATE TABLE IF NOT EXISTS datasets (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(36) NOT NULL,
    datasource_id VARCHAR(36) COMMENT '关联的数据源ID，可选',
    name VARCHAR(200) NOT NULL,
    type ENUM('sql', 'api', 'file') NOT NULL COMMENT '数据集类型：SQL/API/文件导入',
    config JSON COMMENT '数据集配置：SQL查询/API配置/文件配置',
    status TINYINT DEFAULT 1 COMMENT '0-禁用 1-启用',
    created_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_tenant_id (tenant_id),
    INDEX idx_datasource_id (datasource_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (datasource_id) REFERENCES data_sources(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集定义';

-- 数据集字段表
CREATE TABLE IF NOT EXISTS dataset_fields (
    id VARCHAR(36) PRIMARY KEY,
    dataset_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL COMMENT '字段原始名称',
    display_name VARCHAR(100) COMMENT '字段显示名称',
    type ENUM('dimension', 'measure') NOT NULL COMMENT '字段类型：维度/指标',
    data_type ENUM('string', 'number', 'date', 'boolean') NOT NULL COMMENT '数据类型',
    is_computed BOOLEAN DEFAULT FALSE COMMENT '是否为计算字段',
    expression TEXT COMMENT '计算字段表达式，仅计算字段有值',
    is_sortable BOOLEAN DEFAULT TRUE COMMENT '是否可排序',
    is_groupable BOOLEAN DEFAULT TRUE COMMENT '是否可分组（仅维度）',
    default_sort_order ENUM('asc', 'desc', 'none') DEFAULT 'none' COMMENT '默认排序顺序',
    sort_index INT DEFAULT 0 COMMENT '排序索引',
    config JSON COMMENT '字段扩展配置',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_dataset_id (dataset_id),
    INDEX idx_type (type),
    INDEX idx_is_computed (is_computed),
    INDEX idx_sort_index (sort_index),
    FOREIGN KEY (dataset_id) REFERENCES datasets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集字段配置';

-- 数据集源表（用于复杂的多数据源场景，如联邦查询）
CREATE TABLE IF NOT EXISTS dataset_sources (
    id VARCHAR(36) PRIMARY KEY,
    dataset_id VARCHAR(36) NOT NULL,
    source_type ENUM('datasource', 'api', 'file') NOT NULL COMMENT '源类型',
    source_id VARCHAR(36) COMMENT '源ID（datasource_id 或 API endpoint ID）',
    source_config JSON COMMENT '源配置（SQL/API配置/文件路径）',
    join_type ENUM('inner', 'left', 'right', 'full') DEFAULT 'inner' COMMENT '连接类型',
    join_condition TEXT COMMENT '连接条件（仅SQL类型）',
    sort_index INT DEFAULT 0 COMMENT '排序索引',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_dataset_id (dataset_id),
    INDEX idx_source_type (source_type),
    INDEX idx_sort_index (sort_index),
    FOREIGN KEY (dataset_id) REFERENCES datasets(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据集数据源配置';

-- 添加注释
ALTER TABLE datasets COMMENT = '数据集定义：统一的数据抽象层，支持SQL、API、文件导入等多种数据源';
ALTER TABLE dataset_fields COMMENT = '数据集字段配置：包括维度、指标、计算字段的定义';
ALTER TABLE dataset_sources COMMENT = '数据集数据源配置：支持多数据源联合查询';

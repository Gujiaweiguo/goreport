# JimuReport 数据库表结构

## 1. 概述

JimuReport 使用 MySQL 5.7+ 作为数据存储，支持多租户模式。核心表分为以下几类：

- 报表相关表 (`jimu_report*`)
- 仪表盘/大屏相关表 (`onl_drag_*`)
- 字典表 (`jimu_dict*`)
- 示例数据表 (`rep_demo_*`, `test_*`)

## 2. 报表核心表

### 2.1 jimu_report (报表主表)

**用途**: 存储报表模板和配置

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(32) | 主键 |
| code | varchar(50) | 编码（唯一） |
| name | varchar(50) | 报表名称 |
| note | varchar(255) | 说明 |
| status | varchar(10) | 状态 |
| type | varchar(50) | 类型 |
| json_str | longtext | JSON 配置（核心字段，存储报表设计） |
| api_url | varchar(255) | API 请求地址 |
| thumb | text | 缩略图 |
| create_by | varchar(50) | 创建人 |
| create_time | datetime | 创建时间 |
| update_by | varchar(50) | 修改人 |
| update_time | datetime | 修改时间 |
| del_flag | tinyint(1) | 删除标识 (0-正常, 1-已删除) |
| api_method | varchar(255) | 请求方法 (0-get, 1-post) |
| template | tinyint(1) | 是否是模板 (0-是, 1-不是) |
| view_count | bigint | 浏览次数 |
| css_str | text | CSS 增强 |
| js_str | text | JS 增强 |
| py_str | text | Python 增强 |
| tenant_id | varchar(10) | 多租户标识 |
| update_count | int | 乐观锁版本 |
| submit_form | tinyint(1) | 是否填报报表 (0-不是, 1-是) |
| is_multi_sheet | tinyint | 是否多sheet报表 (1-是, 0-否) |

**索引**:
- PRIMARY: id
- UNIQUE: code
- INDEX: create_by, del_flag

### 2.2 jimu_report_data_source (数据源表)

**用途**: 存储报表和仪表盘的数据源配置

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(36) | 主键 |
| name | varchar(100) | 数据源名称 |
| report_id | varchar(100) | 关联的报表 ID |
| code | varchar(100) | 编码 |
| remark | varchar(200) | 备注 |
| db_type | varchar(10) | 数据库类型 (MYSQL5.7, FILES, etc.) |
| db_driver | varchar(100) | 驱动类 |
| db_url | varchar(500) | 数据源地址 |
| db_username | varchar(100) | 用户名 |
| db_password | varchar(100) | 密码（加密） |
| connect_times | int | 连接失败次数 |
| tenant_id | varchar(10) | 多租户标识 |
| type | varchar(10) | 类型 (report-报表, drag-仪表盘) |

**索引**:
- PRIMARY: id
- INDEX: report_id, code

### 2.3 jimu_report_db (报表数据库表)

**用途**: 存储报表关联的数据库配置

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(32) | 主键 |
| report_id | varchar(32) | 报表 ID |
| db_name | varchar(100) | 数据库名 |
| db_type | varchar(10) | 数据库类型 |
| db_key | varchar(50) | 数据库 Key（用于引用） |

### 2.4 jimu_report_export_job (导出任务表)

**用途**: 存储自动化导出任务配置

### 2.5 jimu_report_export_log (导出日志表)

**用途**: 记录导出操作日志

### 2.6 jimu_report_sheet (报表 Sheet 表)

**用途**: 存储多 Sheet 报表的 Sheet 配置

### 2.7 其他报表相关表

| 表名 | 用途 |
|------|------|
| jimu_report_category | 报表分类 |
| jimu_report_db_field | 报表数据库字段 |
| jimu_report_db_param | 报表数据库参数 |
| jimu_report_ext_data | 报表扩展数据 |
| jimu_report_icon_lib | 报表图标库 |
| jimu_report_link | 报表链接 |
| jimu_report_map | 报表映射 |
| jimu_report_share | 报表分享 |

## 3. 仪表盘/大屏核心表

### 3.1 onl_drag_page (仪表盘页面表)

**用途**: 存储大屏/仪表盘页面配置

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(50) | 主键 |
| name | varchar(100) | 界面名称 |
| path | varchar(100) | 访问路径 |
| background_color | varchar(10) | 背景色 |
| background_image | varchar(255) | 背景图 |
| design_type | int | 设计模式 (1-PC, 2-手机, 3-平板) |
| theme | varchar(10) | 主题色 |
| style | varchar(20) | 面板主题 |
| cover_url | varchar(500) | 封面图 |
| des_json | varchar(1000) | 仪表盘主配置 JSON |
| template | longtext | 布局 JSON（核心字段，存储组件布局） |
| protection_code | varchar(32) | 保护码 |
| type | varchar(64) | 文件夹类 |
| iz_template | varchar(10) | 是否模板 (1-是, 0-不是) |
| low_app_id | varchar(50) | 应用 ID |
| tenant_id | int | 租户 ID |
| visits_num | int | 访问次数 |
| del_flag | int | 删除状态 (0-未删除, 1-已删除) |

### 3.2 onl_drag_dataset_head (数据集头表)

**用途**: 存储仪表盘数据集配置

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(32) | 主键 |
| name | varchar(100) | 名称 |
| code | varchar(36) | 编码 |
| parent_id | varchar(36) | 父 ID |
| db_source | varchar(100) | 动态数据源 |
| query_sql | varchar(5000) | 查询数据 SQL |
| content | varchar(1000) | 描述 |
| data_type | varchar(50) | 数据类型 |
| api_method | varchar(10) | API 方法 (get/post) |
| low_app_id | varchar(32) | 应用 ID |
| tenant_id | int | 租户 ID |

### 3.3 onl_drag_dataset_item (数据集项表)

**用途**: 存储数据集的字段配置

### 3.4 onl_drag_dataset_param (数据集参数表)

**用途**: 存储数据集的参数配置

### 3.5 onl_drag_page_comp (页面组件表)

**用途**: 存储仪表盘页面组件的详细配置

### 3.6 onl_drag_comp (组件库表)

**用途**: 存储可用组件的定义

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(32) | 主键 |
| parent_id | varchar(32) | 父 ID |
| comp_name | varchar(50) | 组件名称 |
| comp_type | varchar(20) | 组件类型 |
| icon | varchar(50) | 图标 |
| order_num | int | 排序 |
| type_id | int | 组件类型 |
| comp_config | longtext | 组件配置 |
| status | varchar(2) | 状态 (0-无效, 1-有效) |

### 3.7 onl_drag_share (分享表)

**用途**: 存储仪表盘分享配置

### 3.8 onl_drag_table_relation (表关系表)

**用途**: 存储数据表之间的关系

## 4. 字典表

### 4.1 jimu_dict (字典主表)

**用途**: 存储字典定义

**核心字段**:
| 字段名 | 类型 | 说明 |
|--------|------|------|
| id | varchar(32) | 主键 |
| dict_name | varchar(100) | 字典名称 |
| dict_code | varchar(100) | 字典编码（唯一） |
| description | varchar(255) | 描述 |
| del_flag | int | 删除状态 |
| type | int(1) | 字典类型 (0-string, 1-number) |
| tenant_id | varchar(10) | 多租户标识 |

**索引**:
- PRIMARY: id
- UNIQUE: dict_code

### 4.2 jimu_dict_item (字典项表)

**用途**: 存储字典的具体项

## 5. 示例数据表

以下表仅用于演示，生产环境中应删除或替换：

| 表名 | 用途 |
|------|------|
| rep_demo_dxtj | 演示数据-统计分析 |
| rep_demo_employee | 演示数据-员工 |
| rep_demo_gongsi | 演示数据-公司 |
| rep_demo_jianpiao | 演示数据-发票 |
| rep_demo_xiaoshou | 演示数据-销售 |
| test_customer | 测试客户表 |
| test_order | 测试订单表 |
| test_order_pros | 测试订单产品表 |
| test_resume | 测试简历表 |

## 6. 临时表

| 表名 | 用途 |
|------|------|
| tmp_report_data_1 | 临时报表数据 |
| tmp_report_data_income | 临时收入数据 |

## 7. 多租户支持

多租户通过 `tenant_id` 字段实现，以下表支持多租户：

- jimu_report
- jimu_report_data_source
- jimu_dict
- onl_drag_page
- onl_drag_dataset_head

## 8. 字段命名约定

| 前缀 | 说明 |
|------|------|
| create_by | 创建人 |
| create_time | 创建时间 |
| update_by | 更新人 |
| update_time | 更新时间 |
| del_flag | 删除标识 |
| tenant_id | 租户 ID |

## 9. 数据类型说明

| 类型 | 说明 |
|------|------|
| varchar(32) | 短字符串（ID） |
| varchar(50) | 中等字符串（名称） |
| varchar(100) | 长字符串（描述） |
| varchar(255) | 超长字符串 |
| text | 大文本 |
| longtext | 超大文本（JSON 配置） |
| tinyint(1) | 布尔值 |
| int | 整数 |
| bigint | 大整数 |
| datetime | 日期时间 |

## 10. 核心表依赖关系

```
jimu_report (报表)
    ├── jimu_report_data_source (数据源)
    │   └── jimu_report_db (数据库)
    │       ├── jimu_report_db_field (字段)
    │       └── jimu_report_db_param (参数)
    ├── jimu_report_sheet (Sheet)
    ├── jimu_report_ext_data (扩展数据)
    └── jimu_report_share (分享)

onl_drag_page (仪表盘页面)
    ├── onl_drag_page_comp (组件)
    └── onl_drag_dataset_head (数据集)
        ├── onl_drag_dataset_item (数据项)
        └── onl_drag_dataset_param (参数)

onl_drag_comp (组件库)

jimu_dict (字典)
    └── jimu_dict_item (字典项)
```

## 11. Go 后端迁移注意事项

1. **必须保持表结构不变**: 现有数据库 schema 必须完全兼容
2. **多租户处理**: 需要在查询时考虑 `tenant_id` 过滤
3. **JSON 字段**: `json_str` 和 `template` 字段包含复杂的 JSON 配置，需要正确解析
4. **软删除**: 使用 `del_flag` 进行软删除，而不是物理删除
5. **乐观锁**: 使用 `update_count` 字段处理并发更新
6. **加密字段**: 数据源密码等敏感字段需要正确解密

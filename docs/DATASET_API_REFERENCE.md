# 数据集 API 文档

## 概述

本文档描述数据集（Dataset）相关的 REST API 端点。

**Base URL**: `/api/v1/datasets`

**认证**: 所有 API 请求需要在 Header 中包含 JWT token：
```
Authorization: Bearer <token>
```

## 数据模型

### Dataset

```json
{
  "id": "string",
  "tenantId": "string",
  "datasourceId": "string",
  "name": "string",
  "type": "sql|api|file",
  "config": "object",
  "status": 1,
  "createdBy": "string",
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

### DatasetField

```json
{
  "id": "string",
  "datasetId": "string",
  "name": "string",
  "displayName": "string",
  "type": "dimension|measure|computed",
  "dataType": "string|number|date|boolean",
  "isComputed": false,
  "expression": "string",
  "isSortable": true,
  "isGroupable": true,
  "defaultSortOrder": "asc|desc|none",
  "sortIndex": 0,
  "config": "object",
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

### DatasetSource

```json
{
  "id": "string",
  "datasetId": "string",
  "sourceType": "datasource|api|file",
  "sourceId": "string",
  "sourceConfig": "object",
  "joinType": "inner|left|right|full",
  "joinCondition": "string",
  "sortIndex": 0,
  "createdAt": "datetime",
  "updatedAt": "datetime"
}
```

## API 端点

### 1. 创建数据集

**POST** `/datasets`

创建新的数据集。

**请求体**:

```json
{
  "name": "string (required)",
  "type": "sql|api|file (required)",
  "datasourceId": "string (optional, required for SQL type)",
  "description": "string (optional)",
  "config": {
    // SQL type
    "sql": "SELECT * FROM sales WHERE create_time >= :start_date",
    "params": {
      "start_date": "2024-01-01"
    },

    // API type
    "url": "https://api.example.com/data",
    "method": "GET",
    "headers": {
      "Authorization": "Bearer token"
    },
    "timeout": 30000,

    // File type
    "filePath": "/path/to/data.csv",
    "format": "csv",
    "delimiter": ",",
    "encoding": "UTF-8"
  }
}
```

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "id": "uuid",
    "name": "Sales Dataset",
    "type": "sql",
    "status": 1,
    "createdAt": "2024-02-08T10:00:00Z",
    "updatedAt": "2024-02-08T10:00:00Z"
  },
  "message": "Dataset created successfully"
}
```

失败（400/500）:
```json
{
  "success": false,
  "message": "Error description",
  "timestamp": "2024-02-08T10:00:00Z"
}
```

### 2. 获取数据集列表

**GET** `/datasets`

分页获取数据集列表。

**查询参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| pageSize | int | 否 | 每页大小，默认 20 |
| keyword | string | 否 | 搜索关键词 |
| type | string | 否 | 数据集类型过滤（sql/api/file） |
| status | int | 否 | 状态过滤（1=active, 0=inactive） |

**请求示例**:

```
GET /api/v1/datasets?page=1&pageSize=20&type=sql&status=1
```

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": [
    {
      "id": "uuid-1",
      "name": "Sales Dataset",
      "type": "sql",
      "description": "Sales data from 2024",
      "status": 1,
      "createdAt": "2024-02-08T10:00:00Z",
      "updatedAt": "2024-02-08T10:00:00Z",
      "fields": [
        {
          "id": "field-1",
          "name": "region",
          "type": "dimension",
          "dataType": "string"
        }
      ]
    }
  ],
  "total": 100,
  "page": 1,
  "pageSize": 20
}
```

### 3. 获取数据集详情

**GET** `/datasets/{id}`

获取指定数据集的详细信息。

**路径参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| id | string | 是 | 数据集 ID |

**请求示例**:

```
GET /api/v1/datasets/{dataset-uuid}
```

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "id": "uuid",
    "name": "Sales Dataset",
    "type": "sql",
    "description": "Sales data from 2024",
    "config": {
      "sql": "SELECT * FROM sales"
    },
    "status": 1,
    "fields": [...],
    "sources": [...],
    "createdAt": "2024-02-08T10:00:00Z",
    "updatedAt": "2024-02-08T10:00:00Z"
  }
}
```

失败（404）:
```json
{
  "success": false,
  "message": "Dataset not found"
}
```

### 4. 获取数据集 Schema

**GET** `/datasets/{id}/schema`

获取数据集的维度和指标列表。

**路径参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| id | string | 是 | 数据集 ID |

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "dimensions": [
      {
        "id": "field-1",
        "name": "region",
        "displayName": "地区",
        "type": "dimension",
        "dataType": "string",
        "isSortable": true,
        "isGroupable": true
      }
    ],
    "measures": [
      {
        "id": "field-2",
        "name": "sales_amount",
        "displayName": "销售额",
        "type": "measure",
        "dataType": "number",
        "defaultAggregation": "SUM"
      }
    ]
  }
}
```

### 5. 查询数据集

**POST** `/datasets/{id}/query`

执行数据集查询，支持筛选、排序和分组。

**路径参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| id | string | 是 | 数据集 ID |

**请求体**:

```json
{
  "dimensions": ["region", "product"],
  "measures": ["sales_amount", "quantity"],
  "filters": [
    {
      "field": "region",
      "operator": "=",
      "value": "华东"
    },
    {
      "field": "sales_amount",
      "operator": ">",
      "value": 10000
    }
  ],
  "sortBy": [
    {
      "field": "sales_amount",
      "order": "desc"
    }
  ],
  "groupBy": ["region"],
  "aggregation": {
    "sales_amount": "SUM",
    "quantity": "SUM"
  },
  "limit": 100,
  "offset": 0
}
```

**参数说明**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| dimensions | array | 否 | 维度字段列表 |
| measures | array | 是 | 指标字段列表 |
| filters | array | 否 | 过滤条件列表 |
| filters[].field | string | 是 | 字段名 |
| filters[].operator | string | 是 | 操作符（=, !=, >, <, >=, <=, LIKE, IN） |
| filters[].value | any | 是 | 过滤值 |
| sortBy | array | 否 | 排序规则 |
| sortBy[].field | string | 是 | 排序字段 |
| sortBy[].order | string | 是 | 排序方向（asc, desc） |
| groupBy | array | 否 | 分组字段 |
| aggregation | object | 否 | 聚合函数映射（字段名 -> 函数名） |
| limit | int | 否 | 返回记录数限制 |
| offset | int | 否 | 偏移量 |

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "data": [
      {
        "region": "华东",
        "sales_amount": 150000,
        "quantity": 1500
      },
      {
        "region": "华南",
        "sales_amount": 120000,
        "quantity": 1200
      }
    ],
    "total": 2,
    "limit": 100,
    "offset": 0
  }
}
```

### 6. 更新数据集

**PUT** `/datasets/{id}`

更新数据集信息。

**路径参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| id | string | 是 | 数据集 ID |

**请求体**:

```json
{
  "name": "string (optional)",
  "description": "string (optional)",
  "config": {
    // 同创建请求的 config
  },
  "status": 1
}
```

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "id": "uuid",
    "name": "Updated Dataset Name",
    "type": "sql",
    "status": 1,
    "updatedAt": "2024-02-08T11:00:00Z"
  },
  "message": "Dataset updated successfully"
}
```

失败（404）:
```json
{
  "success": false,
  "message": "Dataset not found"
}
```

### 7. 删除数据集

**DELETE** `/datasets/{id}`

删除指定数据集（软删除）。

**路径参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| id | string | 是 | 数据集 ID |

**响应**:

成功（200）:
```json
{
  "success": true,
  "message": "Dataset deleted successfully"
}
```

失败（404）:
```json
{
  "success": false,
  "message": "Dataset not found"
}
```

## 字段管理 API

### 8. 创建字段

**POST** `/datasets/{datasetId}/fields`

创建数据集字段。

**请求体**:

```json
{
  "name": "string (required)",
  "displayName": "string (optional)",
  "type": "dimension|measure|computed (required)",
  "dataType": "string|number|date|boolean (required)",
  "isComputed": false,
  "expression": "string (required if isComputed=true)",
  "isSortable": true,
  "isGroupable": true,
  "defaultSortOrder": "asc|desc|none",
  "config": {}
}
```

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "id": "field-uuid",
    "datasetId": "dataset-uuid",
    "name": "growth_rate",
    "displayName": "增长率",
    "type": "measure",
    "dataType": "number",
    "isComputed": true,
    "expression": "(sales - last_year_sales) / last_year_sales * 100"
  }
}
```

### 9. 更新字段

**PUT** `/datasets/{datasetId}/fields/{fieldId}`

更新数据集字段。

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "id": "field-uuid",
    "name": "updated_field_name",
    ...
  },
  "message": "Field updated successfully"
}
```

### 10. 删除字段

**DELETE** `/datasets/{datasetId}/fields/{fieldId}`

删除数据集字段。

**响应**:

成功（200）:
```json
{
  "success": true,
  "message": "Field deleted successfully"
}
```

### 11. 获取字段列表

**GET** `/datasets/{datasetId}/fields`

获取数据集的所有字段。

**查询参数**:

| 参数 | 类型 | 必需 | 描述 |
|------|------|------|------|
| type | string | 否 | 字段类型过滤（dimension/measure/computed） |

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": [
    {
      "id": "field-1",
      "name": "region",
      "type": "dimension",
      "dataType": "string"
    },
    {
      "id": "field-2",
      "name": "sales_amount",
      "type": "measure",
      "dataType": "number"
    }
  ]
}
```

## 预览 API

### 12. 预览数据

**POST** `/datasets/{id}/preview`

预览数据集数据（返回前 100 条）。

**响应**:

成功（200）:
```json
{
  "success": true,
  "result": {
    "columns": ["region", "product", "sales_amount", "quantity"],
    "data": [
      {"region": "华东", "product": "A", "sales_amount": 1000, "quantity": 10},
      {"region": "华南", "product": "B", "sales_amount": 2000, "quantity": 20}
    ],
    "total": 100
  }
}
```

## 错误代码

| HTTP 状态码 | 错误类型 | 描述 |
|-------------|----------|------|
| 400 | Bad Request | 请求参数错误 |
| 401 | Unauthorized | 未授权或 token 失效 |
| 403 | Forbidden | 无权限访问 |
| 404 | Not Found | 资源不存在 |
| 409 | Conflict | 资源冲突（如名称重复） |
| 500 | Internal Server Error | 服务器内部错误 |

**错误响应格式**:

```json
{
  "success": false,
  "message": "Error description",
  "code": "ERROR_CODE",
  "timestamp": "2024-02-08T10:00:00Z"
}
```

## 限流

API 实施了速率限制：

- 每个用户每分钟最多 100 次请求
- 超过限制将返回 429 状态码
- 响应头包含限流信息：
  ```
  X-RateLimit-Limit: 100
  X-RateLimit-Remaining: 95
  X-RateLimit-Reset: 1644520000
  ```

## 示例代码

### Python 示例

```python
import requests

base_url = "http://localhost:8085/api/v1"
token = "your-jwt-token"

headers = {
    "Authorization": f"Bearer {token}",
    "Content-Type": "application/json"
}

# 创建数据集
data = {
    "name": "Sales Dataset",
    "type": "sql",
    "datasourceId": "datasource-uuid",
    "config": {
        "sql": "SELECT * FROM sales"
    }
}

response = requests.post(
    f"{base_url}/datasets",
    json=data,
    headers=headers
)

print(response.json())
```

### JavaScript 示例

```javascript
const baseURL = 'http://localhost:8085/api/v1';
const token = 'your-jwt-token';

const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
};

// 获取数据集列表
fetch(`${baseURL}/datasets?page=1&pageSize=20`, {
    headers: headers
})
    .then(response => response.json())
    .then(data => console.log(data));
```

### cURL 示例

```bash
# 获取数据集列表
curl -X GET "http://localhost:8085/api/v1/datasets?page=1&pageSize=20" \
  -H "Authorization: Bearer your-jwt-token"

# 创建数据集
curl -X POST "http://localhost:8085/api/v1/datasets" \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sales Dataset",
    "type": "sql",
    "datasourceId": "datasource-uuid",
    "config": {"sql": "SELECT * FROM sales"}
  }'
```

## 版本

当前 API 版本：v1

API 版本变更将在发布说明中标注，遵循 [语义化版本控制](https://semver.org/)。

## 支持

如有问题，请：
- 提交 [GitHub Issue](https://github.com/gujiaweiguo/goreport/issues)
- 查阅 [用户指南](./DATASET_USER_GUIDE.md)
- 联系技术支持：[weiguogu@163.com](mailto:weiguogu@163.com)

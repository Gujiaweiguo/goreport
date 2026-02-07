# JWT Claim 约定和认证规范

## 1. 概述

Go 后端将使用 JWT (JSON Web Token) 进行认证，替代原有的 Sa-Token 认证机制。为了保持与现有 Go 系统的兼容性，我们将使用标准的 JWT claim 结构。

## 2. JWT Claim 结构

### 2.1 标准 Claims

| Claim | 类型 | 说明 | 示例 |
|-------|------|------|------|
| sub (Subject) | string | 用户 ID | "admin" |
| iss (Issuer) | string | 签发者 | "jimureport-go" |
| aud (Audience) | string | 接收方 | "jimureport-api" |
| exp (Expiration) | int64 | 过期时间 (Unix timestamp) | 1735689600 |
| iat (Issued At) | int64 | 签发时间 (Unix timestamp) | 1735689600 |
| nbf (Not Before) | int64 | 生效时间 (Unix timestamp) | 1735689600 |

### 2.2 自定义 Claims

| Claim | 类型 | 说明 | 示例 | 必需 |
|-------|------|------|------|------|
| username | string | 用户名 | "admin" | 是 |
| roles | string[] | 用户角色数组 | ["admin", "lowdeveloper"] | 否 |
| tenant_id | string | 租户 ID | "1" | 否 |

### 2.3 JWT Payload 示例

```json
{
  "sub": "admin",
  "username": "admin",
  "roles": ["admin", "lowdeveloper", "dbadeveloper"],
  "tenant_id": "1",
  "iss": "jimureport-go",
  "aud": "jimureport-api",
  "iat": 1735689600,
  "exp": 1738281600
}
```

## 3. Token 传递方式

Token 可以通过以下三种方式传递到后端，按优先级顺序处理：

### 3.1 优先级顺序

1. **Authorization Header** (标准方式)
   ```
   Authorization: Bearer <token>
   ```

2. **X-Access-Token Header** (兼容原有 Sa-Token)
   ```
   X-Access-Token: <token>
   ```

3. **Query Parameter** (用于 iframe 嵌入)
   ```
   ?token=<token>
   ```

### 3.2 Token 解析逻辑

```go
func GetToken(r *http.Request) string {
    // 1. 尝试从 Authorization header 获取
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
        return strings.TrimPrefix(authHeader, "Bearer ")
    }

    // 2. 尝试从 X-Access-Token header 获取
    token := r.Header.Get("X-Access-Token")
    if token != "" {
        return token
    }

    // 3. 尝试从 query parameter 获取
    return r.URL.Query().Get("token")
}
```

## 4. 租户 ID 传递

租户 ID 可以通过以下方式传递：

| 方式 | Header/Param | 示例 |
|------|---------------|------|
| Header | X-Tenant-Key | "1" |
| Header | X-Tenant-Id | "1" |
| Query Parameter | tenant_id | ?tenant_id=1 |

### 4.1 租户 ID 解析逻辑

```go
func GetTenantID(r *http.Request) string {
    // 1. 尝试从 X-Tenant-Key header 获取
    tenantID := r.Header.Get("X-Tenant-Key")
    if tenantID != "" {
        return tenantID
    }

    // 2. 尝试从 X-Tenant-Id header 获取
    tenantID = r.Header.Get("X-Tenant-Id")
    if tenantID != "" {
        return tenantID
    }

    // 3. 尝试从 query parameter 获取
    return r.URL.Query().Get("tenant_id")
}
```

## 5. 角色和权限

### 5.1 内置角色

goReport 内置三个角色，JWT 中的 `roles` claim 应包含以下值：

| 角色 | 说明 |
|------|------|
| admin | 管理员 |
| lowdeveloper | 低代码开发者 |
| dbadeveloper | 数据库开发者 |

### 5.2 权限标识符

权限标识符用于控制对特定功能的访问：

| 权限标识符 | 说明 |
|-----------|------|
| drag:datasource:testConnection | 仪表盘数据库连接测试 |
| onl:drag:clear:recovery | 清空回收站 |
| drag:analysis:sql | SQL 解析 |
| drag:design:getTotalData | 仪表盘对 Online 表单展示数据 |
| onl:drag:page:delete | 仪表盘页面删除 |
| drag:dataset:save | 数据集保存 |
| drag:dataset:delete | 数据集删除 |
| drag:datasource:saveOrUpate | 数据源保存/更新 |
| drag:datasource:delete | 数据源删除 |

### 5.3 权限检查逻辑

```go
func HasPermission(token string, requiredPermission string) bool {
    claims, err := ParseToken(token)
    if err != nil {
        return false
    }

    // admin 角色拥有所有权限
    for _, role := range claims.Roles {
        if role == "admin" {
            return true
        }
    }

    // 检查用户是否有特定权限
    return hasPermission(claims.Username, requiredPermission)
}
```

## 6. Token 验证

### 6.1 验证步骤

1. **提取 Token**: 从请求中提取 JWT token
2. **验证签名**: 使用密钥验证 JWT 签名
3. **检查过期时间**: 验证 `exp` claim 是否未过期
4. **解析 Claims**: 提取用户信息和角色
5. **检查权限**: (可选) 验证用户是否有访问权限

### 6.2 验证失败处理

| 错误类型 | HTTP Status Code | 描述 |
|----------|-----------------|------|
| Token 缺失 | 401 Unauthorized | 请求中未包含 token |
| Token 无效 | 401 Unauthorized | Token 签名验证失败 |
| Token 过期 | 401 Unauthorized | Token 已过期 |
| 权限不足 | 403 Forbidden | 用户没有访问权限 |

### 6.3 Token 验证中间件示例

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. 获取 token
        token := GetToken(r)
        if token == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        // 2. 验证 token
        claims, err := ValidateToken(token)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // 3. 将用户信息存入上下文
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 7. 配置项

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| JWT_SECRET | JWT 签名密钥 | - |
| JWT_EXPIRATION | Token 过期时间（秒） | 2592000 (30天) |
| TOKEN_NAME | Token 名称 | X-Access-Token |

## 8. 安全注意事项

1. **密钥管理**: JWT 密钥必须安全存储，建议使用环境变量
2. **HTTPS**: 生产环境必须使用 HTTPS 传输 token
3. **Token 过期**: 设置合理的 token 过期时间
4. **Token 刷新**: 实现 token 刷新机制
5. **权限检查**: 不要仅依赖 JWT 声明中的角色，每次请求都应验证权限

## 9. 与现有 Sa-Token 的兼容性

### 9.1 兼容策略

| 功能 | Sa-Token | JWT (Go) |
|------|-----------|------------|
| Token 存储 | Session/Cookie | 无状态 |
| Token 验证 | 服务端验证 | JWT 签名验证 |
| Token 传递 | Cookie/Query | Header/Query |
| 过期控制 | 配置 timeout | exp claim |

### 9.2 迁移步骤

1. **Dual-Auth 阶段**: 同时支持 Sa-Token 和 JWT
2. **逐步迁移**: 逐步将客户端切换到 JWT
3. **完全迁移**: 移除 Sa-Token 支持

## 10. Go 实现示例

### 10.1 JWT 生成

```go
func GenerateToken(user User) (string, error) {
    claims := jwt.MapClaims{
        "sub":      user.ID,
        "username":  user.Username,
        "roles":     user.Roles,
        "tenant_id": user.TenantID,
        "iss":       "jimureport-go",
        "aud":       "jimureport-api",
        "iat":       time.Now().Unix(),
        "exp":       time.Now().Add(30 * 24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
```

### 10.2 JWT 验证

```go
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return &Claims{
            Subject:   claims["sub"].(string),
            Username:  claims["username"].(string),
            Roles:     claims["roles"].([]string),
            TenantID:  claims["tenant_id"].(string),
        }, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

## 11. 测试用例

### 11.1 Token 生成测试

```bash
# 生成测试 token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

### 11.2 Token 使用测试

```bash
# 使用 Authorization header
curl http://localhost:8080/jmreport/list \
  -H "Authorization: Bearer <token>"

# 使用 X-Access-Token header
curl http://localhost:8080/jmreport/list \
  -H "X-Access-Token: <token>"

# 使用 query parameter
curl http://localhost:8080/jmreport/list?token=<token>
```

### 11.3 租户 ID 传递测试

```bash
# 使用 header
curl http://localhost:8080/jmreport/list \
  -H "X-Access-Token: <token>" \
  -H "X-Tenant-Id: 1"

# 使用 query parameter
curl http://localhost:8080/jmreport/list?token=<token>&tenant_id=1
```

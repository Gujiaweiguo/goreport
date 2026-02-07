# 模块接口定义

## 1. 概述

本文档定义了 goReport Go 后端各个模块的接口契约。每个模块包含 Handler（HTTP 层）、Service（业务逻辑层）、Repository（数据访问层）和 Models（数据模型）。

## 2. Auth 模块

### 2.1 职责

- JWT token 生成和验证
- 用户认证和授权
- 角色和权限管理
- 租户识别

### 2.2 Models

```go
package auth

// User 用户模型
type User struct {
    ID        string   `json:"id" gorm:"primaryKey"`
    Username  string   `json:"username" gorm:"uniqueIndex"`
    Password  string   `json:"-" gorm:"column:password"` // 存储加密后的密码
    Roles     []string `json:"roles"`
    TenantID  *string  `json:"tenant_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Claims JWT Claims
type Claims struct {
    Subject   string   `json:"sub"`
    Username  string   `json:"username"`
    Roles     []string `json:"roles"`
    TenantID  *string  `json:"tenant_id"`
    Issuer    string   `json:"iss"`
    Audience  string   `json:"aud"`
    ExpiresAt int64    `json:"exp"`
    IssuedAt  int64    `json:"iat"`
}

// LoginRequest 登录请求
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Token string `json:"token"`
    User  *User  `json:"user"`
}

// RefreshTokenRequest 刷新 Token 请求
type RefreshTokenRequest struct {
    Token string `json:"token" binding:"required"`
}

// RefreshTokenResponse 刷新 Token 响应
type RefreshTokenResponse struct {
    Token string `json:"token"`
}
```

### 2.3 Repository

```go
package auth

import "gorm.io/gorm"

type UserRepository interface {
    // FindByUsername 根据用户名查找用户
    FindByUsername(username string) (*User, error)

    // FindByID 根据 ID 查找用户
    FindByID(id string) (*User, error)

    // Create 创建用户
    Create(user *User) error

    // Update 更新用户
    Update(user *User) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(username string) (*User, error) {
    var user User
    err := r.db.Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) FindByID(id string) (*User, error) {
    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) Update(user *User) error {
    return r.db.Save(user).Error
}
```

### 2.4 Service

```go
package auth

type AuthService interface {
    // Login 用户登录
    Login(req *LoginRequest) (*LoginResponse, error)

    // RefreshToken 刷新 token
    RefreshToken(token string) (*RefreshTokenResponse, error)

    // ValidateToken 验证 token
    ValidateToken(tokenString string) (*Claims, error)

    // GetCurrentUser 获取当前用户（从上下文）
    GetCurrentUser(c *gin.Context) (*User, error)

    // HasRole 检查用户是否拥有指定角色
    HasRole(user *User, role string) bool

    // HasPermission 检查用户是否拥有指定权限
    HasPermission(user *User, permission string) bool

    // GetTenantID 获取租户 ID（从上下文）
    GetTenantID(c *gin.Context) *string
}

type authService struct {
    userRepo UserRepository
    jwtSecret string
    jwtExpiration int64 // seconds
}

func NewAuthService(userRepo UserRepository, jwtSecret string, jwtExpiration int64) AuthService {
    return &authService{
        userRepo:     userRepo,
        jwtSecret:     jwtSecret,
        jwtExpiration: jwtExpiration,
    }
}

func (s *authService) Login(req *LoginRequest) (*LoginResponse, error) {
    // 1. 查找用户
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // 2. 验证密码（需要实现密码加密/验证）
    if !s.verifyPassword(user.Password, req.Password) {
        return nil, ErrInvalidCredentials
    }

    // 3. 生成 JWT token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token: token,
        User:  user,
    }, nil
}

func (s *authService) RefreshToken(tokenString string) (*RefreshTokenResponse, error) {
    // 1. 验证旧 token
    claims, err := s.ValidateToken(tokenString)
    if err != nil {
        return nil, err
    }

    // 2. 根据用户信息生成新 token
    user, err := s.userRepo.FindByID(claims.Subject)
    if err != nil {
        return nil, ErrUserNotFound
    }

    newToken, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }

    return &RefreshTokenResponse{
        Token: newToken,
    }, nil
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
    // 实现 JWT 验证逻辑
    // ...
}

func (s *authService) GetCurrentUser(c *gin.Context) (*User, error) {
    claims, exists := c.Get("claims")
    if !exists {
        return nil, ErrUnauthorized
    }

    userClaims := claims.(*Claims)
    return s.userRepo.FindByID(userClaims.Subject)
}

func (s *authService) HasRole(user *User, role string) bool {
    for _, r := range user.Roles {
        if r == "admin" || r == role {
            return true
        }
    }
    return false
}

func (s *authService) HasPermission(user *User, permission string) bool {
    // 检查角色权限映射
    // TODO: 实现权限系统
    return s.HasRole(user, "admin")
}

func (s *authService) GetTenantID(c *gin.Context) *string {
    tenantID, exists := c.Get("tenant_id")
    if !exists {
        return nil
    }
    return tenantID.(*string)
}

func (s *authService) generateToken(user *User) (string, error) {
    // 实现 JWT 生成逻辑
    // ...
}

func (s *authService) verifyPassword(hashedPassword, password string) bool {
    // 实现密码验证逻辑
    // ...
}

// 错误定义
var (
    ErrInvalidCredentials = errors.New("invalid username or password")
    ErrUserNotFound     = errors.New("user not found")
    ErrUnauthorized     = errors.New("unauthorized")
)
```

### 2.5 Handler

```go
package auth

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type AuthHandler interface {
    // Login 用户登录
    Login(c *gin.Context)

    // RefreshToken 刷新 token
    RefreshToken(c *gin.Context)

    // Logout 用户登出（可选，JWT 无状态）
    Logout(c *gin.Context)
}

type authHandler struct {
    authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
    return &authHandler{authService: authService}
}

func (h *authHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    resp, err := h.authService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *authHandler) RefreshToken(c *gin.Context) {
    var req RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request body",
        })
        return
    }

    resp, err := h.authService.RefreshToken(req.Token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *authHandler) Logout(c *gin.Context) {
    // JWT 是无状态的，客户端只需要删除 token
    c.JSON(http.StatusOK, gin.H{
        "message": "Logged out successfully",
    })
}
```

## 3. Report 模块

### 3.1 职责

- 报表模板 CRUD
- 报表渲染和预览
- 报表分类管理
- 报表分享管理

### 3.2 Models

```go
package report

// Report 报表模型
type Report struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Code        string    `json:"code" gorm:"uniqueIndex;column:code"`
    Name        string    `json:"name" gorm:"column:name"`
    Note        *string   `json:"note" gorm:"column:note"`
    Status      *string   `json:"status" gorm:"column:status"`
    Type        *string   `json:"type" gorm:"column:type"`
    JSONStr     string    `json:"json_str" gorm:"column:json_str;type:longtext"` // 报表配置 JSON
    APIURL      *string   `json:"api_url" gorm:"column:api_url"`
    Thumb       *string   `json:"thumb" gorm:"column:thumb;type:text"`
    CreateBy    string    `json:"create_by" gorm:"column:create_by"`
    CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
    UpdateBy    string    `json:"update_by" gorm:"column:update_by"`
    UpdateTime  time.Time `json:"update_time" gorm:"column:update_time"`
    DelFlag     *int      `json:"del_flag" gorm:"column:del_flag"`
    APIMethod   *string   `json:"api_method" gorm:"column:api_method"`
    APICode     *string   `json:"api_code" gorm:"column:api_code"`
    Template    *int      `json:"template" gorm:"column:template"`
    ViewCount   int64     `json:"view_count" gorm:"column:view_count"`
    CSSStr      *string   `json:"css_str" gorm:"column:css_str;type:text"`
    JSStr       *string   `json:"js_str" gorm:"column:js_str;type:text"`
    PYStr       *string   `json:"py_str" gorm:"column:py_str;type:text"`
    TenantID    *string   `json:"tenant_id" gorm:"column:tenant_id"`
    UpdateCount int       `json:"update_count" gorm:"column:update_count"`
    SubmitForm  *int      `json:"submit_form" gorm:"column:submit_form"`
    IsMultiSheet *int     `json:"is_multi_sheet" gorm:"column:is_multi_sheet"`
}

// ReportCategory 报表分类
type ReportCategory struct {
    ID         string `json:"id" gorm:"primaryKey"`
    Name       string `json:"name"`
    ParentID   *string `json:"parent_id"`
    CreateBy   string `json:"create_by"`
    CreateTime time.Time `json:"create_time"`
    UpdateBy   string `json:"update_by"`
    UpdateTime time.Time `json:"update_time"`
    DelFlag    *int   `json:"del_flag"`
    TenantID    *string `json:"tenant_id"`
}

// ReportShare 报表分享
type ReportShare struct {
    ID         string `json:"id" gorm:"primaryKey"`
    ReportID   string `json:"report_id" gorm:"column:report_id"`
    ShareType  int    `json:"share_type"`
    CreateBy   string `json:"create_by"`
    CreateTime time.Time `json:"create_time"`
    TenantID    *string `json:"tenant_id"`
}

// CreateReportRequest 创建报表请求
type CreateReportRequest struct {
    Code       string `json:"code" binding:"required"`
    Name       string `json:"name" binding:"required"`
    Note       *string `json:"note"`
    Type       *string `json:"type"`
    JSONStr    string `json:"json_str" binding:"required"`
    CategoryID *string `json:"category_id"`
}

// UpdateReportRequest 更新报表请求
type UpdateReportRequest struct {
    ID       string `json:"id" binding:"required"`
    Code     *string `json:"code"`
    Name     *string `json:"name"`
    Note     *string `json:"note"`
    Status   *string `json:"status"`
    JSONStr  *string `json:"json_str"`
    CSSStr   *string `json:"css_str"`
    JSStr    *string `json:"js_str"`
    PYStr    *string `json:"py_str"`
}

// ListReportRequest 查询报表列表请求
type ListReportRequest struct {
    Page     int    `json:"page" binding:"min=1"`
    PageSize int    `json:"page_size" binding:"min=1,max=100"`
    Keyword *string `json:"keyword"`
    CategoryID *string `json:"category_id"`
    DelFlag *int `json:"del_flag"`
}

// RenderReportRequest 渲染报表请求
type RenderReportRequest struct {
    ID     string                 `json:"id" binding:"required"`
    Params map[string]interface{} `json:"params"`
}

// RenderReportResponse 渲染报表响应
type RenderReportResponse struct {
    HTML  string `json:"html"`
    Data  interface{} `json:"data"`
}
```

### 3.3 Repository

```go
package report

import "gorm.io/gorm"

type ReportRepository interface {
    // Create 创建报表
    Create(report *Report) error

    // Update 更新报表
    Update(report *Report) error

    // Delete 软删除报表
    Delete(id string) error

    // FindByID 根据 ID 查找报表
    FindByID(id string) (*Report, error)

    // FindByCode 根据 code 查找报表
    FindByCode(code string) (*Report, error)

    // List 查询报表列表
    List(params *ListReportRequest) ([]*Report, int64, error)

    // IncrementViewCount 增加浏览次数
    IncrementViewCount(id string) error
}

type reportRepository struct {
    db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
    return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *Report) error {
    return r.db.Create(report).Error
}

func (r *reportRepository) Update(report *Report) error {
    return r.db.Model(&Report{}).Where("id = ?", report.ID).Updates(report).Error
}

func (r *reportRepository) Delete(id string) error {
    return r.db.Model(&Report{}).Where("id = ?", id).Update("del_flag", 1).Error
}

func (r *reportRepository) FindByID(id string) (*Report, error) {
    var report Report
    err := r.db.Where("id = ? AND del_flag = 0", id).First(&report).Error
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *reportRepository) FindByCode(code string) (*Report, error) {
    var report Report
    err := r.db.Where("code = ? AND del_flag = 0", code).First(&report).Error
    if err != nil {
        return nil, err
    }
    return &report, nil
}

func (r *reportRepository) List(params *ListReportRequest) ([]*Report, int64, error) {
    var reports []*Report
    var total int64

    query := r.db.Model(&Report{}).Where("del_flag = ?", 0)

    if params.Keyword != nil && *params.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+*params.Keyword+"%")
    }

    if params.CategoryID != nil {
        query = query.Where("category_id = ?", *params.CategoryID)
    }

    if params.DelFlag != nil {
        query = query.Where("del_flag = ?", *params.DelFlag)
    }

    // 查询总数
    query.Count(&total)

    // 分页查询
    offset := (params.Page - 1) * params.PageSize
    err := query.Offset(offset).Limit(params.PageSize).Find(&reports).Error

    return reports, total, err
}

func (r *reportRepository) IncrementViewCount(id string) error {
    return r.db.Model(&Report{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}
```

### 3.4 Service

```go
package report

type ReportService interface {
    // Create 创建报表
    Create(req *CreateReportRequest, user *User) (*Report, error)

    // Update 更新报表
    Update(req *UpdateReportRequest, user *User) (*Report, error)

    // Delete 删除报表
    Delete(id string, user *User) error

    // Get 获取报表详情
    Get(id string, user *User) (*Report, error)

    // List 查询报表列表
    List(params *ListReportRequest, user *User) (*ReportListResponse, error)

    // Render 渲染报表
    Render(req *RenderReportRequest, user *User) (*RenderReportResponse, error)
}

type ReportListResponse struct {
    List  []*Report `json:"list"`
    Total int64      `json:"total"`
    Page  int        `json:"page"`
    PageSize int     `json:"page_size"`
}

type reportService struct {
    reportRepo ReportRepository
}

func NewReportService(reportRepo ReportRepository) ReportService {
    return &reportService{reportRepo: reportRepo}
}

func (s *reportService) Create(req *CreateReportRequest, user *User) (*Report, error) {
    report := &Report{
        ID:         generateID(),
        Code:       req.Code,
        Name:       req.Name,
        Note:       req.Note,
        Type:       req.Type,
        JSONStr:    req.JSONStr,
        CreateBy:   user.ID,
        CreateTime:  time.Now(),
        DelFlag:     ptrInt(0),
        ViewCount:   0,
        UpdateCount: 0,
        TenantID:    user.TenantID,
    }

    if err := s.reportRepo.Create(report); err != nil {
        return nil, err
    }

    return report, nil
}

func (s *reportService) Update(req *UpdateReportRequest, user *User) (*Report, error) {
    // 1. 查找报表
    report, err := s.reportRepo.FindByID(req.ID)
    if err != nil {
        return nil, ErrReportNotFound
    }

    // 2. 检查权限（TODO）

    // 3. 更新字段
    if req.Code != nil {
        report.Code = *req.Code
    }
    if req.Name != nil {
        report.Name = *req.Name
    }
    if req.JSONStr != nil {
        report.JSONStr = *req.JSONStr
    }
    report.UpdateBy = user.ID
    report.UpdateTime = time.Now()
    report.UpdateCount++

    if err := s.reportRepo.Update(report); err != nil {
        return nil, err
    }

    return report, nil
}

func (s *reportService) Delete(id string, user *User) error {
    // 1. 查找报表
    report, err := s.reportRepo.FindByID(id)
    if err != nil {
        return ErrReportNotFound
    }

    // 2. 检查权限（TODO）

    // 3. 软删除
    return s.reportRepo.Delete(id)
}

func (s *reportService) Get(id string, user *User) (*Report, error) {
    report, err := s.reportRepo.FindByID(id)
    if err != nil {
        return nil, ErrReportNotFound
    }

    // 增加浏览次数
    s.reportRepo.IncrementViewCount(id)

    return report, nil
}

func (s *reportService) List(params *ListReportRequest, user *User) (*ReportListResponse, error) {
    reports, total, err := s.reportRepo.List(params)
    if err != nil {
        return nil, err
    }

    return &ReportListResponse{
        List:     reports,
        Total:     total,
        Page:      params.Page,
        PageSize:  params.PageSize,
    }, nil
}

func (s *reportService) Render(req *RenderReportRequest, user *User) (*RenderReportResponse, error) {
    // 1. 查找报表
    report, err := s.reportRepo.FindByID(req.ID)
    if err != nil {
        return nil, ErrReportNotFound
    }

    // 2. 解析 JSON 配置
    // 3. 根据参数渲染报表
    // TODO: 实现报表渲染逻辑

    return &RenderReportResponse{
        HTML: "<!-- rendered HTML -->",
        Data: map[string]interface{}{
            "config": report.JSONStr,
        },
    }, nil
}

// 错误定义
var (
    ErrReportNotFound = errors.New("report not found")
)

// 辅助函数
func generateID() string {
    return uuid.New().String()
}

func ptrInt(i int) *int {
    return &i
}
```

### 3.5 Handler

```go
package report

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type ReportHandler interface {
    // Create 创建报表
    Create(c *gin.Context)

    // Update 更新报表
    Update(c *gin.Context)

    // Delete 删除报表
    Delete(c *gin.Context)

    // Get 获取报表详情
    Get(c *gin.Context)

    // List 查询报表列表
    List(c *gin.Context)

    // Render 渲染报表
    Render(c *gin.Context)
}

type reportHandler struct {
    reportService ReportService
}

func NewReportHandler(reportService ReportService) ReportHandler {
    return &reportHandler{reportService: reportService}
}

func (h *reportHandler) Create(c *gin.Context) {
    var req CreateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    report, err := h.reportService.Create(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, report)
}

func (h *reportHandler) Update(c *gin.Context) {
    var req UpdateReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    report, err := h.reportService.Update(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, report)
}

func (h *reportHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    err := h.reportService.Delete(id, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func (h *reportHandler) Get(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    report, err := h.reportService.Get(id, user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, report)
}

func (h *reportHandler) List(c *gin.Context) {
    var req ListReportRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 设置默认值
    if req.Page == 0 {
        req.Page = 1
    }
    if req.PageSize == 0 {
        req.PageSize = 20
    }

    user := getCurrentUser(c)
    resp, err := h.reportService.List(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *reportHandler) Render(c *gin.Context) {
    var req RenderReportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    resp, err := h.reportService.Render(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

// 辅助函数
func getCurrentUser(c *gin.Context) *User {
    if user, exists := c.Get("user"); exists {
        return user.(*User)
    }
    return nil
}
```

## 4. Dashboard 模块

### 4.1 职责

- 仪表盘页面 CRUD
- 组件管理
- 数据集管理
- 仪表盘预览

### 4.2 Models

```go
package dashboard

// Page 仪表盘页面
type Page struct {
    ID             string `json:"id" gorm:"primaryKey"`
    Name           string `json:"name" gorm:"column:name"`
    Path           string `json:"path" gorm:"column:path"`
    BackgroundColor *string `json:"background_color" gorm:"column:background_color"`
    BackgroundImage *string `json:"background_image" gorm:"column:background_image"`
    DesignType     int    `json:"design_type" gorm:"column:design_type"`
    Theme          *string `json:"theme" gorm:"column:theme"`
    Style          *string `json:"style" gorm:"column:style"`
    CoverURL       *string `json:"cover_url" gorm:"column:cover_url"`
    DesJSON        string `json:"des_json" gorm:"column:des_json"` // 配置 JSON
    Template       string `json:"template" gorm:"column:template;type:longtext"` // 布局 JSON
    ProtectionCode *string `json:"protection_code" gorm:"column:protection_code"`
    Type           *string `json:"type" gorm:"column:type"`
    IzTemplate     string `json:"iz_template" gorm:"column:iz_template"`
    CreateBy       string `json:"create_by" gorm:"column:create_by"`
    CreateTime     time.Time `json:"create_time" gorm:"column:create_time"`
    UpdateBy       string `json:"update_by" gorm:"column:update_by"`
    UpdateTime     time.Time `json:"update_time" gorm:"column:update_time"`
    TenantID       *int    `json:"tenant_id" gorm:"column:tenant_id"`
    VisitsNum     *int    `json:"visits_num" gorm:"column:visits_num"`
    DelFlag        *int    `json:"del_flag" gorm:"column:del_flag"`
}

// DatasetHead 数据集头
type DatasetHead struct {
    ID        string  `json:"id" gorm:"primaryKey"`
    Name      string  `json:"name" gorm:"column:name"`
    Code      *string `json:"code" gorm:"column:code"`
    ParentID  *string `json:"parent_id" gorm:"column:parent_id"`
    DBSource  *string `json:"db_source" gorm:"column:db_source"`
    QuerySQL  string  `json:"query_sql" gorm:"column:query_sql"`
    Content   *string `json:"content" gorm:"column:content"`
    IzAgent   string  `json:"iz_agent" gorm:"column:iz_agent"`
    DataType  *string `json:"data_type" gorm:"column:data_type"`
    APIMethod *string `json:"api_method" gorm:"column:api_method"`
    CreateBy  string  `json:"create_by" gorm:"column:create_by"`
    TenantID  *int    `json:"tenant_id" gorm:"column:tenant_id"`
}

// PageComponent 页面组件
type PageComponent struct {
    ID        string `json:"id" gorm:"primaryKey"`
    PageID    string `json:"page_id" gorm:"column:page_id"`
    Component string `json:"component" gorm:"column:component"` // JSON 配置
    Config    string `json:"config" gorm:"column:config;type:longtext"`
    CreateBy  string `json:"create_by" gorm:"column:create_by"`
}

// CreatePageRequest 创建页面请求
type CreatePageRequest struct {
    Name           string  `json:"name" binding:"required"`
    DesignType     int     `json:"design_type" binding:"required"`
    DesJSON        string  `json:"des_json"`
    Template       string  `json:"template"`
    ProtectionCode *string `json:"protection_code"`
}

// UpdatePageRequest 更新页面请求
type UpdatePageRequest struct {
    ID             string  `json:"id" binding:"required"`
    Name           *string `json:"name"`
    BackgroundColor *string `json:"background_color"`
    BackgroundImage *string `json:"background_image"`
    Template       *string `json:"template"`
}

// ListPageRequest 查询页面列表请求
type ListPageRequest struct {
    Page     int    `json:"page" binding:"min=1"`
    PageSize int    `json:"page_size" binding:"min=1,max=100"`
    Keyword  *string `json:"keyword"`
    Type     *string `json:"type"`
}
```

### 4.3 Repository

```go
package dashboard

import "gorm.io/gorm"

type PageRepository interface {
    Create(page *Page) error
    Update(page *Page) error
    Delete(id string) error
    FindByID(id string) (*Page, error)
    FindByPath(path string) (*Page, error)
    List(params *ListPageRequest, tenantID *string) ([]*Page, int64, error)
    IncrementVisits(id string) error
}

type DatasetRepository interface {
    Create(dataset *DatasetHead) error
    Update(dataset *DatasetHead) error
    Delete(id string) error
    FindByID(id string) (*DatasetHead, error)
    List(tenantID *int) ([]*DatasetHead, error)
}

type pageRepository struct {
    db *gorm.DB
}

func NewPageRepository(db *gorm.DB) PageRepository {
    return &pageRepository{db: db}
}

func (r *pageRepository) Create(page *Page) error {
    return r.db.Create(page).Error
}

func (r *pageRepository) Update(page *Page) error {
    return r.db.Save(page).Error
}

func (r *pageRepository) Delete(id string) error {
    return r.db.Model(&Page{}).Where("id = ?", id).Update("del_flag", 1).Error
}

func (r *pageRepository) FindByID(id string) (*Page, error) {
    var page Page
    err := r.db.Where("id = ? AND del_flag = 0", id).First(&page).Error
    if err != nil {
        return nil, err
    }
    return &page, nil
}

func (r *pageRepository) FindByPath(path string) (*Page, error) {
    var page Page
    err := r.db.Where("path = ? AND del_flag = 0", path).First(&page).Error
    if err != nil {
        return nil, err
    }
    return &page, nil
}

func (r *pageRepository) List(params *ListPageRequest, tenantID *string) ([]*Page, int64, error) {
    var pages []*Page
    var total int64

    query := r.db.Model(&Page{}).Where("del_flag = ?", 0)

    if params.Keyword != nil && *params.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+*params.Keyword+"%")
    }

    if params.Type != nil {
        query = query.Where("type = ?", *params.Type)
    }

    if tenantID != nil {
        query = query.Where("tenant_id = ?", *tenantID)
    }

    query.Count(&total)

    offset := (params.Page - 1) * params.PageSize
    err := query.Offset(offset).Limit(params.PageSize).Find(&pages).Error

    return pages, total, err
}

func (r *pageRepository) IncrementVisits(id string) error {
    return r.db.Model(&Page{}).Where("id = ?", id).UpdateColumn("visits_num", gorm.Expr("visits_num + ?", 1)).Error
}
```

### 4.4 Service

```go
package dashboard

type DashboardService interface {
    CreatePage(req *CreatePageRequest, user *User) (*Page, error)
    UpdatePage(req *UpdatePageRequest, user *User) (*Page, error)
    DeletePage(id string, user *User) error
    GetPage(id string, user *User) (*Page, error)
    ListPages(params *ListPageRequest, user *User) (*PageListResponse, error)
    ViewPage(path string, user *User) (*Page, error)
}

type PageListResponse struct {
    List     []*Page `json:"list"`
    Total    int64   `json:"total"`
    Page     int     `json:"page"`
    PageSize int     `json:"page_size"`
}

type dashboardService struct {
    pageRepo PageRepository
}

func NewDashboardService(pageRepo PageRepository) DashboardService {
    return &dashboardService{pageRepo: pageRepo}
}

func (s *dashboardService) CreatePage(req *CreatePageRequest, user *User) (*Page, error) {
    page := &Page{
        ID:         generateID(),
        Name:       req.Name,
        DesignType: req.DesignType,
        DesJSON:    req.DesJSON,
        Template:    req.Template,
        ProtectionCode: req.ProtectionCode,
        IzTemplate: "0",
        CreateBy:   user.ID,
        CreateTime: time.Now(),
        DelFlag:    ptrInt(0),
        VisitsNum:  ptrInt(0),
    }

    // 租户 ID 处理
    if user.TenantID != nil {
        tenantID, _ := strconv.Atoi(*user.TenantID)
        page.TenantID = &tenantID
    }

    if err := s.pageRepo.Create(page); err != nil {
        return nil, err
    }

    return page, nil
}

func (s *dashboardService) UpdatePage(req *UpdatePageRequest, user *User) (*Page, error) {
    page, err := s.pageRepo.FindByID(req.ID)
    if err != nil {
        return nil, ErrPageNotFound
    }

    // 更新字段
    if req.Name != nil {
        page.Name = *req.Name
    }
    if req.Template != nil {
        page.Template = *req.Template
    }
    page.UpdateBy = user.ID
    page.UpdateTime = time.Now()

    if err := s.pageRepo.Update(page); err != nil {
        return nil, err
    }

    return page, nil
}

func (s *dashboardService) DeletePage(id string, user *User) error {
    page, err := s.pageRepo.FindByID(id)
    if err != nil {
        return ErrPageNotFound
    }

    // 检查权限（TODO）

    return s.pageRepo.Delete(id)
}

func (s *dashboardService) GetPage(id string, user *User) (*Page, error) {
    page, err := s.pageRepo.FindByID(id)
    if err != nil {
        return nil, ErrPageNotFound
    }

    return page, nil
}

func (s *dashboardService) ViewPage(path string, user *User) (*Page, error) {
    page, err := s.pageRepo.FindByPath(path)
    if err != nil {
        return nil, ErrPageNotFound
    }

    // 增加访问次数
    s.pageRepo.IncrementVisits(page.ID)

    return page, nil
}

func (s *dashboardService) ListPages(params *ListPageRequest, user *User) (*PageListResponse, error) {
    pages, total, err := s.pageRepo.List(params, user.TenantID)
    if err != nil {
        return nil, err
    }

    return &PageListResponse{
        List:     pages,
        Total:     total,
        Page:      params.Page,
        PageSize:  params.PageSize,
    }, nil
}

var (
    ErrPageNotFound = errors.New("page not found")
)
```

### 4.5 Handler

```go
package dashboard

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type DashboardHandler interface {
    CreatePage(c *gin.Context)
    UpdatePage(c *gin.Context)
    DeletePage(c *gin.Context)
    GetPage(c *gin.Context)
    ListPages(c *gin.Context)
    ViewPage(c *gin.Context)
}

type dashboardHandler struct {
    dashboardService DashboardService
}

func NewDashboardHandler(dashboardService DashboardService) DashboardHandler {
    return &dashboardHandler{dashboardService: dashboardService}
}

func (h *dashboardHandler) CreatePage(c *gin.Context) {
    var req CreatePageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    page, err := h.dashboardService.CreatePage(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, page)
}

func (h *dashboardHandler) UpdatePage(c *gin.Context) {
    var req UpdatePageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    page, err := h.dashboardService.UpdatePage(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, page)
}

func (h *dashboardHandler) DeletePage(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    err := h.dashboardService.DeletePage(id, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func (h *dashboardHandler) GetPage(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    page, err := h.dashboardService.GetPage(id, user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, page)
}

func (h *dashboardHandler) ListPages(c *gin.Context) {
    var req ListPageRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 设置默认值
    if req.Page == 0 {
        req.Page = 1
    }
    if req.PageSize == 0 {
        req.PageSize = 20
    }

    user := getCurrentUser(c)
    resp, err := h.dashboardService.ListPages(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *dashboardHandler) ViewPage(c *gin.Context) {
    path := c.Param("path")
    if path == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
        return
    }

    user := getCurrentUser(c)
    page, err := h.dashboardService.ViewPage(path, user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, page)
}
```

## 5. DataSource 模块

### 5.1 职责

- 数据源 CRUD
- 数据源连接测试
- 数据库连接池管理

### 5.2 Models

```go
package datasource

// DataSource 数据源模型
type DataSource struct {
    ID          string `json:"id" gorm:"primaryKey"`
    Name        string `json:"name" gorm:"column:name"`
    ReportID    *string `json:"report_id" gorm:"column:report_id"`
    Code        *string `json:"code" gorm:"column:code"`
    Remark      *string `json:"remark" gorm:"column:remark"`
    DBType      string `json:"db_type" gorm:"column:db_type"`
    DBDriver    string `json:"db_driver" gorm:"column:db_driver"`
    DBURL       string `json:"db_url" gorm:"column:db_url"`
    DBUsername  string `json:"db_username" gorm:"column:db_username"`
    DBPassword  string `json:"db_password" gorm:"column:db_password"` // 加密存储
    ConnectTimes int    `json:"connect_times" gorm:"column:connect_times"`
    TenantID    *string `json:"tenant_id" gorm:"column:tenant_id"`
    Type        string `json:"type" gorm:"column:type"` // report or drag
    CreateBy    string `json:"create_by" gorm:"column:create_by"`
    CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
    UpdateBy    string `json:"update_by" gorm:"column:update_by"`
    UpdateTime  time.Time `json:"update_time" gorm:"column:update_time"`
}

// CreateDataSourceRequest 创建数据源请求
type CreateDataSourceRequest struct {
    Name      string `json:"name" binding:"required"`
    ReportID  *string `json:"report_id"`
    Type      string `json:"type" binding:"required,oneof=report drag"`
    DBType    string `json:"db_type" binding:"required"`
    DBURL     string `json:"db_url" binding:"required"`
    DBUsername string `json:"db_username" binding:"required"`
    DBPassword string `json:"db_password" binding:"required"`
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
    DBType    string `json:"db_type" binding:"required"`
    DBURL     string `json:"db_url" binding:"required"`
    DBUsername string `json:"db_username" binding:"required"`
    DBPassword string `json:"db_password" binding:"required"`
}

// TestConnectionResponse 测试连接响应
type TestConnectionResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}
```

### 5.3 Repository

```go
package datasource

import "gorm.io/gorm"

type DataSourceRepository interface {
    Create(ds *DataSource) error
    Update(ds *DataSource) error
    Delete(id string) error
    FindByID(id string) (*DataSource, error)
    List(tenantID *string) ([]*DataSource, error)
    IncrementConnectTimes(id string) error
}

type dataSourceRepository struct {
    db *gorm.DB
}

func NewDataSourceRepository(db *gorm.DB) DataSourceRepository {
    return &dataSourceRepository{db: db}
}

func (r *dataSourceRepository) Create(ds *DataSource) error {
    return r.db.Create(ds).Error
}

func (r *dataSourceRepository) Update(ds *DataSource) error {
    return r.db.Save(ds).Error
}

func (r *dataSourceRepository) Delete(id string) error {
    return r.db.Delete(&DataSource{}, "id = ?", id).Error
}

func (r *dataSourceRepository) FindByID(id string) (*DataSource, error) {
    var ds DataSource
    err := r.db.Where("id = ?", id).First(&ds).Error
    if err != nil {
        return nil, err
    }
    return &ds, nil
}

func (r *dataSourceRepository) List(tenantID *string) ([]*DataSource, error) {
    var dss []*DataSource
    query := r.db.Model(&DataSource{})
    if tenantID != nil {
        query = query.Where("tenant_id = ?", *tenantID)
    }
    err := query.Find(&dss).Error
    return dss, err
}

func (r *dataSourceRepository) IncrementConnectTimes(id string) error {
    return r.db.Model(&DataSource{}).Where("id = ?", id).UpdateColumn("connect_times", gorm.Expr("connect_times + ?", 1)).Error
}
```

### 5.4 Service

```go
package datasource

type DataSourceService interface {
    Create(req *CreateDataSourceRequest, user *User) (*DataSource, error)
    Update(id string, req *CreateDataSourceRequest, user *User) (*DataSource, error)
    Delete(id string, user *User) error
    Get(id string, user *User) (*DataSource, error)
    List(user *User) ([]*DataSource, error)
    TestConnection(req *TestConnectionRequest) (*TestConnectionResponse, error)
    TestConnectionByID(id string, user *User) (*TestConnectionResponse, error)
}

type dataSourceService struct {
    repo DataSourceRepository
}

func NewDataSourceService(repo DataSourceRepository) DataSourceService {
    return &dataSourceService{repo: repo}
}

func (s *dataSourceService) Create(req *CreateDataSourceRequest, user *User) (*DataSource, error) {
    ds := &DataSource{
        ID:         generateID(),
        Name:       req.Name,
        ReportID:   req.ReportID,
        Type:       req.Type,
        DBType:     req.DBType,
        DBDriver:   getDBDriver(req.DBType),
        DBURL:      req.DBURL,
        DBUsername: req.DBUsername,
        DBPassword: encryptPassword(req.DBPassword),
        CreateBy:   user.ID,
        CreateTime: time.Now(),
        TenantID:   user.TenantID,
    }

    if err := s.repo.Create(ds); err != nil {
        return nil, err
    }

    return ds, nil
}

func (s *dataSourceService) Update(id string, req *CreateDataSourceRequest, user *User) (*DataSource, error) {
    ds, err := s.repo.FindByID(id)
    if err != nil {
        return nil, ErrDataSourceNotFound
    }

    ds.Name = req.Name
    ds.DBType = req.DBType
    ds.DBURL = req.DBURL
    ds.DBUsername = req.DBUsername
    if req.DBPassword != "" {
        ds.DBPassword = encryptPassword(req.DBPassword)
    }
    ds.UpdateBy = user.ID
    ds.UpdateTime = time.Now()

    if err := s.repo.Update(ds); err != nil {
        return nil, err
    }

    return ds, nil
}

func (s *dataSourceService) Delete(id string, user *User) error {
    ds, err := s.repo.FindByID(id)
    if err != nil {
        return ErrDataSourceNotFound
    }

    return s.repo.Delete(id)
}

func (s *dataSourceService) Get(id string, user *User) (*DataSource, error) {
    ds, err := s.repo.FindByID(id)
    if err != nil {
        return nil, ErrDataSourceNotFound
    }

    // 移除密码信息
    ds.DBPassword = ""
    return ds, nil
}

func (s *dataSourceService) List(user *User) ([]*DataSource, error) {
    dss, err := s.repo.List(user.TenantID)
    if err != nil {
        return nil, err
    }

    // 移除所有密码信息
    for _, ds := range dss {
        ds.DBPassword = ""
    }

    return dss, nil
}

func (s *dataSourceService) TestConnection(req *TestConnectionRequest) (*TestConnectionResponse, error) {
    driver := getDBDriver(req.DBType)

    // 测试数据库连接
    db, err := sql.Open(driver, req.DBURL)
    if err != nil {
        s.logTestFailure(req)
        return &TestConnectionResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        s.logTestFailure(req)
        return &TestConnectionResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }

    return &TestConnectionResponse{
        Success: true,
        Message: "Connection successful",
    }, nil
}

func (s *dataSourceService) TestConnectionByID(id string, user *User) (*TestConnectionResponse, error) {
    ds, err := s.repo.FindByID(id)
    if err != nil {
        return nil, ErrDataSourceNotFound
    }

    // 解密密码
    password := decryptPassword(ds.DBPassword)

    req := &TestConnectionRequest{
        DBType:    ds.DBType,
        DBURL:     ds.DBURL,
        DBUsername: ds.DBUsername,
        DBPassword: password,
    }

    return s.TestConnection(req)
}

func (s *dataSourceService) logTestFailure(req *TestConnectionRequest) {
    // TODO: 记录连接失败日志
}

func getDBDriver(dbType string) string {
    drivers := map[string]string{
        "MYSQL5.7":    "mysql",
        "MYSQL5.5":    "mysql",
        "ORACLE":       "oracle",
        "SQLSERVER":    "sqlserver",
        "POSTGRESQL":   "postgres",
    }
    if driver, ok := drivers[dbType]; ok {
        return driver
    }
    return ""
}

func encryptPassword(password string) string {
    // TODO: 实现密码加密
    return password
}

func decryptPassword(encrypted string) string {
    // TODO: 实现密码解密
    return encrypted
}

var (
    ErrDataSourceNotFound = errors.New("datasource not found")
)
```

### 5.5 Handler

```go
package datasource

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type DataSourceHandler interface {
    Create(c *gin.Context)
    Update(c *gin.Context)
    Delete(c *gin.Context)
    Get(c *gin.Context)
    List(c *gin.Context)
    TestConnection(c *gin.Context)
    TestConnectionByID(c *gin.Context)
}

type dataSourceHandler struct {
    service DataSourceService
}

func NewDataSourceHandler(service DataSourceService) DataSourceHandler {
    return &dataSourceHandler{service: service}
}

func (h *dataSourceHandler) Create(c *gin.Context) {
    var req CreateDataSourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    ds, err := h.service.Create(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, ds)
}

func (h *dataSourceHandler) Update(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    var req CreateDataSourceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    ds, err := h.service.Update(id, &req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ds)
}

func (h *dataSourceHandler) Delete(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    err := h.service.Delete(id, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func (h *dataSourceHandler) Get(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    ds, err := h.service.Get(id, user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ds)
}

func (h *dataSourceHandler) List(c *gin.Context) {
    user := getCurrentUser(c)
    dss, err := h.service.List(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"list": dss})
}

func (h *dataSourceHandler) TestConnection(c *gin.Context) {
    var req TestConnectionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    resp, err := h.service.TestConnection(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}

func (h *dataSourceHandler) TestConnectionByID(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    user := getCurrentUser(c)
    resp, err := h.service.TestConnectionByID(id, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, resp)
}
```

## 6. Export 模块

### 6.1 职责

- 报表导出（Excel, PDF, Word, Image）
- 导出任务管理
- 导出日志记录

### 6.2 Models

```go
package export

// ExportJob 导出任务
type ExportJob struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    ReportID    string    `json:"report_id" gorm:"column:report_id"`
    ExportType  string    `json:"export_type" gorm:"column:export_type"`
    Params      string    `json:"params" gorm:"column:params;type:longtext"`
    Status      string    `json:"status" gorm:"column:status"`
    FileURL     *string   `json:"file_url" gorm:"column:file_url"`
    ErrorMsg    *string   `json:"error_msg" gorm:"column:error_msg;type:text"`
    CreateBy    string    `json:"create_by" gorm:"column:create_by"`
    CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
    UpdateTime  time.Time `json:"update_time" gorm:"column:update_time"`
    TenantID    *string   `json:"tenant_id" gorm:"column:tenant_id"`
}

// ExportLog 导出日志
type ExportLog struct {
    ID         string `json:"id" gorm:"primaryKey"`
    ReportID   string `json:"report_id" gorm:"column:report_id"`
    ExportType string `json:"export_type" gorm:"column:export_type"`
    FileSize   int64  `json:"file_size" gorm:"column:file_size"`
    Duration   int64  `json:"duration" gorm:"column:duration"`
    Status     string `json:"status" gorm:"column:status"`
    CreateBy   string `json:"create_by" gorm:"column:create_by"`
    CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

// ExportRequest 导出请求
type ExportRequest struct {
    ReportID   string                 `json:"report_id" binding:"required"`
    ExportType string                 `json:"export_type" binding:"required,oneof=excel pdf word image"`
    Params     map[string]interface{} `json:"params"`
}

// ExportResponse 导出响应
type ExportResponse struct {
    FileURL string `json:"file_url,omitempty"`
    File    []byte `json:"file,omitempty"`
}
```

### 6.3 Service

```go
package export

type ExportService interface {
    Export(req *ExportRequest, user *User) (*ExportResponse, error)
    CreateJob(job *ExportJob) error
    GetJob(id string, user *User) (*ExportJob, error)
    ListJobs(user *User) ([]*ExportJob, error)
}

type exportService struct {
    reportRepo report.ReportRepository
}

func NewExportService(reportRepo report.ReportRepository) ExportService {
    return &exportService{reportRepo: reportRepo}
}

func (s *exportService) Export(req *ExportRequest, user *User) (*ExportResponse, error) {
    // 1. 查找报表
    report, err := s.reportRepo.FindByID(req.ReportID)
    if err != nil {
        return nil, ErrReportNotFound
    }

    // 2. 根据导出类型生成文件
    var file []byte
    var err error

    switch req.ExportType {
    case "excel":
        file, err = s.exportToExcel(report, req.Params)
    case "pdf":
        file, err = s.exportToPDF(report, req.Params)
    case "word":
        file, err = s.exportToWord(report, req.Params)
    case "image":
        file, err = s.exportToImage(report, req.Params)
    default:
        return nil, ErrUnsupportedExportType
    }

    if err != nil {
        return nil, err
    }

    // 3. 记录导出日志
    // TODO: 实现日志记录

    return &ExportResponse{
        File: file,
    }, nil
}

func (s *exportService) exportToExcel(report *report.Report, params map[string]interface{}) ([]byte, error) {
    // TODO: 实现 Excel 导出
    return nil, errors.New("not implemented")
}

func (s *exportService) exportToPDF(report *report.Report, params map[string]interface{}) ([]byte, error) {
    // TODO: 实现 PDF 导出
    return nil, errors.New("not implemented")
}

func (s *exportService) exportToWord(report *report.Report, params map[string]interface{}) ([]byte, error) {
    // TODO: 实现 Word 导出
    return nil, errors.New("not implemented")
}

func (s *exportService) exportToImage(report *report.Report, params map[string]interface{}) ([]byte, error) {
    // TODO: 实现图片导出
    return nil, errors.New("not implemented")
}

var (
    ErrUnsupportedExportType = errors.New("unsupported export type")
)
```

### 6.4 Handler

```go
package export

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type ExportHandler interface {
    Export(c *gin.Context)
}

type exportHandler struct {
    service ExportService
}

func NewExportHandler(service ExportService) ExportHandler {
    return &exportHandler{service: service}
}

func (h *exportHandler) Export(c *gin.Context) {
    var req ExportRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user := getCurrentUser(c)
    resp, err := h.service.Export(&req, user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.Header("Content-Type", getContentType(req.ExportType))
    c.Header("Content-Disposition", "attachment; filename=export."+getFileExtension(req.ExportType))
    c.Data(http.StatusOK, getContentType(req.ExportType), resp.File)
}

func getContentType(exportType string) string {
    types := map[string]string{
        "excel": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        "pdf":   "application/pdf",
        "word":  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        "image": "image/png",
    }
    return types[exportType]
}

func getFileExtension(exportType string) string {
    exts := map[string]string{
        "excel": "xlsx",
        "pdf":   "pdf",
        "word":  "docx",
        "image": "png",
    }
    return exts[exportType]
}
```

## 7. 共享工具函数

```go
package utils

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    auth "jimureport-go/internal/auth"
)

// GenerateID 生成 UUID
func GenerateID() string {
    return uuid.New().String()
}

// GetCurrentUser 从上下文获取当前用户
func GetCurrentUser(c *gin.Context) *auth.User {
    if user, exists := c.Get("user"); exists {
        return user.(*auth.User)
    }
    return nil
}

// GetTenantID 从上下文获取租户 ID
func GetTenantID(c *gin.Context) *string {
    if tenantID, exists := c.Get("tenant_id"); exists {
        return tenantID.(*string)
    }
    return nil
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    data,
    })
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, error string, message string) {
    c.JSON(statusCode, gin.H{
        "success": false,
        "error":   error,
        "message": message,
    })
}
```

## 8. 接口契约总结

| 模块 | HTTP 路由前缀 | 主要功能 |
|------|---------------|---------|
| Auth | /api/v1/auth | 登录、刷新 token、登出 |
| Report | /api/v1/jmreport | 报表 CRUD、渲染、预览 |
| Dashboard | /api/v1/drag | 仪表盘页面管理、组件管理、数据集管理 |
| DataSource | /api/v1/datasource | 数据源 CRUD、连接测试 |
| Export | /api/v1/export | 报表导出（Excel/PDF/Word/Image） |

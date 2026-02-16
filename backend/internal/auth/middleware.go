package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey   contextKey = "userId"
	UsernameKey contextKey = "username"
	TenantIDKey contextKey = "tenantId"
	RolesKey    contextKey = "roles"
)

var publicPaths = []string{
	"/health",
	"/api/v1/auth/login",
	"/jmreport/list",
	"/drag/list",
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, publicPath := range publicPaths {
			if strings.HasPrefix(path, publicPath) {
				c.Next()
				return
			}
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			authHeader = c.Query("token")
		}

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "missing authorization token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if IsTokenRevoked(c.Request.Context(), tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "token revoked",
			})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(UsernameKey), claims.Username)
		c.Set(string(TenantIDKey), claims.TenantID)
		c.Set(string(RolesKey), claims.Roles)

		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	if userID, exists := c.Get(string(UserIDKey)); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(string(UsernameKey)); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

func GetTenantID(c *gin.Context) string {
	if tenantID, exists := c.Get(string(TenantIDKey)); exists {
		if tid, ok := tenantID.(string); ok {
			return tid
		}
	}
	return ""
}

func GetRoles(c *gin.Context) []string {
	if roles, exists := c.Get(string(RolesKey)); exists {
		if r, ok := roles.([]string); ok {
			return r
		}
	}
	return nil
}

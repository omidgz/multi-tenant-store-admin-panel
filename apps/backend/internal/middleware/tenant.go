package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/config"
)

type TenantContextKey string

const TenantIDKey TenantContextKey = "tenant_id"

func TenantMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID string

		if cfg.TenantMode == "header" {
			tenantID = c.GetHeader("X-Tenant-ID")
			if tenantID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "X-Tenant-ID header is required"})
				c.Abort()
				return
			}
		} else if cfg.TenantMode == "jwt" {
			// TODO: Implement JWT parsing and validation
			// For now, stub
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT mode not implemented yet"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid TENANT_MODE"})
			c.Abort()
			return
		}

		// Add tenant_id to context
		ctx := context.WithValue(c.Request.Context(), TenantIDKey, tenantID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// Helper to get tenant ID from context
func GetTenantID(c *gin.Context) string {
	if tenantID, ok := c.Request.Context().Value(TenantIDKey).(string); ok {
		return tenantID
	}
	return ""
}

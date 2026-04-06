package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/config"
)

func TestTenantMiddleware_HeaderMode_Success(t *testing.T) {
	cfg := &config.Config{TenantMode: "header"}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TenantMiddleware(cfg))

	router.GET("/test", func(c *gin.Context) {
		tenantID := GetTenantID(c)
		c.JSON(http.StatusOK, gin.H{"tenant_id": tenantID})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Tenant-ID", "tenant-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"tenant_id":"tenant-123"`)
}

func TestTenantMiddleware_HeaderMode_MissingHeader(t *testing.T) {
	cfg := &config.Config{TenantMode: "header"}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TenantMiddleware(cfg))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "X-Tenant-ID header is required")
}

func TestTenantMiddleware_JWTMode_Stub(t *testing.T) {
	cfg := &config.Config{TenantMode: "jwt"}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TenantMiddleware(cfg))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "JWT mode not implemented yet")
}

func TestTenantMiddleware_InvalidMode(t *testing.T) {
	cfg := &config.Config{TenantMode: "invalid"}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TenantMiddleware(cfg))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req, _ := http.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid TENANT_MODE")
}
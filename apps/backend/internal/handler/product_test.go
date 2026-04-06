package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/config"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/middleware"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/model"
)

// MockProductService is a mock implementation of ProductServiceInterface
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(req *model.CreateProductRequest, tenantID string) (*model.Product, error) {
	args := m.Called(req, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) GetByID(tenantID, productID string) (*model.Product, error) {
	args := m.Called(tenantID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) ListByTenant(tenantID string) ([]model.Product, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductService) Update(req *model.UpdateProductRequest, tenantID, productID string) (*model.Product, error) {
	args := m.Called(req, tenantID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) Delete(tenantID, productID string) error {
	args := m.Called(tenantID, productID)
	return args.Error(0)
}

func TestProductHandler_CreateProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	reqBody := model.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       19.99,
	}
	expectedProduct := &model.Product{
		ID:          "generated-id",
		TenantID:    "tenant-123",
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Price:       reqBody.Price,
	}

	mockService.On("Create", &reqBody, "tenant-123").Return(expectedProduct, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TenantMiddleware(&config.Config{TenantMode: "header"}))
	router.POST("/products", handler.CreateProduct)

	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tenant-ID", "tenant-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expectedProduct.ID, response.ID)
	mockService.AssertExpectations(t)
}

func TestProductHandler_GetProducts_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []model.Product{
		{ID: "1", TenantID: "tenant-123", Name: "Product 1"},
		{ID: "2", TenantID: "tenant-123", Name: "Product 2"},
	}

	mockService.On("ListByTenant", "tenant-123").Return(products, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.TenantMiddleware(&config.Config{TenantMode: "header"}))
	router.GET("/products", handler.GetProducts)

	req, _ := http.NewRequest("GET", "/products", nil)
	req.Header.Set("X-Tenant-ID", "tenant-123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []model.Product
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 2)
	mockService.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
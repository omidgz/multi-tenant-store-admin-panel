package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/model"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(tenantID, productID string) (*model.Product, error) {
	args := m.Called(tenantID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductRepository) ListByTenant(tenantID string) ([]model.Product, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) Update(product *model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(tenantID, productID string) error {
	args := m.Called(tenantID, productID)
	return args.Error(0)
}

func TestProductService_Create_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := &model.CreateProductRequest{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       19.99,
	}
	tenantID := "tenant-123"

	mockRepo.On("Create", mock.AnythingOfType("*model.Product")).Return(nil)

	product, err := service.Create(req, tenantID)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, tenantID, product.TenantID)
	assert.Equal(t, req.Name, product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Create_ValidationError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	req := &model.CreateProductRequest{
		Name:  "",
		Price: 19.99,
	}
	tenantID := "tenant-123"

	product, err := service.Create(req, tenantID)

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestProductService_GetByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	tenantID := "tenant-123"
	productID := "prod-456"
	expectedProduct := &model.Product{ID: productID, TenantID: tenantID, Name: "Test"}

	mockRepo.On("GetByID", tenantID, productID).Return(expectedProduct, nil)

	product, err := service.GetByID(tenantID, productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Update_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	tenantID := "tenant-123"
	productID := "prod-456"
	req := &model.UpdateProductRequest{
		Name:  stringPtr("Updated Name"),
		Price: floatPtr(29.99),
	}

	existingProduct := &model.Product{ID: productID, TenantID: tenantID, Name: "Old"}

	mockRepo.On("GetByID", tenantID, productID).Return(existingProduct, nil)
	mockRepo.On("Update", mock.AnythingOfType("*model.Product")).Return(nil)

	product, err := service.Update(req, tenantID, productID)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Updated Name", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_Delete_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	service := NewProductService(mockRepo)

	tenantID := "tenant-123"
	productID := "prod-456"

	mockRepo.On("Delete", tenantID, productID).Return(nil)

	err := service.Delete(tenantID, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

package service

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/model"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/repository"
)

type ProductServiceInterface interface {
	Create(req *model.CreateProductRequest, tenantID string) (*model.Product, error)
	GetByID(tenantID, productID string) (*model.Product, error)
	ListByTenant(tenantID string) ([]model.Product, error)
	Update(req *model.UpdateProductRequest, tenantID, productID string) (*model.Product, error)
	Delete(tenantID, productID string) error
}

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(req *model.CreateProductRequest, tenantID string) (*model.Product, error) {
	// Basic validation
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("name is required")
	}
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than 0")
	}

	// Create product model
	product := &model.Product{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Call repo
	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetByID(tenantID, productID string) (*model.Product, error) {
	return s.repo.GetByID(tenantID, productID)
}

func (s *ProductService) ListByTenant(tenantID string) ([]model.Product, error) {
	return s.repo.ListByTenant(tenantID)
}

func (s *ProductService) Update(req *model.UpdateProductRequest, tenantID, productID string) (*model.Product, error) {
	// Get existing
	existing, err := s.repo.GetByID(tenantID, productID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			return nil, errors.New("name cannot be empty")
		}
		existing.Name = strings.TrimSpace(*req.Name)
	}
	if req.Description != nil {
		existing.Description = strings.TrimSpace(*req.Description)
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, errors.New("price must be greater than 0")
		}
		existing.Price = *req.Price
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	existing.UpdatedAt = time.Now()

	// Call repo
	err = s.repo.Update(existing)
	if err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *ProductService) Delete(tenantID, productID string) error {
	return s.repo.Delete(tenantID, productID)
}

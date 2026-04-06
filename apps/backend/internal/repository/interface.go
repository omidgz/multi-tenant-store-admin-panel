package repository

import "github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/model"

type ProductRepository interface {
	Create(product *model.Product) error
	GetByID(tenantID, productID string) (*model.Product, error)
	ListByTenant(tenantID string) ([]model.Product, error)
	Update(product *model.Product) error
	Delete(tenantID, productID string) error
}

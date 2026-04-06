package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/omidgz/multi-tenant-store-admin-panel/apps/backend/internal/model"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(product *model.Product) error {
	query := `
		INSERT INTO products (id, tenant_id, name, description, price, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(context.Background(), query,
		product.ID,
		product.TenantID,
		product.Name,
		product.Description,
		product.Price,
		product.ImageURL,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

func (r *PostgresProductRepository) GetByID(tenantID, productID string) (*model.Product, error) {
	query := `SELECT id, tenant_id, name, description, price, image_url, created_at, updated_at 
	          FROM products 
	          WHERE id = $1 AND tenant_id = $2`

	var p model.Product
	err := r.db.QueryRow(context.Background(), query, productID, tenantID).Scan(
		&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresProductRepository) ListByTenant(tenantID string) ([]model.Product, error) {
	query := `SELECT id, tenant_id, name, description, price, image_url, created_at, updated_at 
	          FROM products 
	          WHERE tenant_id = $1 
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(context.Background(), query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *PostgresProductRepository) Update(product *model.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, image_url = $4, updated_at = $5
		WHERE id = $6 AND tenant_id = $7`

	result, err := r.db.Exec(context.Background(), query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageURL,
		product.UpdatedAt,
		product.ID,
		product.TenantID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found or tenant mismatch")
	}
	return nil
}

func (r *PostgresProductRepository) Delete(tenantID, productID string) error {
	query := `DELETE FROM products WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.Exec(context.Background(), query, productID, tenantID)
	return err
}

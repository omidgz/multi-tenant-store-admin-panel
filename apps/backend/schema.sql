-- Demo schema for multi-tenant store admin panel
-- Run this in PostgreSQL to set up the database

-- Enable RLS
ALTER DATABASE store_admin SET row_security = on;

-- Create products table
CREATE TABLE products (
    id VARCHAR(36) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0),
    image_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on tenant_id for performance
CREATE INDEX idx_products_tenant_id ON products(tenant_id);

-- Enable RLS on products table
ALTER TABLE products ENABLE ROW LEVEL SECURITY;

-- RLS policy: users can only see/modify their own tenant's products
-- Note: In production, this would be enforced by the application and DB roles
CREATE POLICY tenant_isolation_policy ON products
    FOR ALL
    USING (tenant_id = current_setting('app.tenant_id', true));

-- Note: For demo, we'll set the tenant_id via application logic, not session variables
-- In production, use proper DB roles and policies
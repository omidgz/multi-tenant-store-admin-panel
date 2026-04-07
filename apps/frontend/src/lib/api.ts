const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  tenant_id: string;
  created_at: string;
  updated_at: string;
}

export interface CreateProductRequest {
  name: string;
  description: string;
  price: number;
}

export class ApiClient {
  private tenantId: string;

  constructor(tenantId: string) {
    this.tenantId = tenantId;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`;
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        'X-Tenant-ID': this.tenantId,
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API error: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  async getProducts(): Promise<Product[]> {
    const result = await this.request<any>('/api/products');
    return Array.isArray(result) ? result : [];
  }

  async createProduct(product: CreateProductRequest): Promise<Product> {
    return this.request('/api/products', {
      method: 'POST',
      body: JSON.stringify(product),
    });
  }

  async updateProduct(id: string, product: Partial<CreateProductRequest>): Promise<Product> {
    return this.request(`/api/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(product),
    });
  }

  async deleteProduct(id: string): Promise<void> {
    return this.request(`/api/products/${id}`, {
      method: 'DELETE',
    });
  }
}
# Multi-Tenant Store Admin Panel on Amazon EKS

## Project Overview

A **multi-tenant SaaS-style admin control panel** for product/store management built on **Amazon EKS (Kubernetes)** using **Application-level Multi-Tenancy**.

All tenants share the same infrastructure (EKS cluster, backend, and frontend), while **data isolation** is strictly enforced at the application and database layers. Each tenant (store/company) can only access and manage their own products, images, and data.

This design prioritizes **scalability**, **cost-efficiency**, and **operational simplicity** while still demonstrating strong multi-tenancy practices.

**Goal**: Showcase how to build a scalable SaaS application on AWS EKS with proper data isolation using shared resources.

## Multi-Tenancy Strategy

**Chosen Model**: **Application-level Multi-Tenancy** (Shared Everything)

### Isolation Approach
- **Infrastructure**: Fully shared (single EKS cluster, single Go backend deployment, single React frontend)
- **Authentication**: AWS Cognito with custom `tenant_id` claim in JWT tokens
- **Application Layer**: Go backend extracts `tenant_id` from JWT and enforces it on every request and database query
- **Database Layer**: PostgreSQL with **Row Level Security (RLS)** policies
- **Storage Layer**: Amazon S3 with tenant-prefixed object keys + strict IAM policies

This model is highly scalable and is commonly used by large SaaS platforms that need to support tens or hundreds of thousands of tenants efficiently.

## High-Level Architecture

``` mermaid
flowchart TD
    Users[Users<br/>Tenant Admins] 
    -->|HTTPS + JWT Token| ALB[AWS ALB + Ingress Controller]

    ALB --> EKS[Amazon EKS Cluster]

    subgraph EKS ["Amazon EKS Cluster"]
        direction TB

        subgraph Shared ["Shared Infrastructure"]
            Backend[Go Backend API - Single Deployment]
            Frontend[React Admin Panel - Single Deployment]
            Cognito[AWS Cognito User Pool]
            Postgres[PostgreSQL with Row Level Security]
            S3[Amazon S3 Bucket with tenant prefixes]
            Obs[Observability Stack]
        end

        TenantFilter[Application-level Tenant Filtering]
    end

    Backend --> TenantFilter
    TenantFilter --> Postgres
    TenantFilter --> S3

    style EKS fill:#e0f2fe,stroke:#0369a1,stroke-width:2px
    style Shared fill:#d1fae5,stroke:#10b981
    style TenantFilter fill:#fef3c7,stroke:#d97706
```

**Figure 1**: High-level architecture using Application-level Multi-Tenancy on a shared Amazon EKS cluster.

## Tech Stack

**Infrastructure**: Amazon EKS, Terraform, Karpenter, ALB Ingress Controller  
**Backend**: Go + Gin or Fiber  
**Frontend**: React + Vite + TypeScript + Tailwind CSS  
**Auth**: AWS Cognito (JWT with `tenant_id` claim)  
**Database**: PostgreSQL with Row Level Security (RLS)  
**Storage**: Amazon S3 with tenant-prefixed keys  

## Core Features

- Tenant-aware authentication via Cognito
- Product CRUD operations fully scoped to tenant
- Secure image upload to S3 with tenant isolation
- Dashboard showing only the tenant's own data
- Strict data separation enforced at both application and database layers

## API Endpoints (Examples)

| Method | Endpoint               | Description                          | Protected |
|--------|------------------------|--------------------------------------|---------|
| POST   | `/api/products`        | Create new product (tenant-scoped)   | Yes     |
| GET    | `/api/products`        | List products (only own tenant)      | Yes     |
| GET    | `/api/products/{id}`   | Get single product                   | Yes     |
| PUT    | `/api/products/{id}`   | Update product                       | Yes     |
| DELETE | `/api/products/{id}`   | Delete product                       | Yes     |

## Scalability & Limitations

- **Scalability**: Excellent — can support 100,000+ tenants efficiently.
- **Advantages**: Low operational overhead, cost-effective, easy horizontal scaling.
- **Trade-offs**: Weaker blast-radius isolation compared to namespace or vCluster models. Requires careful implementation of tenant filtering and RLS.

# Multi-Tenant Store Admin Panel on Amazon EKS

## Project Overview

A **multi-tenant SaaS-style admin control panel** for product and store management built on **Amazon EKS (Kubernetes)** using **Application-level Multi-Tenancy**.

All tenants share the same infrastructure (single EKS cluster, single backend, and single frontend), while **strict data isolation** is enforced at the application and database layers. Each tenant can only view and manage their own products, images, and data.

This design focuses on **scalability**, **cost-efficiency**, and **operational simplicity**.

**Goal**: Demonstrate a modern, scalable SaaS backend on AWS EKS with proper tenant data isolation.

## Multi-Tenancy Strategy

**Chosen Model**: **Application-level Multi-Tenancy** (Shared Everything)

### Isolation Approach
- **Infrastructure**: Fully shared (one EKS cluster, one Go backend deployment, one React frontend)
- **Authentication**: AWS Cognito with custom `tenant_id` claim in JWT tokens
- **Application Layer**: Go backend extracts `tenant_id` from JWT and enforces tenant context on every request
- **Database Layer**: PostgreSQL with **Row Level Security (RLS)** to guarantee data isolation
- **Storage Layer**: Amazon S3 with tenant-prefixed paths (`tenant-{tenant_id}/...`) and strict IAM policies

This approach allows the system to scale efficiently to thousands of tenants without per-tenant infrastructure overhead.

## High-Level Architecture

``` mermaid
flowchart TD
    Users[Users<br/>Tenant Admins] 
    -->|HTTPS + JWT Token| ALB[AWS ALB + Ingress Controller]

    ALB --> EKS[Amazon EKS Cluster]

    subgraph EKS ["Amazon EKS Cluster"]
        direction TB

        subgraph Shared ["Shared Infrastructure"]
            Backend[Go Backend API<br/>- Single Deployment]
            Frontend[React Admin Panel<br/>- Single Deployment]
            Cognito[AWS Cognito User Pool]
            Postgres[PostgreSQL<br/>with Row Level Security]
            S3[Amazon S3 Bucket<br/>with tenant prefixes]
            Obs[Observability Stack]
        end

        TenantFilter[Application-level Tenant Filtering<br/>• Extract tenant_id from JWT<br/>• Filter all queries]
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

**Infrastructure**: 
- Amazon EKS (Kubernetes)
- Terraform (IaC)
- AWS ALB Ingress Controller

**Backend**: 
- Go + Gin/Fiber
- Tenant-aware middleware

**Frontend**: 
- React + Vite + TypeScript + Tailwind CSS

**Authentication**: 
- AWS Cognito (JWT with `tenant_id` claim)

**Data & Storage**: 
- PostgreSQL with Row Level Security (RLS)
- Amazon S3 with tenant-prefixed keys

**Observability**: 
- Prometheus + Grafana + OpenTelemetry

## Core Features

- Secure tenant authentication via Cognito
- Product CRUD operations fully scoped to tenant
- Image upload to S3 with tenant isolation
- Dashboard showing only tenant-specific data
- Strict data separation at both application and database layers

## API Endpoints (Examples)

| Method | Endpoint               | Description                          | Protected |
|--------|------------------------|--------------------------------------|---------|
| POST   | `/api/products`        | Create new product                   | Yes     |
| GET    | `/api/products`        | List tenant's products               | Yes     |
| GET    | `/api/products/{id}`   | Get single product                   | Yes     |
| PUT    | `/api/products/{id}`   | Update product                       | Yes     |
| DELETE | `/api/products/{id}`   | Delete product                       | Yes     |

## Scalability & Limitations

- **Scalability**: Excellent — can support thousands of tenants efficiently due to shared infrastructure.
- **Advantages**: Low cost, simple operations, easy horizontal scaling.
- **Trade-offs**: Weaker infrastructure isolation compared to namespace or vCluster models. Strong application-level and database-level controls are critical.

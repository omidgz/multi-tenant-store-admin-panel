# Multi-Tenant Store Admin Panel

A **multi-tenant SaaS-style admin control panel** for product and store management, built on **Amazon EKS** using **Application-level Multi-Tenancy**.

All tenants share the same infrastructure while maintaining strict data isolation through the application and database layers.

## ✨ Features

- Secure multi-tenant authentication with AWS Cognito
- Tenant-scoped Product CRUD operations
- Image upload with tenant isolation (S3)
- Row Level Security (RLS) at the database layer
- Modern React admin dashboard
- Infrastructure as Code with Terraform

## 🏗️ Architecture

- **Multi-Tenancy Model**: Application-level (Shared Infrastructure)
- **Cluster**: Amazon EKS (Kubernetes)
- **Backend**: Go with tenant middleware
- **Frontend**: React + Vite + TypeScript + Tailwind
- **Database**: PostgreSQL with Row Level Security
- **Storage**: Amazon S3 (tenant-prefixed paths)
- **Auth**: AWS Cognito (JWT with `tenant_id` claim)

See [DESIGN.md](DESIGN.md) for detailed architecture and decisions.

## 🛠️ Tech Stack

- **Infrastructure**: AWS EKS, Terraform
- **Backend**: Go, Gin/Fiber
- **Frontend**: React, Vite, TypeScript, Tailwind CSS
- **Database**: PostgreSQL + RLS
- **Auth**: AWS Cognito
- **IaC**: Terraform

## 📁 Project Structure
```
multi-tenant-store-admin-panel/
├── terraform/              # Infrastructure as Code (EKS + VPC)
├── apps/
│   ├── backend/            # Go API (planned)
│   └── frontend/           # React Admin Panel (planned)
├── k8s/                    # Kubernetes manifests
├── DESIGN.md               # Detailed design document
└── README.md
```


## 🚀 Getting Started

### Prerequisites
- AWS Account with IAM user
- AWS CLI configured
- Terraform installed
- kubectl (for cluster interaction)

### Deploy Infrastructure

```bash
cd terraform
terraform init
terraform plan
terraform apply
```
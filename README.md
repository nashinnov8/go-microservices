# Go E-Commerce Microservices

A production-ready e-commerce microservices platform built with **Go**, **Gin Web Framework**, **sqlc**, and **PostgreSQL (pgx/v5)** using **Clean Architecture** principles.

---

## Architecture Overview

This repository is structured as a **Go Workspaces (`go.work`) Monorepo**, dividing shared domain-agnostic packages from individual domain services.

```text
go-ecommerce/
├── go.work                          # Multi-module Go Workspace
├── Makefile                         # Automation commands
├── deploy/                          # Containerization & orchestration
│   └── docker-compose.yml           # PostgreSQL 16 local setup with auto-migrations
├── pkg/                             # Shared libraries across microservices
│   ├── apperror/                    # Standardized application error definitions
│   ├── database/                    # pgxpool connection & TransactionManager (WithTx)
│   ├── logger/                      # Structured logging (log/slog JSON & text)
│   ├── middleware/                  # Gin middlewares (RequestID, Logger, CORS)
│   └── response/                    # Standardized API JSON response wrappers
└── services/
    ├── product/                     # Product catalog service (Port 8081)
    │   ├── cmd/api/main.go          # Entrypoint with graceful shutdown
    │   ├── migrations/              # SQL schema migrations
    │   ├── sqlc.yaml                # sqlc v2 configuration with pgx/v5
    │   └── internal/                # Clean Architecture layers (domain, dto, service, repo, transport)
    ├── order/                       # Order processing service (Port 8082)
    │   ├── cmd/api/main.go          # Entrypoint with graceful shutdown
    │   ├── migrations/              # Orders and OrderItems schema migrations
    │   ├── sqlc.yaml                # sqlc v2 configuration with pgx/v5
    │   └── internal/                # Clean Architecture layers with Unit of Work (WithTx)
    ├── inventory/                   # Planned inventory service
    ├── payment/                     # Planned payment service
    └── user/                        # Planned authentication and user service
```

---

## Key Design Patterns & Engineering Highlights

1. **Clean / Layered Architecture**:
   - `domain`: Pure business entities and domain error types.
   - `service`: Business logic orchestration without direct dependencies on HTTP frameworks or database drivers.
   - `repository`: Data access implementation powered by type-safe `sqlc` queries and `pgxpool.Pool`.
   - `transport/http`: Gin router, handlers, and middlewares.
2. **Context-based Transaction Manager (`WithTx`)**:
   - Manages atomic multi-table operations (e.g., creating `orders` and `order_items`) cleanly from the `service` layer without leaking `pgx.Tx` into business logic.
   - 100% unit-testable in memory without requiring a live PostgreSQL instance.
3. **Graceful Shutdown**:
   - All services handle OS interrupt signals (`SIGINT`, `SIGTERM`) to cleanly finish in-flight requests and drain database connections before exiting.
4. **Standardized Responses & Error Handling**:
   - Unified API response format (`success`, `data`, `meta`, `error`) and automated HTTP status code mapping.

---

## Getting Started

### Prerequisites

- **Go**: 1.22+ (tested on Go 1.27)
- **Docker & Docker Compose**

### 1. Start Infrastructure (PostgreSQL)

```bash
docker compose -f deploy/docker-compose.yml up -d
```

This starts PostgreSQL on port `5432` and automatically runs initial schema migrations for both `products` and `orders`.

### 2. Run the Services

- **Product Service** (HTTP `:8081`):
  ```bash
  go run -C services/product ./cmd/api
  ```

- **Order Service** (HTTP `:8082`):
  ```bash
  go run -C services/order ./cmd/api
  ```

### 3. Run Unit & Integration Tests

```bash
go test -C services/product -v ./...
go test -C services/order -v ./...
```

---

## API Endpoints

### Product Service (`:8081`)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/healthz` | Health check probe |
| `POST` | `/api/v1/products` | Create a new product |
| `GET` | `/api/v1/products` | List products with pagination & search |
| `GET` | `/api/v1/products/:id` | Get product details by UUID |
| `PUT` | `/api/v1/products/:id` | Update product details |
| `PATCH` | `/api/v1/products/:id/stock` | Atomically adjust inventory stock |
| `DELETE` | `/api/v1/products/:id` | Delete a product |

### Order Service (`:8082`)

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/healthz` | Health check probe |
| `POST` | `/api/v1/orders` | Create an order with items (atomic transaction) |
| `GET` | `/api/v1/orders/:id` | Get order details with items by UUID |

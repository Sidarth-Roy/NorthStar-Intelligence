Backend/
├── cmd/
│   └── api/
│       ├── main.go                 # Entry point & Graceful Shutdown
├── internal/
│   ├── app/
│   │   └── container.go            # Dependency Injection
│   ├── api/
│   │   ├── router.go               # Centralized Route Registry (v1)
│   │   └── v1/
│   │       ├── product_controller.go
│   │       ├── category_controller.go
│   │       ├── customer_controller.go
│   │       ├── employee_controller.go
│   │       ├── shipper_controller.go
│   │       └── order_controller.go      # Handles both Order & OrderDetails
│   ├── service/
│   │   ├── product_service.go
│   │   ├── category_service.go
│   │   ├── customer_service.go
│   │   ├── employee_service.go
│   │   ├── shipper_service.go
│   │   └── order_service.go         # Transactional logic for Orders
│   ├── repository/
│   │   ├── product_repository.go
│   │   ├── category_repository.go
│   │   ├── customer_repository.go
│   │   ├── employee_repository.go
│   │   ├── shipper_repository.go
│   │   └── order_repository.go      # GORM logic with Preloads/Cascades
│   ├── dto/
│   │   ├── product_dto.go
│   │   ├── category_dto.go
│   │   ├── customer_dto.go
│   │   ├── employee_dto.go
│   │   ├── shipper_dto.go
│   │   └── order_dto.go             # Includes Nested OrderDetail DTOs
│   └── middleware/
│       ├── error_handler.go         # Global Exception handling
│       └── request_id.go            # Correlation ID for logging
├── pkg/
│   ├── db/
│   │   └── postgres.go              # GORM Initialization & Connection Pooling
│   ├── logger/
│   │   └── zap.go                   # Structured logging (Uber-Zap)
│   ├── config/
│   │   └── config.go                # Env loader
│   └── model/
│       └── models.go                # All DB Entities & GORM relationships
├── .env                             # Database credentials & App Port
├── go.mod
└── go.sum
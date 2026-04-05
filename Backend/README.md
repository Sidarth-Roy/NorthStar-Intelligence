Backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point & Graceful Shutdown
├── DB/                
│   ├── seeder_main.go              # Your seeding script
│   ├── seeder                      # Your seeding script
│   │   └── seed.go            
│   └── Northwind_Traders_Kaggle_Dataset_CSV/              # Your raw data files
│       └── categories.csv            
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
│   │   ├── product_service_test.go
│   │   ├── category_service_test.go
│   │   ├── customer_service_test.go
│   │   ├── employee_service_test.go
│   │   ├── shipper_service_test.go
│   │   ├── order_service_test.go         
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
│   └── model/                       # All DB Entities & GORM relationships
│       ├── customer.go
│       ├── employee.go
│       ├── model.go
│       ├── order.go
│       ├── product.go
│       └── shipper.go
├── .env                             # Database credentials & App Port
├── go.mod
└── go.sum

Run all tests in the project:
    go test ./...

Run only Service layer tests:
    go test ./internal/service/...

Run with Verbose output (shows pass/fail for each specific test case):
    go test -v ./internal/service/category_service_test.go ./internal/service/category_service.go

Run and check Code Coverage:
    go test -cover ./internal/service/...

Run tests and generate the profile (Ensure the path is correct)
    go test ./internal/service/... -coverprofile=services_coverage.out

Open the HTML representation
    go tool cover -html=services_coverage

Run tests and generate the profile (Ensure the path is correct)
    go test ./internal/repository/... -coverprofile=repository_coverage.out

Open the HTML representation
    go tool cover -html=repository_coverage
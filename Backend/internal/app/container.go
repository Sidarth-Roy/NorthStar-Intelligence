package app

import (
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/api/v1"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"gorm.io/gorm"
)

// AppContainer holds all initialized controllers
type AppContainer struct {
	ProductCtrl  *v1.ProductController
	CategoryCtrl *v1.CategoryController
	CustomerCtrl *v1.CustomerController
	EmployeeCtrl *v1.EmployeeController
	ShipperCtrl  *v1.ShipperController
	OrderCtrl    *v1.OrderController
}

// NewAppContainer initializes all layers of the application
func NewAppContainer(db *gorm.DB) *AppContainer {
	// 1. Initialize Repositories
	prodRepo := repository.NewProductRepo(db)
	catRepo  := repository.NewCategoryRepo(db)
	custRepo := repository.NewCustomerRepo(db)
	empRepo  := repository.NewEmployeeRepo(db)
	shipRepo := repository.NewShipperRepo(db)
	orderRepo := repository.NewOrderRepo(db)

	// 2. Initialize Services
	prodSvc := service.NewProductSvc(prodRepo)
	catSvc  := service.NewCategorySvc(catRepo)
	custSvc := service.NewCustomerSvc(custRepo)
	empSvc  := service.NewEmployeeSvc(empRepo)
	shipSvc := service.NewShipperSvc(shipRepo)
	orderSvc := service.NewOrderSvc(orderRepo)
	
	// 3. Initialize Controllers and return Container
	return &AppContainer{
		ProductCtrl:  v1.NewProductController(prodSvc),
		CategoryCtrl: v1.NewCategoryController(catSvc),
		CustomerCtrl: v1.NewCustomerController(custSvc),
		EmployeeCtrl: v1.NewEmployeeController(empSvc),
		ShipperCtrl:  v1.NewShipperController(shipSvc),
		OrderCtrl:    v1.NewOrderController(orderSvc),
	}
}
package api

import (
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/app"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(deps *app.AppContainer) *gin.Engine {
	r := gin.New()

	// 1. Traceability Middleware
	r.Use(middleware.RequestIDMiddleware())
	
	// Standard Gin Logger for Console
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	// 2. Exception & Logger Middleware
	r.Use(middleware.GlobalExceptionHandler())

	apiV1 := r.Group("/api/v1")
	{
		// Product Routes
		productRoutes := apiV1.Group("/products")
		{
			productRoutes.GET("", deps.ProductCtrl.GetAll)
			productRoutes.GET("/:id", deps.ProductCtrl.GetByID)
			productRoutes.POST("", deps.ProductCtrl.Create)
			productRoutes.PUT("/:id", deps.ProductCtrl.Update)
			productRoutes.DELETE("/:id", deps.ProductCtrl.Delete)
		}
		// Category Routes
		categoryRoutes := apiV1.Group("/categories")
		{
			categoryRoutes.GET("", deps.CategoryCtrl.GetAll)
			categoryRoutes.GET("/:id", deps.CategoryCtrl.GetByID)
			categoryRoutes.POST("", deps.CategoryCtrl.Create)
			categoryRoutes.PUT("/:id", deps.CategoryCtrl.Update)
			categoryRoutes.DELETE("/:id", deps.CategoryCtrl.Delete)
		}
		// Category Routes
		categoryWithProductsRoutes := apiV1.Group("/categorieswithproducts")
		{
			categoryWithProductsRoutes.GET("/:id", deps.CategoryCtrl.GetWithProducts)
		}
		// Customer Routes
		customerRoutes := apiV1.Group("/customers")
		{
			customerRoutes.GET("", deps.CustomerCtrl.GetAll)
			customerRoutes.GET("/:id", deps.CustomerCtrl.GetByID)
			customerRoutes.POST("", deps.CustomerCtrl.Create)
			customerRoutes.PUT("/:id", deps.CustomerCtrl.Update)
			customerRoutes.DELETE("/:id", deps.CustomerCtrl.Delete)
		}
		// Employee Routes
		employeeRoutes := apiV1.Group("/employees")
		{
			employeeRoutes.GET("", deps.EmployeeCtrl.GetAll)
			employeeRoutes.GET("/:id", deps.EmployeeCtrl.GetByID)
			employeeRoutes.POST("", deps.EmployeeCtrl.Create)
			employeeRoutes.PUT("/:id", deps.EmployeeCtrl.Update)
			employeeRoutes.DELETE("/:id", deps.EmployeeCtrl.Delete)
		}
		// Shipper Routes
		shipperRoutes := apiV1.Group("/shippers")
		{
			shipperRoutes.GET("", deps.ShipperCtrl.GetAll)
			shipperRoutes.GET("/:id", deps.ShipperCtrl.GetByID)
			shipperRoutes.POST("", deps.ShipperCtrl.Create)
			shipperRoutes.PUT("/:id", deps.ShipperCtrl.Update)
			shipperRoutes.DELETE("/:id", deps.ShipperCtrl.Delete)
		}
		// Order Routes
		orderRoutes := apiV1.Group("/orders")
		{
			orderRoutes.GET("", deps.OrderCtrl.GetAll)
			orderRoutes.GET("/:id", deps.OrderCtrl.GetByID)
			orderRoutes.POST("", deps.OrderCtrl.Create)
			orderRoutes.PUT("/:id", deps.OrderCtrl.Update)
			orderRoutes.DELETE("/:id", deps.OrderCtrl.Delete)
		}
	}

	return r
}
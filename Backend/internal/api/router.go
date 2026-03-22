package api

import (
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/api/v1"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	prodCtrl *v1.ProductController, 
	catCtrl *v1.CategoryController,
	custCtrl *v1.CustomerController,
	empCtrl *v1.EmployeeController,
	shipCtrl *v1.ShipperController,
	orderCtrl *v1.OrderController,
) *gin.Engine {
	r := gin.New()

	// 1. Traceability Middleware
	r.Use(middleware.RequestIDMiddleware())
	
	// 2. Exception & Logger Middleware
	r.Use(middleware.GlobalExceptionHandler())
	
	// Standard Gin Logger for Console
	r.Use(gin.Logger())

	apiV1 := r.Group("/api/v1")
	{
		// Product Routes
		productRoutes := apiV1.Group("/products")
		{
			productRoutes.GET("", prodCtrl.GetAll)
			productRoutes.GET("/:id", prodCtrl.GetByID)
			productRoutes.POST("", prodCtrl.Create)
			productRoutes.PUT("/:id", prodCtrl.Update)
			productRoutes.DELETE("/:id", prodCtrl.Delete)
		}
		// Category Routes
		categoryRoutes := apiV1.Group("/categories")
		{
			categoryRoutes.GET("", catCtrl.GetAll)
			categoryRoutes.GET("/:id", catCtrl.GetByID)
			categoryRoutes.POST("", catCtrl.Create)
			categoryRoutes.PUT("/:id", catCtrl.Update)
			categoryRoutes.DELETE("/:id", catCtrl.Delete)
		}
		// Customer Routes
		customerRoutes := apiV1.Group("/customers")
		{
			customerRoutes.GET("", custCtrl.GetAll)
			customerRoutes.GET("/:id", custCtrl.GetByID)
			customerRoutes.POST("", custCtrl.Create)
			customerRoutes.PUT("/:id", custCtrl.Update)
			customerRoutes.DELETE("/:id", custCtrl.Delete)
		}
		// Employee Routes
		employeeRoutes := apiV1.Group("/employees")
		{
			employeeRoutes.GET("", empCtrl.GetAll)
			employeeRoutes.GET("/:id", empCtrl.GetByID)
			employeeRoutes.POST("", empCtrl.Create)
			employeeRoutes.PUT("/:id", empCtrl.Update)
			employeeRoutes.DELETE("/:id", empCtrl.Delete)
		}
		// Shipper Routes
		shipperRoutes := apiV1.Group("/shippers")
		{
			shipperRoutes.GET("", shipCtrl.GetAll)
			shipperRoutes.GET("/:id", shipCtrl.GetByID)
			shipperRoutes.POST("", shipCtrl.Create)
			shipperRoutes.PUT("/:id", shipCtrl.Update)
			shipperRoutes.DELETE("/:id", shipCtrl.Delete)
		}
		// Order Routes
		orderRoutes := apiV1.Group("/orders")
		{
			orderRoutes.GET("", orderCtrl.GetAll)
			orderRoutes.GET("/:id", orderCtrl.GetByID)
			orderRoutes.POST("", orderCtrl.Create)
			orderRoutes.PUT("/:id", orderCtrl.Update)
			orderRoutes.DELETE("/:id", orderCtrl.Delete)
		}
	}

	return r
}
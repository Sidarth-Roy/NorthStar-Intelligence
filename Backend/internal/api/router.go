package api

import (
	"net/http"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/app"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/middleware"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/config"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/db"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/dto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Swagger UI HTML template loading via CDN
const swaggerHTML = `<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<title>NorthStar API Docs</title>
	<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
	</head>
	<body>
	<div id="swagger-ui"></div>
	<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
	<script>
		window.onload = () => {
		window.ui = SwaggerUIBundle({
			url: '/docs/openapi.yaml', // Points to your static YAML file
			dom_id: '#swagger-ui',
		});
		};
	</script>
	</body>
	</html>`

func SetupRouter(deps *app.AppContainer) *gin.Engine {
	r := gin.New()

	// 1. CORS CONFIGURATION (Industry Best Practice)
	// In a real production app, move the "AllowOrigins" to your config/env file
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"}, // Allows React to read your Trace ID
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Preflight request caching
	}))

	// 2. Traceability Middleware
	r.Use(middleware.RequestIDMiddleware())
	
	// Standard Gin Logger for Console
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	// 3. Exception & Logger Middleware
	r.Use(middleware.GlobalExceptionHandler())

	// --------------------------------------------------------
	// 📚 SWAGGER DOCS SETUP
	// --------------------------------------------------------
	// Expose the YAML file statically
	r.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")
	
	// Serve the Swagger UI HTML
	r.GET("/swagger", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerHTML))
	})
	// --------------------------------------------------------

	apiV1 := r.Group("/api/v1")
	{

		apiV1.GET("/health", func(c *gin.Context) {
			cfg := config.LoadConfig()
			database := db.GetDB(cfg.DatabaseURL)
			sqlDB, err := database.DB() 
			
        	response := dto.HealthResponse{
        	    Status: "ok",
        	    Details: dto.HealthDetails{
        	        DB: "ok",
        	    },
        	}

        	if err != nil || sqlDB.Ping() != nil {
        	    response.Status = "error"
        	    response.Details.DB = "down"
        	    c.JSON(http.StatusServiceUnavailable, response)
        	    return
        	}

			c.JSON(http.StatusOK, response)
		})

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

		detailRoutes := apiV1.Group("/order-details")
		{
			detailRoutes.GET("", deps.OrderCtrl.ListOrderDetails)
			detailRoutes.GET("/:id", deps.OrderCtrl.GetOrderDetailByID)
			detailRoutes.POST("", deps.OrderCtrl.CreateOrderDetail)
			detailRoutes.PUT("/:id", deps.OrderCtrl.UpdateOrderDetail)
			detailRoutes.DELETE("/:id", deps.OrderCtrl.DeleteOrderDetail)
		}
	}

	return r
}
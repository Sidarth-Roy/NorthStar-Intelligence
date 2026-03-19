// Entry point & Graceful Shutdown
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/controller"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/middleware"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/db"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, pc *controller.ProductController) {
	api := r.Group("/api/v1")
	{
		productGroup := api.Group("/products")
		{
			productGroup.POST("/", pc.Create)
			productGroup.GET("/:id", pc.GetByID)
			// productGroup.PUT("/:id", pc.Update)
			// productGroup.DELETE("/:id", pc.Delete)
		}
	}
}

func main() {
	database := db.GetDB(os.Getenv("DATABASE_URL"))
	r := gin.New()
	
	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.GlobalErrorHandler())

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Dependency Injection
	productRepo := repository.NewProductRepo(database)
	productService := service.NewProductService(productRepo)
	productController := controller.NewProductController(productService)

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/products", productController.Create)
		v1.GET("/products/:id", productController.GetByID)
	}

	// Graceful Shutdown
	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Get().Fatal().Err(err).Msg("Listen error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
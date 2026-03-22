package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/api"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/api/v1"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/repository"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/service"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/config"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/db"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. Load Config & Init Singleton Logger
	cfg := config.LoadConfig()
	logger.InitLogger()
	log := logger.Get()

	// 2. Init Singleton DB
	database := db.GetDB(cfg.DatabaseURL)

	// 3. DI Setup
	// Products
	prodRepo := repository.NewProductRepo(database)
	prodSvc := service.NewProductSvc(prodRepo)
	prodCtrl := v1.NewProductController(prodSvc)

	// Categories (New Addition)
	catRepo := repository.NewCategoryRepo(database)
	catSvc := service.NewCategorySvc(catRepo)
	catCtrl := v1.NewCategoryController(catSvc)

	// Customers
	custRepo := repository.NewCustomerRepo(database)
	custSvc := service.NewCustomerSvc(custRepo)
	custCtrl := v1.NewCustomerController(custSvc)

	// Employees
	empRepo := repository.NewEmployeeRepo(database)
	empSvc := service.NewEmployeeSvc(empRepo)
	empCtrl := v1.NewEmployeeController(empSvc)

	// Shippers
	shipRepo := repository.NewShipperRepo(database)
	shipSvc := service.NewShipperSvc(shipRepo)
	shipCtrl := v1.NewShipperController(shipSvc)

	// Orders
	orderRepo := repository.NewOrderRepo(database)
	orderSvc := service.NewOrderSvc(orderRepo)
	orderCtrl := v1.NewOrderController(orderSvc)
	
	// 4. Router Setup
	router := api.SetupRouter(prodCtrl, catCtrl, custCtrl, empCtrl, shipCtrl, orderCtrl)
	
	// 5. Server Configuration
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// 6. Start Server in Goroutine
	go func() {
		log.Info("🚀 NorthStar API starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed", zap.Error(err))
		}
	}()

	// 7. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Info("Server exited gracefully")
}
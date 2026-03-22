package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/api"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/internal/app"
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
	container := app.NewAppContainer(database)
	
	// 4. Router Setup
	router := api.SetupRouter(container)	

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
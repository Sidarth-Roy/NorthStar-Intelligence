package db

import (
	"sync"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDB(dsn string) *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true, // Cache prepared statements for speed
		})

		if err != nil {
			logger.Get().Fatal("Could not connect to database", zap.Error(err))
		}

		sqlDB, err := db.DB()
		if err != nil {
			logger.Get().Fatal("Failed to get database handle", zap.Error(err))
		}

		// Connection Pooling settings
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)

		logger.Get().Info("Database connection pool initialized")
		instance = db
	})
	return instance
}
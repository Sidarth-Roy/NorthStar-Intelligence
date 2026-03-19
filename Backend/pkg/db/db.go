// The Singleton DB & Pooling
package db

import (
	// "fmt"
	"sync"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

func GetDB(dsn string) *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
            l := logger.Get() 
            l.Fatal().Err(err).Msg("Database connection failed")
		}

		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
		
		instance = db
	})
	return instance
}
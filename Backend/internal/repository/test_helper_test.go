package repository

import (
	"fmt"
	"context"
	// "log"
	"path/filepath"
	"runtime"

	// "path/filepath"
	// "runtime"
	"testing"
	"time"

	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/DB/seeder"
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB starts a real Postgres container, migrates, and seeds data.
func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	// 1. Start Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("northwind_test"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("❌ Failed to start container: %s", err)
	}

	// 2. Get Connection String and Connect via GORM
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("❌ Failed to get connection string: %s", err)
	}

	db, err := gorm.Open(gormPostgres.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatalf("❌ Failed to connect to test db: %s", err)
	}

	// 3. Run Migrations
	// We migrate in specific order to respect Foreign Key constraints
	err = db.AutoMigrate(
		&model.Category{},
		&model.Shipper{},
		&model.Customer{},
		&model.Employee{},
		&model.Product{},
		&model.Order{},
		&model.OrderDetail{},
	)
	if err != nil {
		t.Fatalf("❌ Migration failed: %v", err)
	}

	// 4. Resolve absolute path to CSVs
	// _, b, _, _ := runtime.Caller(0)
	// basepath := filepath.Dir(b)
	// Points to Backend/DB/Northwind_Traders_Kaggle_Dataset_CSV relative to Backend/internal/repository
	// csvDir := filepath.Join(basepath, "../../DB/Northwind_Traders_Kaggle_Dataset_CSV")

	// 5. Use the variables! (Fixes the UnusedVar error)
	// log.Printf("📂 Using CSV data from: %s", csvDir)
	   
	// seeder.SeedCategories(db, csvDir)
	// seeder.SeedShippers(db, csvDir)
	// seeder.SeedProducts(db, csvDir)

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// Go up from Backend/internal/repository to Backend root, then into DB/Northwind_Traders_Kaggle_Dataset_CSV
	csvDir := filepath.Join(basepath, "../../DB/Northwind_Traders_Kaggle_Dataset_CSV")

	// Use the new function signature
	seeder.Seed(db, csvDir)

	// 3. FIX: Reset Sequences so Auto-Increment works	
    ResetSequences(db)

    return db, func() { pgContainer.Terminate(ctx) }
}

func ResetSequences(db *gorm.DB) {
	tables := []string{"categories", "shippers", "employees", "products", "orders", "order_details"}
	for _, table := range tables {
		// This SQL finds the max ID and sets the next auto-increment value to max + 1
		db.Exec(fmt.Sprintf("SELECT setval(pg_get_serial_sequence('%s', 'id'), coalesce(max(id), 1)) FROM %s", table, table))
	}
}
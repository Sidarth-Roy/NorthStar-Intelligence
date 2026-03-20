package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// Replace with your actual module path
	"github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/pkg/model" 

	"golang.org/x/text/encoding/charmap"
)

func main() {
	// Configuration with logging
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "northwind")
	port := getEnv("DB_PORT", "5432")
	sslMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbName, port, sslMode)
	
	log.Println("--- 🔌 Connecting to Database ---")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Critical: Failed to connect to database: %v", err)
	}

	log.Println("--- 🛠️  Running Migrations ---")
	err = db.Migrator().DropTable(&model.OrderDetail{}, &model.Order{}, &model.Product{}, &model.Employee{}, &model.Customer{}, &model.Shipper{}, &model.Category{})
	if err != nil {
		log.Printf("⚠️ Warning during DropTable: %v", err)
	}

	err = db.AutoMigrate(&model.Category{}, &model.Shipper{}, &model.Customer{}, &model.Employee{}, &model.Product{}, &model.Order{}, &model.OrderDetail{})
	if err != nil {
		log.Fatalf("❌ Critical: Migration failed: %v", err)
	}

	// Seeding execution
	seedCategories(db)
	seedShippers(db)
	seedCustomers(db)
	seedEmployees(db)
	seedProducts(db)
	seedOrders(db)
	seedOrderDetails(db)

	log.Println("✅ Database seeding process completed.")
}

// Helper to get environment variables
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func openCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("❌ Skipping file: Cannot open %s: %v", path, err)
		return nil
	}
	defer f.Close()

	decoder := charmap.Windows1252.NewDecoder().Reader(f)
	r := csv.NewReader(decoder)
	rows, err := r.ReadAll()
	if err != nil {
		log.Printf("❌ Error reading CSV %s: %v", path, err)
		return nil
	}
	
	if len(rows) <= 1 {
		return nil
	}
	return rows[1:]
}

func seedCategories(db *gorm.DB) {
	log.Println("📦 Seeding Categories...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/categories.csv")
	for i, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Printf("⚠️  Row %d: Invalid ID %s, skipping", i, row[0])
			continue
		}
		res := db.Create(&model.Category{
			Base:         model.Base{ID: uint(id), Active: true},
			CategoryName: row[1],
			Description:  row[2],
		})
		if res.Error != nil {
			log.Printf("❌ DB Error (Categories): %v", res.Error)
		}
	}
}

func seedShippers(db *gorm.DB) {
	log.Println("🚚 Seeding Shippers...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/shippers.csv")
	for i, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Printf("⚠️  Row %d: Invalid ID %s", i, row[0])
			continue
		}
		db.Create(&model.Shipper{
			Base:        model.Base{ID: uint(id), Active: true},
			CompanyName: row[1],
		})
	}
}

func seedCustomers(db *gorm.DB) {
	log.Println("👥 Seeding Customers...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/customers.csv")
	for _, row := range rows {
		db.Create(&model.Customer{
			Base:         model.Base{Active: true},
			CustomerID:   row[0],
			CompanyName:  row[1],
			ContactName:  row[2],
			ContactTitle: row[3],
			City:         row[4],
			Country:      row[5],
		})
	}
}

func seedEmployees(db *gorm.DB) {
	log.Println("👔 Seeding Employees...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/employees.csv")
	for _, row := range rows {
		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}
		var reportsTo *uint
		if row[5] != "" && row[5] != "NULL" {
			val, err := strconv.Atoi(row[5])
			if err == nil {
				uVal := uint(val)
				reportsTo = &uVal
			}
		}
		db.Create(&model.Employee{
			Base:         model.Base{ID: uint(id), Active: true},
			EmployeeName: row[1],
			Title:        row[2],
			City:         row[3],
			Country:      row[4],
			ReportsTo:    reportsTo,
		})
	}
}

func seedProducts(db *gorm.DB) {
	log.Println("🍎 Seeding Products...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/products.csv")
	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		price, _ := strconv.ParseFloat(row[3], 64)
		disc, _ := strconv.Atoi(row[4])
		catID, _ := strconv.Atoi(row[5])
		
		err := db.Create(&model.Product{
			Base:            model.Base{ID: uint(id), Active: true},
			ProductName:     row[1],
			QuantityPerUnit: row[2],
			UnitPrice:       price,
			Discontinued:    disc,
			CategoryID:      uint(catID),
		}).Error
		if err != nil {
			log.Printf("❌ Failed to create product %s: %v", row[1], err)
		}
	}
}

func seedOrders(db *gorm.DB) {
	log.Println("📝 Seeding Orders...")
	layout := "2006-01-02"
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/orders.csv")
	for _, row := range rows {
		id, _ := strconv.Atoi(row[0])
		eID, _ := strconv.Atoi(row[2])
		oDate, _ := time.Parse(layout, row[3])
		rDate, _ := time.Parse(layout, row[4])
		sID, _ := strconv.Atoi(row[6])
		f, _ := strconv.ParseFloat(row[7], 64)
		
		var sDate *time.Time
		if row[5] != "" && row[5] != "NULL" {
			t, err := time.Parse(layout, row[5])
			if err == nil {
				sDate = &t
			}
		}

		err := db.Create(&model.Order{
			Base:         model.Base{ID: uint(id), Active: true},
			CustomerID:   row[1],
			EmployeeID:   uint(eID),
			OrderDate:    oDate,
			RequiredDate: rDate,
			ShippedDate:  sDate,
			ShipperID:    uint(sID),
			Freight:      f,
		}).Error
		if err != nil {
			log.Printf("❌ Order ID %d insert failed: %v", id, err)
		}
	}
}

func seedOrderDetails(db *gorm.DB) {
	log.Println("🛒 Seeding Order Details...")
	rows := openCSV("Northwind_Traders_Kaggle_Dataset_CSV/order_details.csv")
	for _, row := range rows {
		oID, _ := strconv.Atoi(row[0])
		pID, _ := strconv.Atoi(row[1])
		uP, _ := strconv.ParseFloat(row[2], 64)
		q, _ := strconv.Atoi(row[3])
		d, _ := strconv.ParseFloat(row[4], 64)
		
		db.Create(&model.OrderDetail{
			Base:      model.Base{Active: true},
			OrderID:   uint(oID),
			ProductID: uint(pID),
			UnitPrice: uP,
			Quantity:  q,
			Discount:  d,
		})
	}
}
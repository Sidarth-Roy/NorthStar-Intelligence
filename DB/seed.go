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
// Grab the host from the environment (Docker), or fallback to localhost (Windows Terminal)
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "user"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "northwind"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbName, port, sslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop and Re-create tables
	db.Migrator().DropTable(&model.OrderDetail{}, &model.Order{}, &model.Product{}, &model.Employee{}, &model.Customer{}, &model.Shipper{}, &model.Category{})
	db.AutoMigrate(&model.Category{}, &model.Shipper{}, &model.Customer{}, &model.Employee{}, &model.Product{}, &model.Order{}, &model.OrderDetail{})

	seedCategories(db)
	seedShippers(db)
	seedCustomers(db)
	seedEmployees(db)
	seedProducts(db)
	seedOrders(db)
	seedOrderDetails(db)

	fmt.Println("✅ Database seeded successfully with Active status and Timestamps.")
}

func openCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("❌ Cannot open %s: %v", path, err)
	}
	defer f.Close()

	// MAGIC HAPPENS HERE: Convert Windows-1252 (Standard CSV encoding) to strict UTF-8
	decoder := charmap.Windows1252.NewDecoder().Reader(f)

	r := csv.NewReader(decoder)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatalf("❌ Error reading %s: %v", path, err)
	}
	return rows[1:]
}

func seedCategories(db *gorm.DB) {
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/categories.csv") {
		id, _ := strconv.Atoi(row[0])
		db.Create(&model.Category{
			Base:         model.Base{ID: uint(id), Active: true},
			CategoryName: row[1],
			Description:  row[2],
		})
	}
}

func seedShippers(db *gorm.DB) {
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/shippers.csv") {
		id, _ := strconv.Atoi(row[0])
		db.Create(&model.Shipper{
			Base:        model.Base{ID: uint(id), Active: true},
			CompanyName: row[1],
		})
	}
}

func seedCustomers(db *gorm.DB) {
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/customers.csv") {
		db.Create(&model.Customer{
			Base:         model.Base{Active: true}, // ID is auto-gen
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
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/employees.csv") {
		id, _ := strconv.Atoi(row[0])
		var reportsTo *uint
		if row[5] != "" && row[5] != "NULL" {
			val, _ := strconv.Atoi(row[5])
			uVal := uint(val)
			reportsTo = &uVal
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
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/products.csv") {
		id, _ := strconv.Atoi(row[0])
		price, _ := strconv.ParseFloat(row[3], 64)
		disc, _ := strconv.Atoi(row[4])
		catID, _ := strconv.Atoi(row[5])
		db.Create(&model.Product{
			Base:            model.Base{ID: uint(id), Active: true},
			ProductName:     row[1],
			QuantityPerUnit: row[2],
			UnitPrice:       price,
			Discontinued:    disc,
			CategoryID:      uint(catID),
		})
	}
}

func seedOrders(db *gorm.DB) {
	layout := "2006-01-02"
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/orders.csv") {
		id, _ := strconv.Atoi(row[0])
		eID, _ := strconv.Atoi(row[2])
		oDate, _ := time.Parse(layout, row[3])
		rDate, _ := time.Parse(layout, row[4])
		sID, _ := strconv.Atoi(row[6])
		f, _ := strconv.ParseFloat(row[7], 64)
		
		var sDate *time.Time
		if row[5] != "" && row[5] != "NULL" {
			t, _ := time.Parse(layout, row[5])
			sDate = &t
		}

		db.Create(&model.Order{
			Base:         model.Base{ID: uint(id), Active: true},
			CustomerID:   row[1],
			EmployeeID:   uint(eID),
			OrderDate:    oDate,
			RequiredDate: rDate,
			ShippedDate:  sDate,
			ShipperID:    uint(sID),
			Freight:      f,
		})
	}
}

func seedOrderDetails(db *gorm.DB) {
	for _, row := range openCSV("Northwind_Traders_Kaggle_Dataset_CSV/order_details.csv") {
		oID, _ := strconv.Atoi(row[0])
		pID, _ := strconv.Atoi(row[1])
		uP, _ := strconv.ParseFloat(row[2], 64)
		q, _ := strconv.Atoi(row[3])
		d, _ := strconv.ParseFloat(row[4], 64)
		db.Create(&model.OrderDetail{
			Base:      model.Base{Active: true}, // Surrogate ID auto-gen
			OrderID:   uint(oID),
			ProductID: uint(pID),
			UnitPrice: uP,
			Quantity:  q,
			Discount:  d,
		})
	}
}
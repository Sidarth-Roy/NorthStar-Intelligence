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
)

func main() {
	// 1. Connect to Postgres (Docker)
	dsn := "host=localhost user=user password=password dbname=northwind port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 2. Drop and Re-create tables (Clean start)
	db.Migrator().DropTable(&model.OrderDetail{}, &model.Order{}, &model.Product{}, &model.Employee{}, &model.Customer{}, &model.Shipper{}, &model.Category{})
	db.AutoMigrate(&model.Category{}, &model.Shipper{}, &model.Customer{}, &model.Employee{}, &model.Product{}, &model.Order{}, &model.OrderDetail{})

	// 3. Seed Data (Order matters for Foreign Keys)
	seedCategories(db)
	seedShippers(db)
	seedCustomers(db)
	seedEmployees(db)
	seedProducts(db)
	seedOrders(db)
	seedOrderDetails(db)

	fmt.Println("✅ Database seeded successfully.")
}

// Helper to open CSV
func openCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Cannot open %s: %v", path, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, _ := r.ReadAll()
	return rows[1:] // Skip header
}

func seedCategories(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/categories.csv") {
		id, _ := strconv.Atoi(row[0])
		db.Create(&model.Category{CategoryID: uint(id), CategoryName: row[1], Description: row[2]})
	}
}

func seedShippers(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/shippers.csv") {
		id, _ := strconv.Atoi(row[0])
		db.Create(&model.Shipper{ShipperID: uint(id), CompanyName: row[1]})
	}
}

func seedCustomers(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/customers.csv") {
		db.Create(&model.Customer{CustomerID: row[0], CompanyName: row[1], ContactName: row[2], ContactTitle: row[3], City: row[4], Country: row[5]})
	}
}

func seedEmployees(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/employees.csv") {
		id, _ := strconv.Atoi(row[0])
		var reportsTo *uint
		if row[5] != "" {
			val, _ := strconv.Atoi(row[5])
			uVal := uint(val)
			reportsTo = &uVal
		}
		db.Create(&model.Employee{EmployeeID: uint(id), EmployeeName: row[1], Title: row[2], City: row[3], Country: row[4], ReportsTo: reportsTo})
	}
}

func seedProducts(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/products.csv") {
		id, _ := strconv.Atoi(row[0])
		price, _ := strconv.ParseFloat(row[3], 64)
		disc, _ := strconv.Atoi(row[4])
		catID, _ := strconv.Atoi(row[5])
		db.Create(&model.Product{ProductID: uint(id), ProductName: row[1], QuantityPerUnit: row[2], UnitPrice: price, Discontinued: disc, CategoryID: uint(catID)})
	}
}

func seedOrders(db *gorm.DB) {
	layout := "2006-01-02"
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/orders.csv") {
		id, _ := strconv.Atoi(row[0])
		eID, _ := strconv.Atoi(row[2])
		oDate, _ := time.Parse(layout, row[3])
		rDate, _ := time.Parse(layout, row[4])
		sID, _ := strconv.Atoi(row[6])
		f, _ := strconv.ParseFloat(row[7], 64)
		
		var sDate *time.Time
		if row[5] != "" {
			t, _ := time.Parse(layout, row[5])
			sDate = &t
		}

		db.Create(&model.Order{OrderID: uint(id), CustomerID: row[1], EmployeeID: uint(eID), OrderDate: oDate, RequiredDate: rDate, ShippedDate: sDate, ShipperID: uint(sID), Freight: f})
	}
}

func seedOrderDetails(db *gorm.DB) {
	for _, row := range openCSV("DB/Northwind Traders Kaggle Dataset CSV/order_details.csv") {
		oID, _ := strconv.Atoi(row[0])
		pID, _ := strconv.Atoi(row[1])
		uP, _ := strconv.ParseFloat(row[2], 64)
		q, _ := strconv.Atoi(row[3])
		d, _ := strconv.ParseFloat(row[4], 64)
		db.Create(&model.OrderDetail{OrderID: uint(oID), ProductID: uint(pID), UnitPrice: uP, Quantity: q, Discount: d})
	}
}
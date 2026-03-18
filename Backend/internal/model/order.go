package model

import "time"

type Order struct {
	OrderID      uint          `gorm:"primaryKey" json:"orderID"`
	CustomerID   string        `json:"customerID"`
	EmployeeID   uint          `json:"employeeID"`
	OrderDate    time.Time     `json:"orderDate"`
	RequiredDate time.Time     `json:"requiredDate"`
	ShippedDate  *time.Time    `json:"shippedDate"` // Nullable
	ShipperID    uint          `json:"shipperID"`
	Freight      float64       `json:"freight"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
}

type OrderDetail struct {
	OrderID   uint    `gorm:"primaryKey" json:"orderID"`
	ProductID uint    `gorm:"primaryKey" json:"productID"`
	UnitPrice float64 `json:"unitPrice"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
}
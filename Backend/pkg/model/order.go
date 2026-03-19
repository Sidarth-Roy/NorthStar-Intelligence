package model

import "time"

type Order struct {
	Base
	CustomerID   string        `json:"customerID"`
	EmployeeID   uint          `json:"employeeID"`
	OrderDate    time.Time     `json:"orderDate"`
	RequiredDate time.Time     `json:"requiredDate"`
	ShippedDate  *time.Time    `json:"shippedDate"`
	ShipperID    uint          `json:"shipperID"`
	Freight      float64       `json:"freight"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
}

type OrderDetail struct {
	Base
	OrderID   uint    `json:"orderID"`
	ProductID uint    `json:"productID"`
	UnitPrice float64 `json:"unitPrice"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
}
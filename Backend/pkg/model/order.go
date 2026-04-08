package model

import "time"

type Order struct {
	Base
	CustomerID   string        `json:"customerID"`
	Customer     Customer      `gorm:"foreignKey:CustomerID;references:CustomerID" json:"-"` // Added
	EmployeeID   uint          `json:"employeeID"`
	Employee     Employee      `gorm:"foreignKey:EmployeeID" json:"-"` // Added
	ShipperID    uint          `json:"shipperID"`
	Shipper      Shipper       `gorm:"foreignKey:ShipperID" json:"-"`  // Added
	OrderDate    time.Time     `gorm:"type:date" json:"orderDate"`
    RequiredDate time.Time     `gorm:"type:date" json:"requiredDate"`
    ShippedDate  *time.Time    `gorm:"type:date" json:"shippedDate"`
	Freight      float64       `json:"freight"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID;references:ID" json:"orderDetails"`
}

type OrderDetail struct {
	Base
	OrderID   uint    `json:"orderID"`
	ProductID uint    `json:"productID"`
	Product   Product `gorm:"foreignKey:ProductID"` // Added
	UnitPrice float64 `json:"unitPrice"`
	Quantity  int     `json:"quantity"`
	Discount  float64 `json:"discount"`
}
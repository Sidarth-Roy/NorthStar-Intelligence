package model

type Category struct {
	CategoryID   uint      `gorm:"primaryKey" json:"categoryID"`
	CategoryName string    `json:"categoryName"`
	Description  string    `json:"description"`
	Products     []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ProductID       uint    `gorm:"primaryKey" json:"productID"`
	ProductName     string  `json:"productName"`
	QuantityPerUnit string  `json:"quantityPerUnit"`
	UnitPrice       float64 `json:"unitPrice"`
	Discontinued    int     `json:"discontinued"`
	CategoryID      uint    `json:"categoryID"`
}
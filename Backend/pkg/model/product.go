package model

type Category struct {
	Base
	CategoryName string    `gorm:"unique" json:"categoryName"`
	Description  string    `json:"description"`
	Products     []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	Base
	ProductName     string  `gorm:"unique" json:"productName"`
	QuantityPerUnit string  `json:"quantityPerUnit"`
	UnitPrice       float64 `json:"unitPrice"`
	Discontinued    int     `json:"discontinued"`
	CategoryID      uint    `json:"categoryID"`

	Category        Category `gorm:"foreignKey:CategoryID" json:"category"`
}
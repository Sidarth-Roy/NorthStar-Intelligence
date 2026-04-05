package model

type Category struct {
	Base
	CategoryName string    `gorm:"column:category_name;unique" json:"categoryName"`
	Description  string    `gorm:"column:description" json:"description"`
	Products     []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	Base
	ProductName     string  `gorm:"column:product_name;unique" json:"productName"`
	QuantityPerUnit string  `gorm:"column:quantity_per_unit" json:"quantityPerUnit"`
	UnitPrice       float64 `gorm:"column:unit_price" json:"unitPrice"`
	Discontinued    int     `gorm:"column:discontinued" json:"discontinued"`
	CategoryID      uint    `gorm:"column:category_id" json:"categoryID"`

	Category        Category `gorm:"foreignKey:CategoryID" json:"category"`
}
package model

type Customer struct {
	Base
	CustomerID   string  `gorm:"unique" json:"customerID"` // The string ID from CSV
	CompanyName  string  `gorm:"unique" json:"companyName"`
	ContactName  string  `json:"contactName"`
	ContactTitle string  `json:"contactTitle"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	Orders       []Order `gorm:"foreignKey:CustomerID;references:CustomerID"`
}
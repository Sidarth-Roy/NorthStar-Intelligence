package model

type Customer struct {
	CustomerID   string  `gorm:"primaryKey" json:"customerID"`
	CompanyName  string  `json:"companyName"`
	ContactName  string  `json:"contactName"`
	ContactTitle string  `json:"contactTitle"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	Orders       []Order `gorm:"foreignKey:CustomerID"`
}
package model

type Shipper struct {
	Base
	CompanyName string  `json:"companyName"`
	Orders      []Order `gorm:"foreignKey:ShipperID"`
}
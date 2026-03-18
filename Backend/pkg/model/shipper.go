package model

type Shipper struct {
	ShipperID   uint    `gorm:"primaryKey" json:"shipperID"`
	CompanyName string  `json:"companyName"`
	Orders      []Order `gorm:"foreignKey:ShipperID"`
}
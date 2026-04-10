package model

type Employee struct {
	Base
	EmployeeName string    `json:"employeeName"`
	Title        string    `json:"title"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	ReportsTo    *uint     `json:"reportsTo"`
	Manager      *Employee `gorm:"foreignKey:ReportsTo" json:"manager"`
	Orders       []Order   `gorm:"foreignKey:EmployeeID"`
}
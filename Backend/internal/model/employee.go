package model

type Employee struct {
	EmployeeID   uint      `gorm:"primaryKey" json:"employeeID"`
	EmployeeName string    `json:"employeeName"`
	Title        string    `json:"title"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	ReportsTo    *uint     `json:"reportsTo"` // Pointer for nulls
	Orders       []Order   `gorm:"foreignKey:EmployeeID"`
}
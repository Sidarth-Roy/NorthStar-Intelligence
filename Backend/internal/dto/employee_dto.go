package dto

// type EmployeeUpsertReq struct {
// 	EmployeeName string `json:"employeeName" binding:"required,min=3"`
// 	Title        string `json:"title" binding:"required"`
// 	City         string `json:"city"`
// 	Country      string `json:"country"`
// 	ReportsTo    *uint  `json:"reportsTo"` // Pointer allows null in JSON
// }

type EmployeeInsertReq struct {
	EmployeeName string `json:"employeeName" binding:"required,min=3"`
	Title        string `json:"title" binding:"required"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ReportsTo    *uint  `json:"reportsTo"` // Pointer allows null in JSON
}

type EmployeeUpdateReq struct {
	EmployeeName string `json:"employeeName" binding:"required,min=3"`
	Title        string `json:"title" binding:"required"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ReportsTo    *uint  `json:"reportsTo"` // Pointer allows null in JSON
}

type OrderForEmployeeNestedResponse struct {
	ID           uint      `json:"id"`
	CustomerID   string    `json:"customerID"`
	CompanyName  string    `json:"companyName"`
	OrderDate    string	   `json:"orderDate"`
	RequiredDate string	   `json:"requiredDate"`
	ShippedDate  string	   `json:"shippedDate"`
	ShipperID    uint      `json:"shipperID"`
	ShipperName  string    `json:"shipperName"`
	Freight      float64   `json:"freight"`
	Active       bool      `json:"active"`
}

type EmployeeResponse struct {
	ID           	uint   	`json:"id"`
	EmployeeName 	string 	`json:"employeeName"`
	Title        	string 	`json:"title"`
	City         	string 	`json:"city"`
	Country      	string 	`json:"country"`
	ReportsTo    	*uint  	`json:"reportsTo"`
	ReportsToName 	string 	`json:"reportsToName"`
	Active       	bool   	`json:"active"`
	// ModifiedAt   string `json:"modifiedAt"`
	Orders			[]OrderForEmployeeNestedResponse `json:"orders"`
}
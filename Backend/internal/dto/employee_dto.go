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

type EmployeeResponse struct {
	ID           	uint   	`json:"id"`
	EmployeeName 	string 	`json:"employeeName"`
	Title        	string 	`json:"title"`
	City         	string 	`json:"city"`
	Country      	string 	`json:"country"`
	ReportsTo    	*uint  	`json:"reportsTo"`
	Active       	bool   	`json:"active"`
	// ModifiedAt   string `json:"modifiedAt"`
}
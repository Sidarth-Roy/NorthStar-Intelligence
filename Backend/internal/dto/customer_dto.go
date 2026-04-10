package dto

// type CustomerUpsertReq struct {
// 	CustomerID   string `json:"customerID" binding:"required,min=5,max=5"` // Specific for Northwind IDs
// 	CompanyName  string `json:"companyName" binding:"required,min=3"`
// 	ContactName  string `json:"contactName"`
// 	ContactTitle string `json:"contactTitle"`
// 	City         string `json:"city"`
// 	Country      string `json:"country"`
// }

type CustomerInsertReq struct {
	CustomerID   string `json:"customerID" binding:"required,min=5,max=5"` // Specific for Northwind IDs
	CompanyName  string `json:"companyName" binding:"required,min=3"`
	ContactName  string `json:"contactName"`
	ContactTitle string `json:"contactTitle"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Active       bool   `json:"active"`
}

type CustomerUpdateReq struct {
	CustomerID   string `json:"customerID" binding:"required,min=5,max=5"` // Specific for Northwind IDs
	CompanyName  string `json:"companyName" binding:"required,min=3"`
	ContactName  string `json:"contactName" binding:"required"`
	ContactTitle string `json:"contactTitle" binding:"required"`
	City         string `json:"city" binding:"required"`
	Country      string `json:"country" binding:"required"`
	Active       bool   `json:"active"`
}

type CustomerResponse struct {
	ID           uint   `json:"id"`
	CustomerID   string `json:"customerID"`
	CompanyName  string `json:"companyName"`
	ContactName  string `json:"contactName"`
	ContactTitle string `json:"contactTitle"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Active       bool   `json:"active"`
	ModifiedAt   string `json:"modifiedAt"`
}
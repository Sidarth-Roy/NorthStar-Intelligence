package dto

// type ShipperUpsertReq struct {
// 	CompanyName string `json:"companyName" binding:"required,min=2"`
// }

type ShipperInsertReq struct {
	CompanyName string `json:"companyName" binding:"required,min=2"`
	Active      bool   `json:"active"`
}

type ShipperUpdateReq struct {
	CompanyName string `json:"companyName" binding:"required,min=2"`
	Active      bool   `json:"active"`
}

type OrderForShipperNestedResponse struct {
	ID           uint                  `json:"id"`
	CustomerID   uint                  `json:"customerID"`
	CompanyName  string                `json:"customerName"`
	OrderDate    string	               `json:"orderDate"`
	RequiredDate string	               `json:"requiredDate"`
	ShippedDate  string	               `json:"shippedDate"`
	ShipperID    uint                  `json:"shipperID"`
	ShipperName  string                `json:"shipperName"`
	Freight      float64               `json:"freight"`
	Active       bool                  `json:"active"`
}

type ShipperResponse struct {
	ID          uint   `json:"id"`
	CompanyName string `json:"companyName"`
	Active      bool   `json:"active"`
	Orders      []OrderForShipperNestedResponse `json:"orders"`
	// ModifiedAt  string `json:"modifiedAt"`
}
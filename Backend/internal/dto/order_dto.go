package dto

// import "time"

type OrderDetailReq struct {
	ProductID uint    `json:"productID" binding:"required"`
	UnitPrice float64 `json:"unitPrice" binding:"required,gt=0"`
	Quantity  int     `json:"quantity" binding:"required,gt=0"`
	Discount  float64 `json:"discount"`
}

// type OrderUpsertReq struct {
// 	CustomerID   string           `json:"customerID" binding:"required"`
// 	EmployeeID   uint             `json:"employeeID" binding:"required"`
// 	OrderDate    time.Time        `json:"orderDate" binding:"required"`
// 	RequiredDate time.Time        `json:"requiredDate" binding:"required"`
// 	ShippedDate  *time.Time       `json:"shippedDate"`
// 	ShipperID    uint             `json:"shipperID" binding:"required"`
// 	Freight      float64          `json:"freight" binding:"gte=0"`
// 	OrderDetails []OrderDetailReq `json:"orderDetails" binding:"required,dive,required"` // 'dive' validates the nested slice
// }

type OrderInsertReq struct {
	OrderID	  uint             `json:"orderID" binding:"required"`
	CustomerID   string           `json:"customerID" binding:"required"`
	EmployeeID   uint             `json:"employeeID" binding:"required"`
	OrderDate    string        `json:"orderDate" binding:"required"`
	RequiredDate string        `json:"requiredDate" binding:"required"`
	ShippedDate   string       `json:"shippedDate"`
	ShipperID    uint             `json:"shipperID" binding:"required"`
	Freight      float64          `json:"freight" binding:"gte=0"`
	OrderDetails []OrderDetailReq `json:"orderDetails" binding:"required,dive,required"` // 'dive' validates the nested slice
}

type OrderUpdateReq struct {
	CustomerID   string           `json:"customerID" binding:"required"`
	EmployeeID   uint             `json:"employeeID" binding:"required"`
	OrderDate    string        `json:"orderDate" binding:"required"`
	RequiredDate string        `json:"requiredDate" binding:"required"`
	ShippedDate   string       `json:"shippedDate"`
	ShipperID    uint             `json:"shipperID" binding:"required"`
	Freight      float64          `json:"freight" binding:"gte=0"`
	OrderDetails []OrderDetailReq `json:"orderDetails" binding:"required,dive,required"` // 'dive' validates the nested slice
}

type OrderDetailResponse struct {
	ID			uint    `json:"id"`
	ProductID	uint    `json:"productID"`
	ProductName string  `json:"productName"`
	UnitPrice 	float64 `json:"unitPrice"`
	Quantity  	int     `json:"quantity"`
	Discount  	float64 `json:"discount"`
}

type OrderResponse struct {
	ID           uint                  `json:"id"`
	CustomerID   string                `json:"customerID"`
	CustomerName string                `json:"customerName"`
	EmployeeID   uint                  `json:"employeeID"`
	EmployeeName string                `json:"employeeName"`
	OrderDate    string	               `json:"orderDate"`
	RequiredDate string	               `json:"requiredDate"`
	ShippedDate  string	               `json:"shippedDate"`
	ShipperID    uint                  `json:"shipperID"`
	ShipperName  string                `json:"shipperName"`
	Freight      float64               `json:"freight"`
	Active       bool                  `json:"active"`
	// ModifiedAt   string                `json:"modifiedAt"`
	OrderDetails []OrderDetailResponse `json:"orderDetails"`
}
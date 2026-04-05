package dto

type ProductUpsertReq struct {
	ProductName     string  `json:"productName" binding:"required,min=3"`
	QuantityPerUnit string  `json:"quantityPerUnit"`
	UnitPrice       float64 `json:"unitPrice" binding:"required,gt=0"`
	CategoryID      uint    `json:"categoryID" binding:"required"`
	Discontinued    int     `json:"discontinued"`
}

type ProductResponse struct {
	ID              uint    `json:"id"`
	ProductName     string  `json:"productName"`
	UnitPrice       float64 `json:"unitPrice"`
	QuantityPerUnit string  `json:"quantityPerUnit"`
	CategoryID      uint    `json:"categoryID"`
	Active          bool    `json:"active"`
	// ModifiedAt      string  `json:"modifiedAt"`
}
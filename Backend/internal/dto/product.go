package dto

import "time"

// ProductCreateRequest - Input Validation DTO
type ProductCreateRequest struct {
	ProductName     string  `json:"productName" binding:"required,min=2,max=100"`
	QuantityPerUnit string  `json:"quantityPerUnit" binding:"required"`
	UnitPrice       float64 `json:"unitPrice" binding:"required,gt=0"`
	CategoryID      uint    `json:"categoryID" binding:"required"`
	Discontinued    int     `json:"discontinued" binding:"oneof=0 1"`
}

// ProductResponse - Outbound Path DTO
type ProductResponse struct {
	ID              uint      `json:"id"`
	ProductName     string    `json:"productName"`
	UnitPrice       float64   `json:"unitPrice"`
	QuantityPerUnit string    `json:"quantityPerUnit"`
	CategoryID      uint      `json:"categoryID"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"createdAt"`
	ModifiedAt      time.Time `json:"modifiedAt"`
}
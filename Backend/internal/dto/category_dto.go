package dto

type CategoryUpsertReq struct {
	CategoryName string `json:"categoryName" binding:"required,min=2"`
	Description  string `json:"description"`
}

type CategoryResponse struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"categoryName"`
	Description  string `json:"description"`
	Active       bool   `json:"active"`
	// ModifiedAt   string `json:"modifiedAt"`
}

type ProductForCategoryResponse struct {
	ID              uint    `json:"id"`
	ProductName     string  `json:"productName"`
	UnitPrice       float64 `json:"unitPrice"`
	QuantityPerUnit string  `json:"quantityPerUnit"`
	Active          bool    `json:"active"`
}

type CategoryWithProductsResponse struct {
	ID           uint            `json:"id"`
	CategoryName string          `json:"categoryName"`
	Description  string          `json:"description"`
	Active       bool            `json:"active"`
	Products     []ProductForCategoryResponse `json:"products"`
}
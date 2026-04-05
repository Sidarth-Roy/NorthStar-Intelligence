package dto

// type CategoryUpsertReq struct {
// 	CategoryName string `json:"categoryName" binding:"required,min=2"`
// 	Description  string `json:"description"`
// }


// Required for POST /categories
type CategoryCreateReq struct {
	CategoryName string `json:"categoryName" binding:"required,min=2"`
	Description  string `json:"description"`
}

// Optional for PUT or PATCH /categories/:id
type CategoryUpdateReq struct {
	// Use 'omitempty' so it's not required, but if provided, must be min 2 chars
	CategoryName string `json:"categoryName" binding:"omitempty,min=2"`
	Description  string `json:"description"`
	// Use a pointer for bools so you can tell the difference between "false" and "missing"
	Active       *bool  `json:"active"` 
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
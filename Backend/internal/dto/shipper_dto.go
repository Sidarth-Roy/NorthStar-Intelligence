package dto

type ShipperUpsertReq struct {
	CompanyName string `json:"companyName" binding:"required,min=2"`
}

type ShipperResponse struct {
	ID          uint   `json:"id"`
	CompanyName string `json:"companyName"`
	Active      bool   `json:"active"`
	ModifiedAt  string `json:"modifiedAt"`
}
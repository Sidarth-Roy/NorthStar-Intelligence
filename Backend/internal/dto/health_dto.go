package dto

type HealthDetails struct {
	DB string `json:"DB"`
}

type HealthResponse struct {
	Status  string        `json:"status"`
	Details HealthDetails `json:"details"`
}

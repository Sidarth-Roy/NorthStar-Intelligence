package model

import (
	"time"
)

// Base adds the requested metadata to every table
type Base struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Active    bool           `gorm:"default:true" json:"active"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:modified_at" json:"modifiedAt"`
}
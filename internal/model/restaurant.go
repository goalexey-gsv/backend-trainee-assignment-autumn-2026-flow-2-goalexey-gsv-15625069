package model

import "time"

type RestaurantStatus string

const (
	RestaurantStatusActive   RestaurantStatus = "active"
	RestaurantStatusInactive RestaurantStatus = "inactive"
)

type Restaurant struct {
	ID          string           `json:"id"`
	ExternalID  string           `json:"external_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Address     string           `json:"address"`
	Status      RestaurantStatus `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

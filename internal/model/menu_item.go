package model

import "time"

type MenuItem struct {
	ID           string    `json:"id"`
	RestaurantID string    `json:"restaurant_id"`
	ExternalID   string    `json:"external_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int64     `json:"price"`
	Available    bool      `json:"available"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

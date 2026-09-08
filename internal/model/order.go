package model

import "time"

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusRejected   OrderStatus = "rejected"
	OrderStatusPreparing  OrderStatus = "preparing"
	OrderStatusReady      OrderStatus = "ready"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID              string      `json:"id"`
	RestaurantID    string      `json:"restaurant_id"`
	UserID          string      `json:"user_id"`
	Status          OrderStatus `json:"status"`
	TotalPrice      int64       `json:"total_price"`
	DeliveryAddress string      `json:"delivery_address"`
	Items           []OrderItem `json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         string `json:"id"`
	OrderID    string `json:"order_id"`
	MenuItemID string `json:"menu_item_id"`
	Name       string `json:"name"`
	Price      int64  `json:"price"`
	Quantity   int    `json:"quantity"`
}

package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
)

type CreateOrderItem struct {
	MenuItemID uuid.UUID
	Quantity   int
}

type OrderRepository interface {
	Create(
		ctx context.Context,
		restaurantID uuid.UUID,
		userID uuid.UUID,
		items []CreateOrderItem,
		deliveryAddress string,
	) (model.Order, error)

	UpdateStatus(
		ctx context.Context,
		orderID uuid.UUID,
		status model.OrderStatus,
	) error

	GetByID(
		ctx context.Context,
		orderID uuid.UUID,
	) (model.Order, error)

	GetByUserID(
		ctx context.Context,
		userID uuid.UUID,
	) ([]model.Order, error)
}

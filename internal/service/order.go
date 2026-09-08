package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/integration/restaurant"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/repository"
)

type OrderService struct {
	repository       repository.OrderRepository
	restaurantClient *restaurant.Client
}

func NewOrderService(
	repository repository.OrderRepository,
	restaurantClient *restaurant.Client,
) *OrderService {
	return &OrderService{
		repository:       repository,
		restaurantClient: restaurantClient,
	}
}

type CreateOrderInput struct {
	RestaurantID    uuid.UUID
	UserID          uuid.UUID
	DeliveryAddress string
	Items           []repository.CreateOrderItem
}

func (s *OrderService) Create(
	ctx context.Context,
	input CreateOrderInput,
) (model.Order, error) {
	order, err := s.repository.Create(
		ctx,
		input.RestaurantID,
		input.UserID,
		input.Items,
		input.DeliveryAddress,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}

	restaurantItems := make([]restaurant.OrderItem, 0, len(order.Items))

	for _, item := range order.Items {
		restaurantItems = append(restaurantItems, restaurant.OrderItem{
			ExternalID: item.MenuItemID,
			Name:       item.Name,
			Price:      item.Price,
			Quantity:   item.Quantity,
		})
	}

	response, err := s.restaurantClient.CreateOrder(
		ctx,
		restaurant.CreateOrderRequest{
			OrderID: order.ID,
			Items:   restaurantItems,
		},
	)
	if err != nil {
		_ = s.repository.UpdateStatus(
			ctx,
			uuid.MustParse(order.ID),
			model.OrderStatusRejected,
		)

		return model.Order{}, fmt.Errorf(
			"restaurant service: %w",
			err,
		)
	}

	status := model.OrderStatus(response.Status)

	if err := s.repository.UpdateStatus(
		ctx,
		uuid.MustParse(order.ID),
		status,
	); err != nil {
		return model.Order{}, fmt.Errorf(
			"update order status: %w",
			err,
		)
	}

	order.Status = status

	return order, nil
}

func (s *OrderService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (model.Order, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *OrderService) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Order, error) {
	return s.repository.GetByUserID(ctx, userID)
}

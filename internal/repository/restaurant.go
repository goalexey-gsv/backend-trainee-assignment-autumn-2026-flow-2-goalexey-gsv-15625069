package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
)

type RestaurantRepository interface {
	List(ctx context.Context) ([]model.Restaurant, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Restaurant, error)
	GetMenu(ctx context.Context, restaurantID uuid.UUID) ([]model.MenuItem, error)
}

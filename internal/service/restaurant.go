package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/repository"
)

type RestaurantService struct {
	repository repository.RestaurantRepository
}

func NewRestaurantService(repository repository.RestaurantRepository) *RestaurantService {
	return &RestaurantService{repository: repository}
}

func (s *RestaurantService) List(
	ctx context.Context,
) ([]model.Restaurant, error) {
	return s.repository.List(ctx)
}

func (s *RestaurantService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (model.Restaurant, error) {
	restaurant, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Restaurant{}, fmt.Errorf("get restaurant: %w", err)
	}

	return restaurant, nil
}

func (s *RestaurantService) GetMenu(
	ctx context.Context,
	id uuid.UUID,
) ([]model.MenuItem, error) {
	return s.repository.GetMenu(ctx, id)
}

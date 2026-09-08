package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
)

type RestaurantRepository struct {
	db *pgxpool.Pool
}

func NewRestaurantRepository(db *pgxpool.Pool) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) List(ctx context.Context) ([]model.Restaurant, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			external_id,
			name,
			description,
			address,
			status,
			created_at,
			updated_at
		FROM restaurants
		WHERE status = 'active'
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("list restaurants: %w", err)
	}
	defer rows.Close()

	restaurants := make([]model.Restaurant, 0)

	for rows.Next() {
		var restaurant model.Restaurant

		if err := rows.Scan(
			&restaurant.ID,
			&restaurant.ExternalID,
			&restaurant.Name,
			&restaurant.Description,
			&restaurant.Address,
			&restaurant.Status,
			&restaurant.CreatedAt,
			&restaurant.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan restaurant: %w", err)
		}

		restaurants = append(restaurants, restaurant)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate restaurants: %w", err)
	}

	return restaurants, nil
}

func (r *RestaurantRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (model.Restaurant, error) {
	var restaurant model.Restaurant

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			external_id,
			name,
			description,
			address,
			status,
			created_at,
			updated_at
		FROM restaurants
		WHERE id = $1
		  AND status = 'active'
	`, id).Scan(
		&restaurant.ID,
		&restaurant.ExternalID,
		&restaurant.Name,
		&restaurant.Description,
		&restaurant.Address,
		&restaurant.Status,
		&restaurant.CreatedAt,
		&restaurant.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.Restaurant{}, model.ErrNotFound
	}

	if err != nil {
		return model.Restaurant{}, fmt.Errorf("get restaurant: %w", err)
	}

	return restaurant, nil
}

func (r *RestaurantRepository) GetMenu(
	ctx context.Context,
	restaurantID uuid.UUID,
) ([]model.MenuItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			restaurant_id,
			external_id,
			name,
			description,
			price,
			available,
			created_at,
			updated_at
		FROM menu_items
		WHERE restaurant_id = $1
		  AND available = TRUE
		ORDER BY name
	`, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("get menu: %w", err)
	}
	defer rows.Close()

	items := make([]model.MenuItem, 0)

	for rows.Next() {
		var item model.MenuItem

		if err := rows.Scan(
			&item.ID,
			&item.RestaurantID,
			&item.ExternalID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Available,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan menu item: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate menu: %w", err)
	}

	return items, nil
}

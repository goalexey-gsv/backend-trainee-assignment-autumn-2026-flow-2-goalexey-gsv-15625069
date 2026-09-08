package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/model"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-goalexey-gsv-15625069/internal/repository"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(
	ctx context.Context,
	restaurantID uuid.UUID,
	userID uuid.UUID,
	requestItems []repository.CreateOrderItem,
	deliveryAddress string,
) (model.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Order{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ids := make([]uuid.UUID, 0, len(requestItems))

	for _, item := range requestItems {
		ids = append(ids, item.MenuItemID)
	}

	rows, err := tx.Query(ctx, `
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
		WHERE id = ANY($1::uuid[])
		  AND restaurant_id = $2
		FOR UPDATE
	`, ids, restaurantID)
	if err != nil {
		return model.Order{}, fmt.Errorf("get menu items: %w", err)
	}
	defer rows.Close()

	menuItems := make(map[uuid.UUID]model.MenuItem)

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
			return model.Order{}, fmt.Errorf("scan menu item: %w", err)
		}

		menuItems[uuid.MustParse(item.ID)] = item
	}

	if err := rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("iterate menu items: %w", err)
	}

	if len(menuItems) != len(requestItems) {
		return model.Order{}, fmt.Errorf("one or more menu items not found")
	}

	orderItems := make([]model.OrderItem, 0, len(requestItems))
	var totalPrice int64

	for _, requestItem := range requestItems {
		menuItem, ok := menuItems[requestItem.MenuItemID]
		if !ok {
			return model.Order{}, fmt.Errorf(
				"menu item %s not found",
				requestItem.MenuItemID,
			)
		}

		if !menuItem.Available {
			return model.Order{}, fmt.Errorf(
				"menu item %s is unavailable",
				requestItem.MenuItemID,
			)
		}

		itemTotal := menuItem.Price * int64(requestItem.Quantity)
		totalPrice += itemTotal

		orderItems = append(orderItems, model.OrderItem{
			MenuItemID: menuItem.ID,
			Name:       menuItem.Name,
			Price:      menuItem.Price,
			Quantity:   requestItem.Quantity,
		})
	}

	var order model.Order

	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			restaurant_id,
			user_id,
			status,
			total_price,
			delivery_address
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			restaurant_id,
			user_id,
			status,
			total_price,
			delivery_address,
			created_at,
			updated_at
	`,
		restaurantID,
		userID,
		model.OrderStatusPending,
		totalPrice,
		deliveryAddress,
	).Scan(
		&order.ID,
		&order.RestaurantID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.DeliveryAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, fmt.Errorf("create order: %w", err)
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID

		_, err = tx.Exec(ctx, `
			INSERT INTO order_items (
				order_id,
				menu_item_id,
				name,
				price,
				quantity
			)
			VALUES ($1, $2, $3, $4, $5)
		`,
			order.ID,
			orderItems[i].MenuItemID,
			orderItems[i].Name,
			orderItems[i].Price,
			orderItems[i].Quantity,
		)
		if err != nil {
			return model.Order{}, fmt.Errorf("create order item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Order{}, fmt.Errorf("commit transaction: %w", err)
	}

	order.Items = orderItems

	return order, nil
}

func (r *OrderRepository) UpdateStatus(
	ctx context.Context,
	orderID uuid.UUID,
	status model.OrderStatus,
) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE orders
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
	`, orderID, status)
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (r *OrderRepository) GetByID(
	ctx context.Context,
	orderID uuid.UUID,
) (model.Order, error) {
	var order model.Order

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			restaurant_id,
			user_id,
			status,
			total_price,
			delivery_address,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`, orderID).Scan(
		&order.ID,
		&order.RestaurantID,
		&order.UserID,
		&order.Status,
		&order.TotalPrice,
		&order.DeliveryAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return model.Order{}, model.ErrNotFound
	}

	if err != nil {
		return model.Order{}, fmt.Errorf("get order: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			order_id,
			menu_item_id,
			name,
			price,
			quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`, orderID)
	if err != nil {
		return model.Order{}, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	order.Items = make([]model.OrderItem, 0)

	for rows.Next() {
		var item model.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.MenuItemID,
			&item.Name,
			&item.Price,
			&item.Quantity,
		); err != nil {
			return model.Order{}, fmt.Errorf("scan order item: %w", err)
		}

		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return model.Order{}, fmt.Errorf("iterate order items: %w", err)
	}

	return order, nil
}

func (r *OrderRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			restaurant_id,
			user_id,
			status,
			total_price,
			delivery_address,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get user orders: %w", err)
	}
	defer rows.Close()

	orders := make([]model.Order, 0)

	for rows.Next() {
		var order model.Order

		if err := rows.Scan(
			&order.ID,
			&order.RestaurantID,
			&order.UserID,
			&order.Status,
			&order.TotalPrice,
			&order.DeliveryAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate user orders: %w", err)
	}

	return orders, nil
}

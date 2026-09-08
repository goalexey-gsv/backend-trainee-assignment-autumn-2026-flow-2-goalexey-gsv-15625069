-- +goose Up

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id),
    user_id UUID NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'pending',
    total_price BIGINT NOT NULL,
    delivery_address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT orders_status_check
        CHECK (
            status IN (
                'pending',
                'accepted',
                'rejected',
                'preparing',
                'ready',
                'delivering',
                'delivered',
                'cancelled'
            )
        ),

    CONSTRAINT orders_total_price_check
        CHECK (total_price >= 0)
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id UUID NOT NULL REFERENCES menu_items(id),
    name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    quantity INTEGER NOT NULL,

    CONSTRAINT order_items_price_check
        CHECK (price >= 0),

    CONSTRAINT order_items_quantity_check
        CHECK (quantity > 0)
);

CREATE INDEX idx_orders_restaurant_id
    ON orders (restaurant_id);

CREATE INDEX idx_orders_user_id
    ON orders (user_id);

CREATE INDEX idx_orders_status
    ON orders (status);

CREATE INDEX idx_order_items_order_id
    ON order_items (order_id);


-- +goose Down

DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

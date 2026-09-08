-- +goose Up

INSERT INTO restaurants (
    id,
    external_id,
    name,
    description,
    address,
    status
)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'restaurant-001',
    'Burger House',
    'Бургеры, картофель и напитки',
    'Москва, ул. Пушкина, 10',
    'active'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menu_items (
    id,
    restaurant_id,
    external_id,
    name,
    description,
    price,
    available
)
VALUES
(
    '21111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    'burger-001',
    'Классический бургер',
    'Говяжья котлета, сыр, салат и соус',
    45000,
    TRUE
),
(
    '22222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    'burger-002',
    'Чизбургер',
    'Говяжья котлета, двойной сыр и соус',
    52000,
    TRUE
),
(
    '23333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111',
    'fries-001',
    'Картофель фри',
    'Картофель фри, стандартная порция',
    18000,
    TRUE
),
(
    '24444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111',
    'cola-001',
    'Кола',
    'Газированный напиток 0.5 л',
    12000,
    TRUE
)
ON CONFLICT (id) DO NOTHING;

-- +goose Down

DELETE FROM menu_items
WHERE restaurant_id = '11111111-1111-1111-1111-111111111111';

DELETE FROM restaurants
WHERE id = '11111111-1111-1111-1111-111111111111';
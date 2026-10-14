-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    total_amount,
    status
) VALUES (
    $1, $2, $3
)
RETURNING id, user_id, total_amount, status, created_at, updated_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id,
    product_id,
    quantity,
    unit_price,
    subtotal
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, order_id, product_id, quantity, unit_price, subtotal, created_at;

-- name: GetOrderByID :one
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE id = $1
LIMIT 1;

-- name: ListOrderItemsByOrderID :many
SELECT id, order_id, product_id, quantity, unit_price, subtotal, created_at
FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;

-- name: ListOrdersByUserID :many
SELECT id, user_id, total_amount, status, created_at, updated_at
FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, total_amount, status, created_at, updated_at;

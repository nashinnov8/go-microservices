-- name: CreateProduct :one
INSERT INTO products (
    sku,
    name,
    description,
    price,
    stock,
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, sku, name, description, price, stock, is_active, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, sku, name, description, price, stock, is_active, created_at, updated_at
FROM products
WHERE id = $1
LIMIT 1;

-- name: GetProductBySKU :one
SELECT id, sku, name, description, price, stock, is_active, created_at, updated_at
FROM products
WHERE sku = $1
LIMIT 1;

-- name: ListProducts :many
SELECT id, sku, name, description, price, stock, is_active, created_at, updated_at
FROM products
WHERE ($1::boolean IS NULL OR is_active = $1)
  AND ($2::text = '' OR name ILIKE '%' || $2 || '%' OR sku ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountProducts :one
SELECT COUNT(*)
FROM products
WHERE ($1::boolean IS NULL OR is_active = $1)
  AND ($2::text = '' OR name ILIKE '%' || $2 || '%' OR sku ILIKE '%' || $2 || '%');

-- name: UpdateProduct :one
UPDATE products
SET
    name = $2,
    description = $3,
    price = $4,
    is_active = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, sku, name, description, price, stock, is_active, created_at, updated_at;

-- name: UpdateProductStock :one
UPDATE products
SET
    stock = stock + $2,
    updated_at = NOW()
WHERE id = $1 AND (stock + $2) >= 0
RETURNING id, sku, name, description, price, stock, is_active, created_at, updated_at;

-- name: DeleteProduct :execrows
DELETE FROM products
WHERE id = $1;

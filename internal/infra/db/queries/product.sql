-- name: GetProductByActivityID :one
SELECT * FROM product
WHERE activity_id = ?;

-- name: ListProducts :many
SELECT * FROM product
WHERE
  (? = '' OR platform = ?) AND
  (? = '' OR region = ?) AND
  (? IS NULL OR sales_status = ?) AND
  (? IS NULL OR ? = 0 OR activity_create_time >= datetime('now', '-7 days')) AND
  activity_id NOT IN (SELECT activity_id FROM blocked_product)
ORDER BY activity_create_time DESC;

-- name: ListProductsWithBlockedStatus :many
SELECT p.* FROM product p
INNER JOIN blocked_product b ON p.activity_id = b.activity_id
ORDER BY p.activity_create_time DESC;

-- name: CreateProduct :exec
INSERT INTO product (
  activity_id, platform, region, title, shop_name,
  original_price, current_price, sales_status, activity_create_time
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateProduct :exec
UPDATE product
SET current_price = ?,
    sales_status = ?,
    update_time = datetime('now')
WHERE id = ?;

-- name: UpdateProductByActivityID :exec
UPDATE product
SET current_price = ?,
    sales_status = ?,
    update_time = datetime('now')
WHERE activity_id = ?;

-- name: DeleteProduct :exec
DELETE FROM product WHERE id = ?;

-- name: DeleteByActivityIDs :exec
-- Delete multiple products by activity IDs
-- Note: IN clause with multiple values handled in Go code
DELETE FROM product WHERE activity_id = ?;

-- name: DeleteByPlatform :exec
DELETE FROM product WHERE platform = ?;

-- name: CountByPlatform :one
SELECT COUNT(*) FROM product WHERE platform = ?;

-- name: GetBlockedProduct :one
SELECT * FROM blocked_product
WHERE activity_id = ? AND user_id = ?;

-- name: ListBlockedProductsByUser :many
SELECT activity_id FROM blocked_product WHERE user_id = ?;

-- name: ListAllBlockedProducts :many
SELECT * FROM blocked_product;

-- name: CreateBlockedProduct :exec
INSERT OR IGNORE INTO blocked_product (activity_id, user_id)
VALUES (?, ?);

-- name: DeleteBlockedProduct :exec
DELETE FROM blocked_product WHERE activity_id = ? AND user_id = ?;

-- name: ExistsBlockedProduct :one
SELECT COUNT(*) > 0 AS is_blocked FROM blocked_product WHERE activity_id = ? AND user_id = ?;

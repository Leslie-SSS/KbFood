-- name: GetTrendByActivityIDAndDate :one
SELECT * FROM product_price_trend
WHERE activity_id = ? AND record_date = ?;

-- name: ListTrendsByActivityID :many
SELECT * FROM product_price_trend
WHERE activity_id = ?
ORDER BY record_date ASC;

-- name: CreateTrend :exec
INSERT INTO product_price_trend (activity_id, price, record_date)
VALUES (?, ?, ?)
ON CONFLICT (activity_id, record_date) DO UPDATE
SET price = MIN(product_price_trend.price, excluded.price);

-- name: DeleteTrendsByActivityIDs :exec
-- Delete multiple trends by activity IDs
-- Note: IN clause with multiple values handled in Go code
DELETE FROM product_price_trend WHERE activity_id = ?;

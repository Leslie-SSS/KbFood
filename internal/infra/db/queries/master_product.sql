-- name: GetMasterProductByID :one
SELECT * FROM master_product
WHERE id = ?;

-- name: ListMasterProductsByRegion :many
SELECT * FROM master_product
WHERE region = ?
ORDER BY update_time DESC;

-- name: ListMasterProductsByPlatform :many
SELECT * FROM master_product
WHERE platform = ?
ORDER BY update_time DESC;

-- name: ListMasterProductsByRegionAndPlatform :many
SELECT * FROM master_product
WHERE region = ? AND platform = ?
ORDER BY update_time DESC;

-- name: ListAllMasterProducts :many
SELECT * FROM master_product
ORDER BY update_time DESC;

-- name: CreateMasterProduct :exec
INSERT INTO master_product (id, region, platform, standard_title, price, status, trust_score)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: UpdateMasterProduct :exec
UPDATE master_product
SET price = ?,
    status = ?,
    trust_score = ?,
    update_time = datetime('now')
WHERE id = ?;

-- name: UpdateMasterProductPlatform :exec
UPDATE master_product
SET platform = ?,
    update_time = datetime('now')
WHERE id = ?;

-- name: DeleteMasterProduct :exec
DELETE FROM master_product WHERE id = ?;

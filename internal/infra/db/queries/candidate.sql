-- name: GetCandidateByID :one
SELECT * FROM candidate_item
WHERE id = ?;

-- name: ListCandidatesByRegion :many
SELECT * FROM candidate_item
WHERE region = ?
ORDER BY last_seen_time DESC;

-- name: ListAllCandidates :many
SELECT * FROM candidate_item
ORDER BY last_seen_time DESC;

-- name: CreateCandidate :exec
INSERT INTO candidate_item (group_key, region, title_votes, total_occurrences, last_price, last_status, first_seen_time, last_seen_time)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateCandidate :exec
UPDATE candidate_item
SET title_votes = ?,
    total_occurrences = ?,
    last_price = ?,
    last_status = ?,
    last_seen_time = ?,
    update_time = datetime('now')
WHERE id = ?;

-- name: DeleteCandidate :exec
DELETE FROM candidate_item WHERE id = ?;

-- name: DeleteCandidatesByIDs :exec
-- Delete multiple candidates by IDs
-- Note: IN clause with multiple values handled in Go code
DELETE FROM candidate_item WHERE id = ?;

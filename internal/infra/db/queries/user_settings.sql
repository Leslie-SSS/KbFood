-- name: GetUserSettings :one
SELECT * FROM user_settings WHERE user_id = ?;

-- name: UpsertUserSettings :exec
INSERT INTO user_settings (user_id, bark_key)
VALUES (?, ?)
ON CONFLICT (user_id) DO UPDATE
SET bark_key = excluded.bark_key,
    update_time = datetime('now');

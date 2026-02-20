-- name: GetNotification :one
SELECT * FROM notification_config
WHERE activity_id = ? AND user_id = ?;

-- name: ListNotificationsByUser :many
SELECT * FROM notification_config WHERE user_id = ?;

-- name: ListAllNotifications :many
SELECT * FROM notification_config;

-- name: UpsertNotification :exec
INSERT INTO notification_config (activity_id, user_id, target_price, last_notify_time)
VALUES (?, ?, ?, ?)
ON CONFLICT (activity_id, user_id) DO UPDATE
SET target_price = excluded.target_price,
    last_notify_time = excluded.last_notify_time,
    update_time = datetime('now');

-- name: UpdateNotificationNotifyTime :exec
UPDATE notification_config
SET last_notify_time = datetime('now'),
    update_time = datetime('now')
WHERE activity_id = ? AND user_id = ?;

-- name: DeleteNotification :exec
DELETE FROM notification_config WHERE activity_id = ? AND user_id = ?;

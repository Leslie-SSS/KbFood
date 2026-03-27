package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	dbinfra "kbfood/internal/infra/db"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type userSettingsRow struct {
	UserID     string
	BarkKey    string
	CreateTime string
	UpdateTime string
}

type notificationRow struct {
	ActivityID     string
	UserID         string
	TargetPrice    float64
	LastNotifyTime string
	CreateTime     string
	UpdateTime     string
}

type blockedProductRow struct {
	ActivityID string
	UserID     string
	CreateTime string
}

// UserDataMigrator migrates legacy Bark-key-based user data onto the stable client ID.
func UserDataMigrator(database *dbinfra.Pool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			currentUserID := GetUserID(c)
			legacyUserID := GetLegacyUserID(c)

			if err := migrateLegacyUserData(c.Request().Context(), database, currentUserID, legacyUserID); err != nil {
				log.Error().
					Err(err).
					Str("currentUserID", currentUserID).
					Str("legacyUserID", legacyUserID).
					Msg("failed to migrate legacy user data")
			}

			return next(c)
		}
	}
}

func migrateLegacyUserData(
	ctx context.Context,
	database *dbinfra.Pool,
	currentUserID, legacyUserID string,
) error {
	currentUserID = strings.TrimSpace(currentUserID)
	legacyUserID = strings.TrimSpace(legacyUserID)

	if database == nil || currentUserID == "" || legacyUserID == "" || currentUserID == legacyUserID {
		return nil
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	hasLegacyData, err := hasLegacyUserData(ctx, tx, legacyUserID)
	if err != nil {
		return err
	}
	if !hasLegacyData {
		return nil
	}

	currentSettings, err := loadUserSettingsRow(ctx, tx, currentUserID)
	if err != nil {
		return err
	}
	legacySettings, err := loadUserSettingsRow(ctx, tx, legacyUserID)
	if err != nil {
		return err
	}

	mergedSettings := mergeUserSettings(currentUserID, currentSettings, legacySettings)
	mergedNotifications, err := mergeNotifications(ctx, tx, currentUserID, legacyUserID)
	if err != nil {
		return err
	}
	mergedBlockedProducts, err := mergeBlockedProducts(ctx, tx, currentUserID, legacyUserID)
	if err != nil {
		return err
	}

	if err := replaceUserSettings(ctx, tx, currentUserID, legacyUserID, mergedSettings); err != nil {
		return err
	}
	if err := replaceNotifications(ctx, tx, currentUserID, legacyUserID, mergedNotifications); err != nil {
		return err
	}
	if err := replaceBlockedProducts(ctx, tx, currentUserID, legacyUserID, mergedBlockedProducts); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration tx: %w", err)
	}

	log.Info().
		Str("currentUserID", currentUserID).
		Str("legacyUserID", legacyUserID).
		Int("notifications", len(mergedNotifications)).
		Int("blockedProducts", len(mergedBlockedProducts)).
		Bool("migratedSettings", mergedSettings != nil).
		Msg("migrated legacy user data")

	return nil
}

func hasLegacyUserData(ctx context.Context, tx *sql.Tx, legacyUserID string) (bool, error) {
	var count int
	row := tx.QueryRowContext(ctx, `
		SELECT
			(SELECT COUNT(*) FROM user_settings WHERE user_id = ?) +
			(SELECT COUNT(*) FROM notification_config WHERE user_id = ?) +
			(SELECT COUNT(*) FROM blocked_product WHERE user_id = ?)
	`, legacyUserID, legacyUserID, legacyUserID)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("check legacy user data: %w", err)
	}

	return count > 0, nil
}

func loadUserSettingsRow(ctx context.Context, tx *sql.Tx, userID string) (*userSettingsRow, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT user_id, bark_key, create_time, update_time
		FROM user_settings
		WHERE user_id = ?
	`, userID)

	var settings userSettingsRow
	if err := row.Scan(&settings.UserID, &settings.BarkKey, &settings.CreateTime, &settings.UpdateTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("load user settings: %w", err)
	}

	return &settings, nil
}

func loadNotifications(ctx context.Context, tx *sql.Tx, currentUserID, legacyUserID string) ([]notificationRow, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT activity_id, user_id, target_price, COALESCE(last_notify_time, ''), create_time, update_time
		FROM notification_config
		WHERE user_id IN (?, ?)
	`, currentUserID, legacyUserID)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []notificationRow
	for rows.Next() {
		var item notificationRow
		if err := rows.Scan(
			&item.ActivityID,
			&item.UserID,
			&item.TargetPrice,
			&item.LastNotifyTime,
			&item.CreateTime,
			&item.UpdateTime,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		notifications = append(notifications, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}

	return notifications, nil
}

func loadBlockedProducts(ctx context.Context, tx *sql.Tx, currentUserID, legacyUserID string) ([]blockedProductRow, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT activity_id, user_id, create_time
		FROM blocked_product
		WHERE user_id IN (?, ?)
	`, currentUserID, legacyUserID)
	if err != nil {
		return nil, fmt.Errorf("query blocked products: %w", err)
	}
	defer rows.Close()

	var blockedProducts []blockedProductRow
	for rows.Next() {
		var item blockedProductRow
		if err := rows.Scan(&item.ActivityID, &item.UserID, &item.CreateTime); err != nil {
			return nil, fmt.Errorf("scan blocked product: %w", err)
		}
		blockedProducts = append(blockedProducts, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate blocked products: %w", err)
	}

	return blockedProducts, nil
}

func mergeUserSettings(currentUserID string, current, legacy *userSettingsRow) *userSettingsRow {
	preferred := current
	if shouldReplaceSettings(preferred, legacy) {
		preferred = legacy
	}
	if preferred == nil {
		return nil
	}

	return &userSettingsRow{
		UserID:     currentUserID,
		BarkKey:    strings.TrimSpace(preferred.BarkKey),
		CreateTime: firstNonEmpty(minTimestamp(currentCreateTime(current), currentCreateTime(legacy)), preferred.CreateTime),
		UpdateTime: firstNonEmpty(maxTimestamp(currentUpdateTime(current), currentUpdateTime(legacy)), preferred.UpdateTime),
	}
}

func currentCreateTime(settings *userSettingsRow) string {
	if settings == nil {
		return ""
	}
	return settings.CreateTime
}

func currentUpdateTime(settings *userSettingsRow) string {
	if settings == nil {
		return ""
	}
	return settings.UpdateTime
}

func shouldReplaceSettings(existing, candidate *userSettingsRow) bool {
	if candidate == nil {
		return false
	}
	if existing == nil {
		return true
	}

	existingHasBark := strings.TrimSpace(existing.BarkKey) != ""
	candidateHasBark := strings.TrimSpace(candidate.BarkKey) != ""

	if candidateHasBark && !existingHasBark {
		return true
	}
	if !candidateHasBark && existingHasBark {
		return false
	}

	return compareTimestamp(candidate.UpdateTime, existing.UpdateTime) > 0
}

func mergeNotifications(
	ctx context.Context,
	tx *sql.Tx,
	currentUserID, legacyUserID string,
) ([]notificationRow, error) {
	rows, err := loadNotifications(ctx, tx, currentUserID, legacyUserID)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]notificationRow, len(rows))
	for _, item := range rows {
		existing, ok := merged[item.ActivityID]
		if !ok || shouldReplaceNotification(currentUserID, existing, item) {
			item.UserID = currentUserID
			merged[item.ActivityID] = item
		}
	}

	result := make([]notificationRow, 0, len(merged))
	for _, item := range merged {
		result = append(result, item)
	}

	return result, nil
}

func shouldReplaceNotification(currentUserID string, existing, candidate notificationRow) bool {
	if compareTimestamp(candidate.UpdateTime, existing.UpdateTime) > 0 {
		return true
	}
	if compareTimestamp(candidate.UpdateTime, existing.UpdateTime) < 0 {
		return false
	}
	return candidate.UserID == currentUserID && existing.UserID != currentUserID
}

func mergeBlockedProducts(
	ctx context.Context,
	tx *sql.Tx,
	currentUserID, legacyUserID string,
) ([]blockedProductRow, error) {
	rows, err := loadBlockedProducts(ctx, tx, currentUserID, legacyUserID)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]blockedProductRow, len(rows))
	for _, item := range rows {
		item.UserID = currentUserID
		existing, ok := merged[item.ActivityID]
		if !ok || compareTimestamp(item.CreateTime, existing.CreateTime) < 0 {
			merged[item.ActivityID] = item
		}
	}

	result := make([]blockedProductRow, 0, len(merged))
	for _, item := range merged {
		result = append(result, item)
	}

	return result, nil
}

func replaceUserSettings(
	ctx context.Context,
	tx *sql.Tx,
	currentUserID, legacyUserID string,
	settings *userSettingsRow,
) error {
	if settings != nil {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO user_settings (user_id, bark_key, create_time, update_time)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(user_id) DO UPDATE SET
				bark_key = excluded.bark_key,
				update_time = excluded.update_time
		`, currentUserID, settings.BarkKey, defaultTimestamp(settings.CreateTime), defaultTimestamp(settings.UpdateTime))
		if err != nil {
			return fmt.Errorf("upsert migrated user settings: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM user_settings WHERE user_id = ?`, legacyUserID); err != nil {
		return fmt.Errorf("delete legacy user settings: %w", err)
	}

	return nil
}

func replaceNotifications(
	ctx context.Context,
	tx *sql.Tx,
	currentUserID, legacyUserID string,
	rows []notificationRow,
) error {
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM notification_config
		WHERE user_id IN (?, ?)
	`, currentUserID, legacyUserID); err != nil {
		return fmt.Errorf("delete existing notifications for migration: %w", err)
	}

	for _, item := range rows {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO notification_config (
				activity_id, user_id, target_price, last_notify_time, create_time, update_time
			) VALUES (?, ?, ?, ?, ?, ?)
		`,
			item.ActivityID,
			currentUserID,
			item.TargetPrice,
			nullIfEmpty(item.LastNotifyTime),
			defaultTimestamp(item.CreateTime),
			defaultTimestamp(item.UpdateTime),
		)
		if err != nil {
			return fmt.Errorf("insert migrated notification %s: %w", item.ActivityID, err)
		}
	}

	return nil
}

func replaceBlockedProducts(
	ctx context.Context,
	tx *sql.Tx,
	currentUserID, legacyUserID string,
	rows []blockedProductRow,
) error {
	if _, err := tx.ExecContext(ctx, `
		DELETE FROM blocked_product
		WHERE user_id IN (?, ?)
	`, currentUserID, legacyUserID); err != nil {
		return fmt.Errorf("delete existing blocked products for migration: %w", err)
	}

	for _, item := range rows {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO blocked_product (activity_id, user_id, create_time)
			VALUES (?, ?, ?)
		`, item.ActivityID, currentUserID, defaultTimestamp(item.CreateTime))
		if err != nil {
			return fmt.Errorf("insert migrated blocked product %s: %w", item.ActivityID, err)
		}
	}

	return nil
}

func compareTimestamp(left, right string) int {
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)

	if left == right {
		return 0
	}
	if left == "" {
		return -1
	}
	if right == "" {
		return 1
	}

	leftTime, leftOK := parseTimestamp(left)
	rightTime, rightOK := parseTimestamp(right)
	if leftOK && rightOK {
		switch {
		case leftTime.After(rightTime):
			return 1
		case leftTime.Before(rightTime):
			return -1
		default:
			return 0
		}
	}

	return strings.Compare(left, right)
}

func parseTimestamp(value string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, true
		}
	}

	return time.Time{}, false
}

func minTimestamp(left, right string) string {
	switch compareTimestamp(left, right) {
	case 1:
		return right
	case -1:
		return left
	default:
		return firstNonEmpty(left, right)
	}
}

func maxTimestamp(left, right string) string {
	switch compareTimestamp(left, right) {
	case 1:
		return left
	case -1:
		return right
	default:
		return firstNonEmpty(left, right)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func defaultTimestamp(value string) string {
	if strings.TrimSpace(value) == "" {
		return time.Now().Format("2006-01-02 15:04:05")
	}
	return value
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}

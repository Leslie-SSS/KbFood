package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"kbfood/internal/domain/entity"
	"kbfood/internal/domain/repository"
	"kbfood/internal/interface/http/dto"
	"kbfood/internal/interface/http/middleware"

	"github.com/labstack/echo/v4"
)

// UserHandler handles user-related requests
type UserHandler struct {
	userSettingsRepo repository.UserSettingsRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(userSettingsRepo repository.UserSettingsRepository) *UserHandler {
	return &UserHandler{
		userSettingsRepo: userSettingsRepo,
	}
}

// SaveSettings handles POST /api/user/settings
func (h *UserHandler) SaveSettings(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	if userID == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "用户标识缺失，请刷新页面后重试"))
	}

	var params struct {
		BarkKey string `json:"barkKey"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Invalid request body"))
	}

	// Normalize Bark key - extract device key from URL if needed
	barkKey := normalizeBarkKey(params.BarkKey)

	if barkKey == "" {
		return c.JSON(http.StatusBadRequest, dto.Error(400, "Bark Key 不能为空"))
	}

	settings := &entity.UserSettings{
		UserID:  userID,
		BarkKey: barkKey,
	}

	if err := h.userSettingsRepo.Upsert(ctx, settings); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to save settings"))
	}

	return c.JSON(http.StatusOK, dto.Success(map[string]string{
		"barkKey": barkKey,
	}))
}

// GetSettings handles GET /api/user/settings
func (h *UserHandler) GetSettings(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	if userID == "" {
		return c.JSON(http.StatusOK, dto.Success(map[string]string{
			"barkKey": "",
		}))
	}

	settings, err := h.userSettingsRepo.Get(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Error(500, "Failed to get settings"))
	}

	if settings == nil {
		return c.JSON(http.StatusOK, dto.Success(map[string]string{
			"barkKey": "",
		}))
	}

	return c.JSON(http.StatusOK, dto.Success(map[string]string{
		"barkKey": settings.BarkKey,
	}))
}

// normalizeBarkKey extracts device key from full URL or returns key as-is
func normalizeBarkKey(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// If full URL, extract the last non-empty path segment as device key.
	if strings.HasPrefix(strings.ToLower(input), "http") {
		parsed, err := url.Parse(input)
		if err == nil {
			path := strings.Trim(parsed.Path, "/")
			if path == "" {
				return ""
			}

			parts := strings.Split(path, "/")
			for i := len(parts) - 1; i >= 0; i-- {
				if seg := strings.TrimSpace(parts[i]); seg != "" {
					return seg
				}
			}
			return ""
		}

		trimmed := strings.Trim(input, "/")
		parts := strings.Split(trimmed, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if seg := strings.TrimSpace(parts[i]); seg != "" {
				return seg
			}
		}
		return ""
	}
	return input
}

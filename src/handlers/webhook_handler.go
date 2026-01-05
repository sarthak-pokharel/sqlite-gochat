package handlers

import (
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"

	"github.com/labstack/echo/v4"
)

// WebhookHandler handles incoming webhooks from external platforms
type WebhookHandler struct {
	service services.WebhookService
}

func NewWebhookHandler(service services.WebhookService) *WebhookHandler {
	return &WebhookHandler{
		service: service,
	}
}

// HandleWebhook handles POST /api/v1/webhooks/:channelId/:platform
func (h *WebhookHandler) HandleWebhook(c echo.Context) error {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid channel ID")
	}

	platform := c.Param("platform")
	if platform == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "platform is required")
	}

	// Parse generic webhook payload
	var payload map[string]interface{}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
	}

	// Determine event type (platform-specific logic would go here)
	eventType := "message" // Default
	if et, ok := payload["event_type"].(string); ok {
		eventType = et
	}

	// Process webhook
	if err := h.service.ProcessWebhook(channelID, eventType, payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "processed",
	})
}

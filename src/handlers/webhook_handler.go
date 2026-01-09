package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/utils"

	"github.com/go-chi/chi/v5"
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

// HandleWebhook handles POST /api/v1/webhooks/{channelId}/{platform}
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.ParseInt(chi.URLParam(r, "channelId"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	platform := chi.URLParam(r, "platform")
	if platform == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "platform is required")
		return
	}

	// Parse generic webhook payload
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid payload")
		return
	}

	// Determine event type (platform-specific logic would go here)
	eventType := "message" // Default
	if et, ok := payload["event_type"].(string); ok {
		eventType = et
	}

	// Process webhook
	if err := h.service.ProcessWebhook(channelID, eventType, payload); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"status": "processed",
	})
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// MessageHandler handles message HTTP requests
type MessageHandler struct {
	service   services.MessageService
	validator *validator.Validate
}

func NewMessageHandler(service services.MessageService) *MessageHandler {
	return &MessageHandler{
		service:   service,
		validator: validator.New(),
	}
}

// SendMessage handles POST /api/v1/conversations/{id}/messages
func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	var req struct {
		Content     string             `json:"content" validate:"required"`
		MessageType models.MessageType `json:"message_type"`
		MediaURL    *string            `json:"media_url,omitempty"`
		Metadata    *string            `json:"metadata,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if req.MessageType == "" {
		req.MessageType = models.MessageTypeText
	}

	message, err := h.service.SendOutgoingMessage(&services.SendOutgoingMessageRequest{
		ConversationID: conversationID,
		Content:        req.Content,
		MessageType:    req.MessageType,
		MediaURL:       req.MediaURL,
		Metadata:       req.Metadata,
	})

	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, message)
}

// GetHistory handles GET /api/v1/conversations/{id}/messages
func (h *MessageHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	conversationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	beforeStr := r.URL.Query().Get("before")

	var before *int64
	if beforeStr != "" {
		b, err := strconv.ParseInt(beforeStr, 10, 64)
		if err == nil {
			before = &b
		}
	}

	messages, err := h.service.GetMessageHistory(conversationID, limit, offset, before)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":   messages,
		"limit":  limit,
		"offset": offset,
	})
}

// MarkDelivered handles POST /api/v1/messages/{id}/delivered
func (h *MessageHandler) MarkDelivered(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	if err := h.service.MarkDelivered(messageID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "marked as delivered",
	})
}

// MarkRead handles POST /api/v1/messages/{id}/read
func (h *MessageHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	messageID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid message ID")
		return
	}

	if err := h.service.MarkRead(messageID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "marked as read",
	})
}

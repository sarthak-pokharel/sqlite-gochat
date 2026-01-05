package handlers

import (
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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

// SendMessage handles POST /api/v1/conversations/:id/messages
func (h *MessageHandler) SendMessage(c echo.Context) error {
	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	var req struct {
		Content     string             `json:"content" validate:"required"`
		MessageType models.MessageType `json:"message_type"`
		MediaURL    *string            `json:"media_url,omitempty"`
		Metadata    *string            `json:"metadata,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, message)
}

// GetHistory handles GET /api/v1/conversations/:id/messages
func (h *MessageHandler) GetHistory(c echo.Context) error {
	conversationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	beforeStr := c.QueryParam("before")

	var before *int64
	if beforeStr != "" {
		beforeVal, err := strconv.ParseInt(beforeStr, 10, 64)
		if err == nil {
			before = &beforeVal
		}
	}

	if limit == 0 {
		limit = 50
	}

	messages, err := h.service.GetMessageHistory(conversationID, limit, offset, before)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   messages,
		"limit":  limit,
		"offset": offset,
	})
}

// MarkDelivered handles POST /api/v1/messages/:id/delivered
func (h *MessageHandler) MarkDelivered(c echo.Context) error {
	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid message ID")
	}

	if err := h.service.MarkDelivered(messageID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "message marked as delivered",
	})
}

// MarkRead handles POST /api/v1/messages/:id/read
func (h *MessageHandler) MarkRead(c echo.Context) error {
	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid message ID")
	}

	if err := h.service.MarkRead(messageID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "message marked as read",
	})
}

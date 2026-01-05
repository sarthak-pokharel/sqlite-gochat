package handlers

import (
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"

	"github.com/labstack/echo/v4"
)

// ConversationHandler handles conversation HTTP requests
type ConversationHandler struct {
	service services.ConversationService
}

func NewConversationHandler(service services.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: service,
	}
}

// GetByID handles GET /api/v1/conversations/:id
func (h *ConversationHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	conversation, err := h.service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "conversation not found")
	}

	return c.JSON(http.StatusOK, conversation)
}

// List handles GET /api/v1/channels/:channelId/conversations
func (h *ConversationHandler) List(c echo.Context) error {
	channelID, err := strconv.ParseInt(c.Param("channelId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid channel ID")
	}

	statusStr := c.QueryParam("status")
	var status *models.ConversationStatus
	if statusStr != "" {
		s := models.ConversationStatus(statusStr)
		status = &s
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	conversations, err := h.service.ListByChannel(channelID, status, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   conversations,
		"limit":  limit,
		"offset": offset,
	})
}

// Assign handles POST /api/v1/conversations/:id/assign
func (h *ConversationHandler) Assign(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	var req struct {
		AssigneeID string `json:"assignee_id" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.AssigneeID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "assignee_id is required")
	}

	if err := h.service.Assign(id, req.AssigneeID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "conversation assigned successfully",
	})
}

// UpdateStatus handles PATCH /api/v1/conversations/:id/status
func (h *ConversationHandler) UpdateStatus(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	var req struct {
		Status models.ConversationStatus `json:"status" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "conversation status updated successfully",
	})
}

// UpdatePriority handles PATCH /api/v1/conversations/:id/priority
func (h *ConversationHandler) UpdatePriority(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid conversation ID")
	}

	var req struct {
		Priority models.ConversationPriority `json:"priority" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.service.UpdatePriority(id, req.Priority); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "conversation priority updated successfully",
	})
}

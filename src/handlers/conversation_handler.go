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

// ConversationHandler handles conversation HTTP requests
type ConversationHandler struct {
	service   services.ConversationService
	validator *validator.Validate
}

func NewConversationHandler(service services.ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service:   service,
		validator: validator.New(),
	}
}

// GetByID handles GET /api/v1/conversations/{id}
func (h *ConversationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	conversation, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "conversation not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, conversation)
}

// List handles GET /api/v1/channels/{channelId}/conversations
func (h *ConversationHandler) List(w http.ResponseWriter, r *http.Request) {
	channelID, err := strconv.ParseInt(chi.URLParam(r, "channelId"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	statusStr := r.URL.Query().Get("status")
	var status *models.ConversationStatus
	if statusStr != "" {
		s := models.ConversationStatus(statusStr)
		status = &s
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit == 0 {
		limit = 20
	}

	conversations, err := h.service.ListByChannel(channelID, status, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":   conversations,
		"limit":  limit,
		"offset": offset,
	})
}

// Assign handles POST /api/v1/conversations/{id}/assign
func (h *ConversationHandler) Assign(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	var req struct {
		AssigneeID string `json:"assignee_id" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AssigneeID == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "assignee_id is required")
		return
	}

	if err := h.service.Assign(id, req.AssigneeID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "conversation assigned successfully",
	})
}

// UpdateStatus handles PATCH /api/v1/conversations/{id}/status
func (h *ConversationHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	var req struct {
		Status models.ConversationStatus `json:"status" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Status == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "status is required")
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "conversation status updated successfully",
	})
}

// UpdatePriority handles PATCH /api/v1/conversations/{id}/priority
func (h *ConversationHandler) UpdatePriority(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid conversation ID")
		return
	}

	var req struct {
		Priority models.ConversationPriority `json:"priority" validate:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Priority == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "priority is required")
		return
	}

	if err := h.service.UpdatePriority(id, req.Priority); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "conversation priority updated successfully",
	})
}

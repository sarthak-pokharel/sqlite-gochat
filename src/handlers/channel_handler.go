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

// ChannelHandler handles channel HTTP requests
type ChannelHandler struct {
	service   services.ChannelService
	validator *validator.Validate
}

func NewChannelHandler(service services.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create handles POST /api/v1/channels
func (h *ChannelHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	channel, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, channel)
}

// GetByID handles GET /api/v1/channels/{id}
func (h *ChannelHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	channel, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "channel not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, channel)
}

// ListByOrganization handles GET /api/v1/organizations/{orgId}/channels
func (h *ChannelHandler) ListByOrganization(w http.ResponseWriter, r *http.Request) {
	orgID, err := strconv.ParseInt(chi.URLParam(r, "orgId"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid organization ID")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	channels, err := h.service.ListByOrganization(orgID, limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":   channels,
		"limit":  limit,
		"offset": offset,
	})
}

// Update handles PATCH /api/v1/channels/{id}
func (h *ChannelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	var req models.UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.service.Update(id, &req); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}

// UpdateStatus handles PATCH /api/v1/channels/{id}/status
func (h *ChannelHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	var req struct {
		Status models.ChannelStatus `json:"status" validate:"required,oneof=active inactive error pending"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "updated"})
}

// Delete handles DELETE /api/v1/channels/{id}
func (h *ChannelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid channel ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "deleted"})
}

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

// OrganizationHandler handles organization HTTP requests
type OrganizationHandler struct {
	service   services.OrganizationService
	validator *validator.Validate
}

func NewOrganizationHandler(service services.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{
		service:   service,
		validator: validator.New(),
	}
}

// Create handles POST /api/v1/organizations
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		utils.ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	org, err := h.service.Create(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, org)
}

// GetByID handles GET /api/v1/organizations/{id}
func (h *OrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid organization ID")
		return
	}

	org, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "organization not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, org)
}

// GetBySlug handles GET /api/v1/organizations/slug/{slug}
func (h *OrganizationHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "slug is required")
		return
	}

	org, err := h.service.GetBySlug(slug)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "organization not found")
		return
	}

	utils.JSONResponse(w, http.StatusOK, org)
}

// List handles GET /api/v1/organizations
func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	orgs, err := h.service.List(limit, offset)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"data":   orgs,
		"limit":  limit,
		"offset": offset,
	})
}

// Update handles PATCH /api/v1/organizations/{id}
func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid organization ID")
		return
	}

	var req models.UpdateOrganizationRequest
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

	utils.JSONResponse(w, http.StatusOK, map[string]string{
		"message": "organization updated successfully",
	})
}

// Delete handles DELETE /api/v1/organizations/{id}
func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid organization ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

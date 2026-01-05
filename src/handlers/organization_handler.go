package handlers

import (
	"net/http"
	"strconv"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
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
func (h *OrganizationHandler) Create(c echo.Context) error {
	var req models.CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	org, err := h.service.Create(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, org)
}

// GetByID handles GET /api/v1/organizations/:id
func (h *OrganizationHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid organization ID")
	}

	org, err := h.service.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "organization not found")
	}

	return c.JSON(http.StatusOK, org)
}

// GetBySlug handles GET /api/v1/organizations/slug/:slug
func (h *OrganizationHandler) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "slug is required")
	}

	org, err := h.service.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "organization not found")
	}

	return c.JSON(http.StatusOK, org)
}

// List handles GET /api/v1/organizations
func (h *OrganizationHandler) List(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	if limit == 0 {
		limit = 20
	}

	orgs, err := h.service.List(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":   orgs,
		"limit":  limit,
		"offset": offset,
	})
}

// Update handles PATCH /api/v1/organizations/:id
func (h *OrganizationHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid organization ID")
	}

	var req models.UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	if err := h.service.Update(id, &req); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "organization updated successfully",
	})
}

// Delete handles DELETE /api/v1/organizations/:id
func (h *OrganizationHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid organization ID")
	}

	if err := h.service.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

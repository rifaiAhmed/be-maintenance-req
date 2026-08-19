package api

import (
	"backend-test/internal/interfaces"
	"backend-test/internal/models/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct{ Service interfaces.RequestService }

func NewRequestHandler(service interfaces.RequestService) *RequestHandler {
	return &RequestHandler{Service: service}
}

// List godoc
// @Summary List maintenance requests
// @Tags Maintenance Requests
// @Produce json
// @Param search query string false "Search keyword"
// @Param status query string false "Status"
// @Param priority query string false "Priority"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /requests [get]
func (h *RequestHandler) List(c *gin.Context) {
	page, limit := pagination(c)
	filter := dto.RequestFilter{Search: c.Query("search"), Status: c.Query("status"), Priority: c.Query("priority"), Page: page, Limit: limit}
	items, meta, err := h.Service.List(c.Request.Context(), actor(c), filter)
	if err != nil {
		respondError(c, err)
		return
	}
	respondList(c, items, meta)
}

// Get godoc
// @Summary Get maintenance request by ID
// @Tags Maintenance Requests
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /requests/{id} [get]
func (h *RequestHandler) Get(c *gin.Context) {
	item, err := h.Service.Get(c.Request.Context(), actor(c), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", item)
}

// Create godoc
// @Summary Create maintenance request
// @Tags Maintenance Requests
// @Accept json
// @Produce json
// @Param request body dto.CreateMaintenanceRequest true "Request payload"
// @Success 201 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Security BearerAuth
// @Router /requests [post]
func (h *RequestHandler) Create(c *gin.Context) {
	var req dto.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request data", "errors": err.Error()})
		return
	}
	item, err := h.Service.Create(c.Request.Context(), actor(c), req)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusCreated, "request created successfully", item)
}

// Update godoc
// @Summary Update maintenance request
// @Tags Maintenance Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.UpdateMaintenanceRequest true "Request payload"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Security BearerAuth
// @Router /requests/{id} [put]
func (h *RequestHandler) Update(c *gin.Context) {
	var req dto.UpdateMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request data", "errors": err.Error()})
		return
	}
	item, err := h.Service.Update(c.Request.Context(), actor(c), c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "request updated successfully", item)
}

// Review godoc
// @Summary Review maintenance request
// @Tags Maintenance Requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param request body dto.ReviewMaintenanceRequest true "Review payload"
// @Success 200 {object} map[string]interface{} "success"
// @Failure 400 {object} map[string]interface{} "invalid request"
// @Security BearerAuth
// @Router /requests/{id}/review [patch]
func (h *RequestHandler) Review(c *gin.Context) {
	var req dto.ReviewMaintenanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "status must be Approved or Rejected", "errors": err.Error()})
		return
	}
	item, err := h.Service.Review(c.Request.Context(), actor(c), c.Param("id"), req)
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "request reviewed successfully", item)
}

// Delete godoc
// @Summary Delete maintenance request
// @Tags Maintenance Requests
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /requests/{id} [delete]
func (h *RequestHandler) Delete(c *gin.Context) {
	if err := h.Service.Delete(c.Request.Context(), actor(c), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "request deleted successfully", nil)
}

// Dashboard godoc
// @Summary Get dashboard summary
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{} "success"
// @Security BearerAuth
// @Router /dashboard [get]
func (h *RequestHandler) Dashboard(c *gin.Context) {
	counts, recent, err := h.Service.Dashboard(c.Request.Context(), actor(c))
	if err != nil {
		respondError(c, err)
		return
	}
	respond(c, http.StatusOK, "success", gin.H{"counts": counts, "recentRequests": recent})
}

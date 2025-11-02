package http

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driver"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SLAHandler is a thin HTTP adapter that delegates to the SLAService
type SLAHandler struct {
	service driver.SLAService
}

// NewSLAHandler creates a new SLA handler
func NewSLAHandler(service driver.SLAService) *SLAHandler {
	return &SLAHandler{
		service: service,
	}
}

// GetSLAs returns all SLAs
func (h *SLAHandler) GetSLAs(c *gin.Context) {
	slas, err := h.service.RetrieveSLAs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slas)
}

// GetSLA returns an SLA by id
func (h *SLAHandler) GetSLA(c *gin.Context) {
	id := c.Param("id")
	sla, err := h.service.RetrieveSLA(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "sla not found"})
		return
	}
	c.JSON(http.StatusOK, sla)
}

// PostSLA creates a new SLA
func (h *SLAHandler) PostSLA(c *gin.Context) {
	var payload domain.SLA
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sla, err := h.service.CreateSLA(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sla)
}

// PutSLA updates an existing SLA
func (h *SLAHandler) PutSLA(c *gin.Context) {
	id := c.Param("id")
	var payload domain.SLA
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.ID != "" && payload.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in path and body must match (or omit body id)"})
		return
	}
	payload.ID = id

	sla, err := h.service.UpdateSLA(payload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sla)
}

// DeleteSLA removes an SLA by id
func (h *SLAHandler) DeleteSLA(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSLA(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

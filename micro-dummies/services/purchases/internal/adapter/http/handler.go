package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/application"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/purchases/internal/core/port/driver"
)

type PurchaseHandler struct {
	service driver.PurchaseService
}

func NewPurchaseHandler(service driver.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{service: service}
}

func PongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *PurchaseHandler) GetPurchases(c *gin.Context) {
	purchases, err := h.service.RetrievePurchases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchases)
}

func (h *PurchaseHandler) GetPurchase(c *gin.Context) {
	id := c.Param("id")

	purchase, err := h.service.RetrievePurchase(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if purchase == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) PostPurchase(c *gin.Context) {
	var req struct {
		Price    float64 `json:"price" binding:"required"`
		Quantity int     `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.service.CreatePurchase(req.Price, req.Quantity)
	if err != nil {
		if errors.Is(err, application.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, purchase)
}

func (h *PurchaseHandler) PutPurchase(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Price    float64 `json:"price" binding:"required"`
		Quantity int     `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	purchase, err := h.service.UpdatePurchase(id, req.Price, req.Quantity)
	if err != nil {
		if errors.Is(err, application.ErrPurchaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, application.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchase)
}

func (h *PurchaseHandler) DeletePurchase(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeletePurchase(id)
	if err != nil {
		if errors.Is(err, application.ErrPurchaseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

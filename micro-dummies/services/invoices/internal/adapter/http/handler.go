package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/application"
	"github.com/medisupply/medisupply-infrastructure/micro-dummies/services/invoices/internal/core/port/driver"
)

type InvoiceHandler struct {
	service driver.InvoiceService
}

func NewInvoiceHandler(service driver.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func PongHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	invoices, err := h.service.RetrieveInvoices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoices)
}

func (h *InvoiceHandler) GetInvoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.service.RetrieveInvoice(id)
	if err != nil {
		if errors.Is(err, application.ErrInvoiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if invoice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) PostInvoice(c *gin.Context) {
	var req struct {
		PurchasesID []string `json:"purchases" binding:"required"`
		Buyer       string   `json:"buyer" binding:"required"`
		Subtotal    float64  `json:"subtotal" binding:"required"`
		Discount    float64  `json:"discount"`
		Taxes       float64  `json:"taxes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := h.service.CreateInvoice(req.PurchasesID, req.Buyer, req.Subtotal, req.Discount, req.Taxes)
	if err != nil {
		if errors.Is(err, application.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

func (h *InvoiceHandler) PutInvoice(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PurchasesID []string `json:"purchases" binding:"required"`
		Buyer       string   `json:"buyer" binding:"required"`
		Subtotal    float64  `json:"subtotal" binding:"required"`
		Discount    float64  `json:"discount"`
		Taxes       float64  `json:"taxes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := h.service.UpdateInvoice(id, req.PurchasesID, req.Buyer, req.Subtotal, req.Discount, req.Taxes)
	if err != nil {
		if errors.Is(err, application.ErrInvoiceNotFound) {
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

	c.JSON(http.StatusOK, invoice)
}

func (h *InvoiceHandler) DeleteInvoice(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteInvoice(id)
	if err != nil {
		if errors.Is(err, application.ErrInvoiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

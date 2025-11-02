package http

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driver"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomerHandler is a thin HTTP adapter that delegates to the CustomerService
type CustomerHandler struct {
	service driver.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(service driver.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

// GetCustomers returns all customers
func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	customers, err := h.service.RetrieveCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// GetCustomer returns a customer by id
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.service.RetrieveCustomer(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// PostCustomer creates a new customer
func (h *CustomerHandler) PostCustomer(c *gin.Context) {
	var payload domain.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.service.CreateCustomer(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// PutCustomer updates an existing customer
func (h *CustomerHandler) PutCustomer(c *gin.Context) {
	id := c.Param("id")
	var payload domain.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.ID != "" && payload.ID != id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id in path and body must match (or omit body id)"})
		return
	}
	payload.ID = id

	customer, err := h.service.UpdateCustomer(payload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer removes a customer by id
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

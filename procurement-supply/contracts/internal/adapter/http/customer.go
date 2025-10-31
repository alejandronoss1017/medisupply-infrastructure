package http

import (
	"contracts/internal/core/domain"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CustomerHandler manages in-memory CRUD for customers
type CustomerHandler struct {
	mu        sync.RWMutex
	customers map[string]domain.Customer
	idSeq     int64
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{
		customers: make(map[string]domain.Customer),
		idSeq:     time.Now().UnixNano(),
	}
}

func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	list := make([]domain.Customer, 0, len(h.customers))
	for _, v := range h.customers {
		list = append(list, v)
	}
	c.JSON(http.StatusOK, list)
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id := c.Param("id")
	h.mu.RLock()
	defer h.mu.RUnlock()
	if v, ok := h.customers[id]; ok {
		c.JSON(http.StatusOK, v)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
}

func (h *CustomerHandler) PostCustomer(c *gin.Context) {
	var payload domain.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.ID == "" {
		h.mu.Lock()
		h.idSeq++
		payload.ID = strconv.FormatInt(h.idSeq, 10)
		h.customers[payload.ID] = payload
		h.mu.Unlock()
	} else {
		h.mu.Lock()
		if _, exists := h.customers[payload.ID]; exists {
			h.mu.Unlock()
			c.JSON(http.StatusConflict, gin.H{"error": "customer with given id already exists"})
			return
		}
		h.customers[payload.ID] = payload
		h.mu.Unlock()
	}
	c.JSON(http.StatusCreated, payload)
}

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
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.customers[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	h.customers[id] = payload
	c.JSON(http.StatusOK, payload)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.customers[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	delete(h.customers, id)
	c.Status(http.StatusNoContent)
}

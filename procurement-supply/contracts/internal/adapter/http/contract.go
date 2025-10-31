package http

import (
	"contracts/internal/core/domain"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	mu        sync.RWMutex
	contracts map[string]domain.Contract
	idSeq     int64
}

func NewContractHandler() *ContractHandler {
	return &ContractHandler{
		contracts: make(map[string]domain.Contract),
		idSeq:     time.Now().UnixNano(),
	}
}

// GetContracts returns all contracts
func (h *ContractHandler) GetContracts(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	list := make([]domain.Contract, 0, len(h.contracts))
	for _, v := range h.contracts {
		list = append(list, v)
	}
	c.JSON(http.StatusOK, list)
}

// GetContract returns a contract by id
func (h *ContractHandler) GetContract(c *gin.Context) {
	id := c.Param("id")
	h.mu.RLock()
	defer h.mu.RUnlock()
	if v, ok := h.contracts[id]; ok {
		c.JSON(http.StatusOK, v)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
}

// PostContract creates a new contract
func (h *ContractHandler) PostContract(c *gin.Context) {
	var payload domain.Contract
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if empty
	if payload.ID == "" {
		h.mu.Lock()
		h.idSeq++
		payload.ID = strconv.FormatInt(h.idSeq, 10)
		h.contracts[payload.ID] = payload
		h.mu.Unlock()
	} else {
		// If provided ID already exists, reject
		h.mu.Lock()
		if _, exists := h.contracts[payload.ID]; exists {
			h.mu.Unlock()
			c.JSON(http.StatusConflict, gin.H{"error": "contract with given id already exists"})
			return
		}
		h.contracts[payload.ID] = payload
		h.mu.Unlock()
	}

	c.JSON(http.StatusCreated, payload)
}

// PutContract updates an existing contract (full replacement)
func (h *ContractHandler) PutContract(c *gin.Context) {
	id := c.Param("id")
	var payload domain.Contract
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
	if _, ok := h.contracts[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	h.contracts[id] = payload
	c.JSON(http.StatusOK, payload)
}

// DeleteContract removes a contract by id
func (h *ContractHandler) DeleteContract(c *gin.Context) {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.contracts[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	delete(h.contracts, id)
	c.Status(http.StatusNoContent)
}

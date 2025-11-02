package http

import (
	"contracts/internal/core/domain"
	"contracts/internal/core/port/driver"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContractHandler is a thin HTTP adapter that delegates to the ContractService
type ContractHandler struct {
	service driver.ContractService
}

// NewContractHandler creates a new contract handler
func NewContractHandler(service driver.ContractService) *ContractHandler {
	return &ContractHandler{
		service: service,
	}
}

// GetContracts returns all contracts
func (h *ContractHandler) GetContracts(c *gin.Context) {
	contracts, err := h.service.RetrieveContracts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contracts)
}

// GetContract returns a contract by id
func (h *ContractHandler) GetContract(c *gin.Context) {
	id := c.Param("id")
	contract, err := h.service.RetrieveContract(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "contract not found"})
		return
	}
	c.JSON(http.StatusOK, contract)
}

// PostContract creates a new contract
func (h *ContractHandler) PostContract(c *gin.Context) {
	var payload domain.Contract
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract, err := h.service.CreateContract(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, contract)
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

	contract, err := h.service.UpdateContract(payload)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// DeleteContract removes a contract by id
func (h *ContractHandler) DeleteContract(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteContract(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

package http

import (
	"fmt"
	"net/http"
	"novelties/internal/core/domain"
	"novelties/internal/core/port/driven"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NoveltyHandler struct {
	writer driven.BlockchainWriter
	logger *zap.SugaredLogger
}

func NewNoveltyHandler(writer driven.BlockchainWriter, logger *zap.SugaredLogger) *NoveltyHandler {
	return &NoveltyHandler{
		writer: writer,
		logger: logger,
	}
}

func (h *NoveltyHandler) PostNovelty(c *gin.Context) {
	var payload domain.Novelty
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receipt, err := h.writer.CheckSLA(c, payload.ContractID, payload.SLAID, payload.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to check SLA to blockchain: %w", err)})
		return
	}
	if receipt.Status != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("blockchain transaction failed (status: %d, tx: %s)", receipt.Status, receipt.TxHash)})
		return
	}

	h.logger.Infow("Successfully checked SLA on blockchain",
		"contract_id", payload.ContractID,
		"sla_id", payload.SLAID,
		"tx_hash", receipt.TxHash,
		"block_number", receipt.BlockNumber,
		"gas_used", receipt.GasUsed,
	)

	c.JSON(http.StatusCreated, payload)

}

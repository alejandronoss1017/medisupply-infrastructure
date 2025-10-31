package http

import (
	"contracts/internal/core/domain"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// SLAHandler manages in-memory CRUD for SLAs
type SLAHandler struct {
	mu    sync.RWMutex
	slas  map[string]domain.SLA
	idSeq int64
}

func NewSLAHandler() *SLAHandler {
	return &SLAHandler{
		slas:  make(map[string]domain.SLA),
		idSeq: time.Now().UnixNano(),
	}
}

func (h *SLAHandler) GetSLAs(c *gin.Context) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	list := make([]domain.SLA, 0, len(h.slas))
	for _, v := range h.slas {
		list = append(list, v)
	}
	c.JSON(http.StatusOK, list)
}

func (h *SLAHandler) GetSLA(c *gin.Context) {
	id := c.Param("id")
	h.mu.RLock()
	defer h.mu.RUnlock()
	if v, ok := h.slas[id]; ok {
		c.JSON(http.StatusOK, v)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "sla not found"})
}

func (h *SLAHandler) PostSLA(c *gin.Context) {
	var payload domain.SLA
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.ID == "" {
		h.mu.Lock()
		h.idSeq++
		payload.ID = strconv.FormatInt(h.idSeq, 10)
		h.slas[payload.ID] = payload
		h.mu.Unlock()
	} else {
		h.mu.Lock()
		if _, exists := h.slas[payload.ID]; exists {
			h.mu.Unlock()
			c.JSON(http.StatusConflict, gin.H{"error": "sla with given id already exists"})
			return
		}
		h.slas[payload.ID] = payload
		h.mu.Unlock()
	}
	c.JSON(http.StatusCreated, payload)
}

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
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.slas[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "sla not found"})
		return
	}
	h.slas[id] = payload
	c.JSON(http.StatusOK, payload)
}

func (h *SLAHandler) DeleteSLA(c *gin.Context) {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.slas[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "sla not found"})
		return
	}
	delete(h.slas, id)
	c.Status(http.StatusNoContent)
}

package http

import (
	"net/http"
	"suppliers/internal/core/domain"
	"suppliers/internal/core/port/driver"

	"github.com/gin-gonic/gin"
)

type MedicineHandler struct {
	service driver.MedicineService
}

func NewMedicineHandler(service driver.MedicineService) *MedicineHandler {
	return &MedicineHandler{service: service}
}

func (h *MedicineHandler) GetMedicines(c *gin.Context) {

	medicines, err := h.service.RetrieveMedicines()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, medicines)
}

func (h *MedicineHandler) GetMedicine(c *gin.Context) {

	id := c.Param("id")

	medicine, err := h.service.RetrieveMedicine(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, medicine)
}

func (h *MedicineHandler) PostMedicine(c *gin.Context) {
	var medicine domain.Medicine

	// Bind JSON request to Medicine struct and validate
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Create the medicine using the service
	createdMedicine, err := h.service.CreateMedicine(&medicine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medicine: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdMedicine)
}

func (h *MedicineHandler) PutMedicine(c *gin.Context) {
	id := c.Param("id")
	var medicine domain.Medicine

	// Bind JSON request to Medicine struct and validate
	if err := c.ShouldBindJSON(&medicine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Update the medicine using the service
	updatedMedicine, err := h.service.UpdateMedicine(id, &medicine)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medicine: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedMedicine)
}

func (h *MedicineHandler) DeleteMedicine(c *gin.Context) {
	id := c.Param("id")

	// Delete the medicine using the service
	err := h.service.DeleteMedicine(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete medicine: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}

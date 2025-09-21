package handlers

import (
	"go_pz3/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PerformanceHandler містить логіку для обробки запитів до вистав
type PerformanceHandler struct {
	DB *gorm.DB
}

func NewPerformanceHandler(db *gorm.DB) *PerformanceHandler {
	return &PerformanceHandler{DB: db}
}

// CreatePerformance створює нову виставу
func (h *PerformanceHandler) CreatePerformance(c *gin.Context) {
	var performance models.Performance
	if err := c.ShouldBindJSON(&performance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&performance); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, performance)
}

// GetPerformances повертає список усіх вистав
func (h *PerformanceHandler) GetPerformances(c *gin.Context) {
	var performances []models.Performance
	if result := h.DB.Find(&performances); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, performances)
}

// GetPerformance повертає одну виставу за ID
func (h *PerformanceHandler) GetPerformance(c *gin.Context) {
	var performance models.Performance
	id := c.Param("id")

	if result := h.DB.First(&performance, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Performance not found"})
		return
	}

	c.JSON(http.StatusOK, performance)
}

// UpdatePerformance оновлює дані вистави
func (h *PerformanceHandler) UpdatePerformance(c *gin.Context) {
	var performance models.Performance
	id := c.Param("id")

	if result := h.DB.First(&performance, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Performance not found"})
		return
	}

	if err := c.ShouldBindJSON(&performance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Save(&performance)
	c.JSON(http.StatusOK, performance)
}

// DeletePerformance видаляє виставу
func (h *PerformanceHandler) DeletePerformance(c *gin.Context) {
	id := c.Param("id")
	if result := h.DB.Delete(&models.Performance{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Performance not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

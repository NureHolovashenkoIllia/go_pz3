package handlers

import (
	"go_pz3/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RehearsalHandler struct {
	DB *gorm.DB
}

func NewRehearsalHandler(db *gorm.DB) *RehearsalHandler {
	return &RehearsalHandler{DB: db}
}

// CreateRehearsalInput - Структура для вхідних даних при створенні репетиції
type CreateRehearsalInput struct {
	PerformanceID uint   `json:"performance_id" binding:"required"`
	DateTime      string `json:"date_time" binding:"required"`
	ActorIDs      []uint `json:"actor_ids" binding:"required"`
}

// CreateRehearsal створює репетицію та призначає на неї акторів
func (h *RehearsalHandler) CreateRehearsal(c *gin.Context) {
	var input CreateRehearsalInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Створюємо саму репетицію
	rehearsal := models.Rehearsal{
		PerformanceID: input.PerformanceID,
		DateTime:      parseTime(input.DateTime),
	}
	if result := h.DB.Create(&rehearsal); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Знаходимо акторів за їх ID
	var actors []*models.Actor
	if err := h.DB.Find(&actors, input.ActorIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find actors"})
		return
	}

	// Прив'язуємо акторів до репетиції
	if err := h.DB.Model(&rehearsal).Association("Actors").Append(actors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate actors"})
		return
	}

	// Використовуємо Preload, щоб повернути репетицію разом з акторами
	h.DB.Preload("Actors").First(&rehearsal, rehearsal.ID)
	c.JSON(http.StatusCreated, rehearsal)
}

// GetRehearsals повертає всі репетиції з акторами, що беруть участь
func (h *RehearsalHandler) GetRehearsals(c *gin.Context) {
	var rehearsals []models.Rehearsal
	// Preload("Actors") завантажує пов'язаних акторів для кожної репетиції
	if result := h.DB.Preload("Actors").Find(&rehearsals); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, rehearsals)
}

// GetRehearsal повертає одну репетицію за ID
func (h *RehearsalHandler) GetRehearsal(c *gin.Context) {
	var rehearsal models.Rehearsal
	id := c.Param("id")

	if result := h.DB.Preload("Actors").First(&rehearsal, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rehearsal not found"})
		return
	}

	c.JSON(http.StatusOK, rehearsal)
}

// UpdateRehearsal оновлює дані репетиції
func (h *RehearsalHandler) UpdateRehearsal(c *gin.Context) {
	var rehearsal models.Rehearsal
	id := c.Param("id")

	if result := h.DB.First(&rehearsal, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rehearsal not found"})
		return
	}

	if err := c.ShouldBindJSON(&rehearsal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Save(&rehearsal)
	c.JSON(http.StatusOK, rehearsal)
}

// DeleteRehearsal видаляє репетицію
func (h *RehearsalHandler) DeleteRehearsal(c *gin.Context) {
	id := c.Param("id")
	if result := h.DB.Delete(&models.Rehearsal{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rehearsal not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// Допоміжна функція для парсингу часу (для простоти)
func parseTime(timeStr string) time.Time {
	layout := "2006-01-02T15:04:05Z" // ISO 8601
	t, _ := time.Parse(layout, timeStr)
	return t
}

package handlers

import (
	"errors"
	"go_pz3/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	DB *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{DB: db}
}

// ActorRoleCountResult - Структура для повернення результату про актора та кількість його ролей
type ActorRoleCountResult struct {
	ActorID   uint   `json:"actor_id"`
	FullName  string `json:"full_name"`
	RoleCount int    `json:"role_count"`
}

// GetMostActiveActor повертає актора з найбільшою кількістю ролей
func (h *AnalyticsHandler) GetMostActiveActor(c *gin.Context) {
	var result ActorRoleCountResult

	err := h.DB.Model(&models.Actor{}).
		Select("actors.id as actor_id, actors.full_name, count(roles.id) as role_count").
		Joins("join roles on actors.id = roles.actor_id").
		Group("actors.id, actors.full_name").
		Order("role_count desc").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No actors with roles found to analyze"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "A database error occurred", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetLeastActiveActor повертає актора з найменшою кількістю ролей
func (h *AnalyticsHandler) GetLeastActiveActor(c *gin.Context) {
	var result ActorRoleCountResult

	err := h.DB.Model(&models.Actor{}).
		Select("actors.id as actor_id, actors.full_name, count(roles.id) as role_count").
		Joins("left join roles on actors.id = roles.actor_id").
		Group("actors.id, actors.full_name").
		Order("role_count asc").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No actors found to analyze"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "A database error occurred", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PerformanceActorCountResult - Структура для повернення результату про виставу та кількість акторів
type PerformanceActorCountResult struct {
	PerformanceID uint   `json:"performance_id"`
	Title         string `json:"title"`
	ActorCount    int    `json:"actor_count"`
}

// GetPerformanceWithMostActors повертає виставу з найбільшою кількістю залучених акторів
func (h *AnalyticsHandler) GetPerformanceWithMostActors(c *gin.Context) {
	var result PerformanceActorCountResult

	err := h.DB.Model(&models.Performance{}).
		Select("performances.id as performance_id, performances.title, count(roles.id) as actor_count").
		Joins("join roles on performances.id = roles.performance_id").
		Group("performances.id, performances.title").
		Order("actor_count desc").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No performances with actors found to analyze"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "A database error occurred", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// PerformanceRehearsalCountResult - Структура для повернення результату про виставу та кількість репетицій
type PerformanceRehearsalCountResult struct {
	PerformanceID  uint   `json:"performance_id"`
	Title          string `json:"title"`
	RehearsalCount int    `json:"rehearsal_count"`
}

// GetMostFrequentPerformance повертає виставу, для якої заплановано найбільше репетицій
func (h *AnalyticsHandler) GetMostFrequentPerformance(c *gin.Context) {
	var result PerformanceRehearsalCountResult

	err := h.DB.Model(&models.Performance{}).
		Select("performances.id as performance_id, performances.title, count(rehearsals.id) as rehearsal_count").
		Joins("join rehearsals on performances.id = rehearsals.performance_id").
		Group("performances.id, performances.title").
		Order("rehearsal_count desc").
		First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No rehearsals found to analyze"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "A database error occurred", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

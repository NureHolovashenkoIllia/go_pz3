package handlers

import (
	"go_pz3/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActorHandler містить логіку для обробки запитів до акторів
type ActorHandler struct {
	DB *gorm.DB
}

func NewActorHandler(db *gorm.DB) *ActorHandler {
	return &ActorHandler{DB: db}
}

// CreateActor створює нового актора
func (h *ActorHandler) CreateActor(c *gin.Context) {
	var actor models.Actor
	if err := c.ShouldBindJSON(&actor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&actor); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, actor)
}

// GetActors повертає список усіх акторів
func (h *ActorHandler) GetActors(c *gin.Context) {
	var actors []models.Actor
	if result := h.DB.Find(&actors); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, actors)
}

// GetActor повертає одного актора за ID
func (h *ActorHandler) GetActor(c *gin.Context) {
	var actor models.Actor
	id := c.Param("id")

	if result := h.DB.First(&actor, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Actor not found"})
		return
	}

	c.JSON(http.StatusOK, actor)
}

// UpdateActor оновлює дані актора
func (h *ActorHandler) UpdateActor(c *gin.Context) {
	var actor models.Actor
	id := c.Param("id")

	if result := h.DB.First(&actor, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Actor not found"})
		return
	}

	if err := c.ShouldBindJSON(&actor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Save(&actor)
	c.JSON(http.StatusOK, actor)
}

// DeleteActor видаляє актора
func (h *ActorHandler) DeleteActor(c *gin.Context) {
	id := c.Param("id")
	if result := h.DB.Delete(&models.Actor{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Actor not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

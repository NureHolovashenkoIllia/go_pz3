package handlers

import (
	"go_pz3/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleHandler містить логіку для обробки запитів до ролей
type RoleHandler struct {
	DB *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

// CreateRole створює новий запис про роль (призначає актора на роль у виставі)
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Тут можна додати перевірки: чи існують такий actor_id та performance_id

	if result := h.DB.Create(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetRoles повертає список усіх ролей
func (h *RoleHandler) GetRoles(c *gin.Context) {
	var roles []models.Role
	if result := h.DB.Find(&roles); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetRole повертає одну роль за ID
func (h *RoleHandler) GetRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if result := h.DB.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// UpdateRole оновлює дані ролі
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var role models.Role
	id := c.Param("id")

	if result := h.DB.First(&role, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.DB.Save(&role)
	c.JSON(http.StatusOK, role)
}

// DeleteRole видаляє роль
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if result := h.DB.Delete(&models.Role{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

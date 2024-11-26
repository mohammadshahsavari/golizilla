package http

import (
	"golizilla/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RBACHandler struct {
	rbacService service.RBACService
}

func NewRBACHandler(rbacService service.RBACService) *RBACHandler {
	return &RBACHandler{rbacService}
}

func (h *RBACHandler) CreateRole(c *gin.Context) {
	var request struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.rbacService.CreateRole(request.Name, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// Additional endpoints like AssignPermissionToRole, AssignRoleToUser, etc.

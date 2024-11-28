package rbac

import (
	"github.com/gin-gonic/gin"
	"golizilla/handler/rbac"
)

func RegisterRoleRoutes(router *gin.Engine, roleHandler *rbac.RoleHandler) {
    r := router.Group("/roles")
    r.POST("/", roleHandler.CreateRole)
    r.GET("/", roleHandler.GetRoles)
}

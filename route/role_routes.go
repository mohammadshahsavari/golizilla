package route

import (
	"golizilla/handler/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.Engine, roleHandler *middleware.RoleHandler) {
	r := router.Group("/roles")
	r.POST("/", roleHandler.CreateRole)
	r.GET("/", roleHandler.GetRoles)
}

package routes

import (
	"modify-framework-mod/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, DB *gorm.DB) {

	r.Group("/api", utils.JWTAuthMiddleware(

	// add your routes here...
	// ex: api.POST(Method)("/endpoint", contollers.PackageName)
	))

}

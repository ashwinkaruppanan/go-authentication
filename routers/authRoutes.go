package routers

import (
	"ashwin.com/go-auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("/signup", controllers.SignUp())
	routes.POST("/signin", controllers.SignIn())
}

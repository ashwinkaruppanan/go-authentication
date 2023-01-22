package routers

import (
	"ashwin.com/go-auth/controllers"
	"ashwin.com/go-auth/middlewares"
	"github.com/gin-gonic/gin"
)

func ServiceRoutes(routes *gin.Engine) {
	routes.Use(middlewares.Authenticate())
	routes.GET("/", controllers.Home())
}

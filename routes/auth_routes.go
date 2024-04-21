package routes

import (
	"github.com/dastardlyjockey/ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	route.POST("/", controllers.SignUp)
	route.GET("/sign-in", controllers.SignIn)
}

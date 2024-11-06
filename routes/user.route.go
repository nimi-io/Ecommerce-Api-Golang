package routes

import (
	// "Ecommerce-Apis/controllers"

	con "Ecommerce-Api/controllers"
	mid "Ecommerce-Api/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incommingRoutes *gin.RouterGroup) {
	authGroup := incommingRoutes.Group("")

	authGroup.GET("/user", mid.AuthMiddleware(), con.GetUser)
	//  authGroup.POST("/login", controllers.SignIn)

}

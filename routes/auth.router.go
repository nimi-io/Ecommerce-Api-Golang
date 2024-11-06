package routes

import (
	// "Ecommerce-Apis/controllers"

	"Ecommerce-Api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incommingRoutes *gin.RouterGroup) {
	authGroup := incommingRoutes.Group("/auth")

	authGroup.POST("/register", controllers.Signup)
	authGroup.POST("/login", controllers.SignIn)

}

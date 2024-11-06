package routes

import (
	"Ecommerce-Api/controllers"
	mid "Ecommerce-Api/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(incomingRoutes *gin.RouterGroup) {
	adminGroup := incomingRoutes.Group("/admin", mid.AuthMiddleware(), mid.IsAdmin())
	productGroup := adminGroup.Group("/products")

	productGroup.POST("/create", controllers.CreateProduct)
	productGroup.GET("/:id", controllers.GetSingleProductById)
	productGroup.GET("", controllers.GetProducts)

	productGroup.PATCH("/:id", controllers.UpdateProduct)
	productGroup.DELETE("/:id", controllers.DeleteProduct)
}

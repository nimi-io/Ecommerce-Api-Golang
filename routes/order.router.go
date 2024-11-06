package routes

import (
	// "Ecommerce-Apis/controllers"

	"Ecommerce-Api/controllers"
	mid "Ecommerce-Api/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(incommingRoutes *gin.RouterGroup) {
	orderGroup := incommingRoutes.Group("/order", mid.AuthMiddleware())

	orderGroup.POST("/", controllers.PlaceOrder)
	orderGroup.GET("/list", controllers.ListMyOrders)
	// orderGroup.GET("/list/all", controllers.ListAllOrders)
	orderGroup.POST("cancel/:id", controllers.CancelOrder)

	orderGroup.PATCH("status/:id", mid.IsAdmin(), controllers.UpdateOrderStatus)

}

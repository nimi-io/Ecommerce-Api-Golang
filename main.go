package main

// @title E-Commerce API
// @version 1.0
// @description This is a sample e-commerce API.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	// routes "Auth-API/routes"
	db "Ecommerce-Api/database"
	_ "Ecommerce-Api/docs"
	"Ecommerce-Api/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"os"
	"time"
)

var startTime time.Time

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	startTime = time.Now()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db.DatabaseConnection(connStr)

	r := gin.New()
	r.Use(gin.Logger())

	apiVersion1Group := r.Group("/api/v1")

	routes.AuthRoutes(apiVersion1Group)
	routes.UserRoutes(apiVersion1Group)
	routes.AdminRoutes(apiVersion1Group)
	routes.OrderRoutes(apiVersion1Group)

	r.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	uptime := time.Since(startTime)

	r.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access granted for api-2",
			"message": "This is a message from api-2",
			"status":  "success",
			"code":    200,
			"error":   nil,
			"data":    nil,
			"uptime": gin.H{
				"seconds": uptime.Seconds(),
				"minutes": uptime.Minutes(),
				"hours":   uptime.Hours(),
			},
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + port)
}

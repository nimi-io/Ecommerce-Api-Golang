package controllers

import (
	db "Ecommerce-Api/database"
	"database/sql"
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	var fetchedUser db.User

	user, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	usr, ok := user.(map[string]interface{})
	if !ok {
		log.Println(usr)
		log.Println(reflect.TypeOf(usr))
		log.Println("User data is not of type jwt.MapClaims")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID, ok := usr["id"].(float64)
	if !ok {
		log.Println(reflect.TypeOf(usr["id"]))
		log.Println("User ID is not of expected type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err := db.DB.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", userID).
		Scan(&fetchedUser.ID, &fetchedUser.Username, &fetchedUser.Email, &fetchedUser.CreatedAt, &fetchedUser.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			log.Println("Database error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User found",
		"status":  http.StatusOK,
		"user":    fetchedUser,
	})
}

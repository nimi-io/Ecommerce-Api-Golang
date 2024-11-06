package controllers

import (
	models "Ecommerce-Api/database"
	help "Ecommerce-Api/helpers"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User to register"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func Signup(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if user.Email == "" || user.Password == "" || user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password and Username are required"})
		return
	}
	hashedPassword, err := help.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	checkQuery := "SELECT COUNT(*) FROM users WHERE email = $1"
	var count int
	err = models.DB.QueryRowContext(ctx, checkQuery, user.Email).Scan(&count)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	err = models.DB.QueryRowContext(ctx, query, user.Username, user.Email, hashedPassword).Scan(&user.ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login godoc
// @Summary Login a user
// @Description Login with email and password to get a JWT
// @Tags Users
// @Accept json
// @Produce json
// @Param login body map[string]string true "Email and Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func SignIn(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(user)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if user.Email == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both Email and password are required"})
		return
	}

	query := "SELECT id, username , password,created_at,role FROM users WHERE email = $1"
	var storedPassword string
	// var dbUser model.User
	//  err := model.DB.QueryRowContext(ctx, query, user.Email).Scan(&user.ID, &storedPassword,)

	err := models.DB.QueryRowContext(ctx, query, user.Email).Scan(&user.ID, &user.Username, &storedPassword, &user.CreatedAt, &user.Role)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	match := help.VerifyPassword(user.Password, storedPassword)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := help.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"role":       user.Role,
		},
		"token": token,
	})
}

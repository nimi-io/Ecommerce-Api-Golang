package controllers

import (
	db "Ecommerce-Api/database"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PlaceOrder(c *gin.Context) {
	user, exists := c.Get("user")
	log.Println("User data:", user)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	usr, ok := user.(map[string]interface{})
	if !ok {
		log.Println("User data is not of expected type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID, ok := usr["id"].(float64)
	if !ok {
		log.Println("User ID is not of expected type:", reflect.TypeOf(usr["id"]))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var requestBody struct {
		Products []struct {
			ProductID uint `json:"product_id"`
			Quantity  int  `json:"quantity"`
		} `json:"products"`
	}
	log.Println("Received request body:", requestBody)
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	order := db.Order{
		UserID:    uint(userID),
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO orders (user_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	err := db.DB.QueryRow(query, order.UserID, order.Status, order.CreatedAt, order.UpdatedAt).Scan(&order.ID)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	for _, product := range requestBody.Products {
		_, err := db.DB.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
		`, order.ID, product.ProductID, product.Quantity, time.Now(), time.Now())

		if err != nil {
			log.Println("Database error while inserting order items:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Order placed successfully",
		"status":  http.StatusCreated,
		"order":   order,
	})
}

func ListMyOrders(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	usr, ok := user.(map[string]interface{})
	if !ok {
		log.Println("User data is not of expected type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID, ok := usr["id"].(float64)
	if !ok {
		log.Println("User ID is not of expected type:", reflect.TypeOf(usr["id"]))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var orders []db.Order
	query := `
		SELECT o.id, o.user_id, o.status, o.created_at, o.updated_at
		FROM orders o
		WHERE o.user_id = $1
		ORDER BY o.created_at DESC
		LIMIT $2 OFFSET $3
	`
	orderRows, err := db.DB.Query(query, uint(userID), limit, offset)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer orderRows.Close()

	for orderRows.Next() {
		var order db.Order
		err := orderRows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			log.Println("Database error while scanning orders:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Retrieve the order items for each order
		var orderItems []db.OrderProduct
		orderItemsQuery := `
			SELECT oi.product_id, oi.quantity, p.id, p.name, p.description, p.price, p.stock, p.created_at, p.updated_at
			FROM order_items oi
			INNER JOIN products p ON oi.product_id = p.id
			WHERE oi.order_id = $1
		`
		itemRows, err := db.DB.Query(orderItemsQuery, order.ID)
		if err != nil {
			log.Println("Database error while retrieving order items:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer itemRows.Close()

		// Collect all order items with product details
		for itemRows.Next() {
			var orderItem db.OrderProduct
			var product db.Product
			err := itemRows.Scan(&orderItem.ProductID, &orderItem.Quantity, &product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
			if err != nil {
				log.Println("Database error while scanning order items:", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			orderItem.Product = product
			orderItems = append(orderItems, orderItem)
		}

		// Assign order items to the order
		order.Products = orderItems
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Orders retrieved successfully",
		"status":  http.StatusOK,
		"orders":  orders,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
		},
	},
	)
}

func CancelOrder(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	usr, ok := user.(map[string]interface{})
	if !ok {
		log.Println("User data is not of expected type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID, ok := usr["id"].(float64)
	if !ok {
		log.Println("User ID is not of expected type:", reflect.TypeOf(usr["id"]))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	orderID := c.Param("id")

	var order db.Order
	query := `SELECT id, user_id, status FROM orders WHERE id = $1 AND user_id = $2`
	err := db.DB.QueryRow(query, orderID, uint(userID)).Scan(&order.ID, &order.UserID, &order.Status)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already cancelled"})
		return
	}

	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order cannot be cancelled"})
		return
	}

	_, err = db.DB.Exec(`UPDATE orders SET status = 'cancelled' WHERE id = $1`, order.ID)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"status":  http.StatusOK,
		"order":   order,
	})
}

// Update the status
func UpdateOrderStatus(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	usr, ok := user.(map[string]interface{})
	if !ok {
		log.Println("User data is not of expected type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	userID, ok := usr["id"].(float64)
	if !ok {
		log.Println("User ID is not of expected type:", reflect.TypeOf(usr["id"]))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	orderID := c.Param("id")

	var order db.Order
	query := `SELECT id, user_id, status FROM orders WHERE id = $1 AND user_id = $2`
	err := db.DB.QueryRow(query, orderID, uint(userID)).Scan(&order.ID, &order.UserID, &order.Status)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == "delivered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already delivered"})
		return
	}

	_, err = db.DB.Exec(`UPDATE orders SET status = 'delivered' WHERE id = $1`, order.ID)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"status":  http.StatusOK,
		"order":   order,
	})
}

package controllers

import (
	db "Ecommerce-Api/database"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	var products []db.Product

	// Pagination parameters
	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Query with limit and offset for pagination
	rows, err := db.DB.Query("SELECT id, name, description, price, stock, created_at, updated_at FROM products LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product db.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			log.Println("Database error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Products found",
		"status":   http.StatusOK,
		"products": products,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": len(products),
		},
	})
}
func GetSingleProductById(c *gin.Context) {
	productId := c.Param("id")
	var product db.Product

	err := db.DB.QueryRow("SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1", productId).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product found",
		"status":  http.StatusOK,
		"product": product,
	})
}

func CreateProduct(c *gin.Context) {
	var product db.Product

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	// check if name exist
	var existingProduct db.Product
	if err := db.DB.QueryRow("SELECT id FROM products WHERE name = $1", product.Name).Scan(&existingProduct.ID); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product with the same name already exists"})
		return
	}

	query := `INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id, name, created_at, updated_at`

	err := db.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Stock).Scan(&product.ID, &product.Name, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"status":  http.StatusCreated,
		"product": product,
	})
}

func UpdateProduct(c *gin.Context) {
	var product db.Product

	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := `UPDATE products SET name = $1, description = $2, price = $3, quantity = $4 WHERE id = $5 RETURNING updated_at`

	err := db.DB.QueryRow(query, product.Name, product.Description, product.Price, product.Stock, product.ID).Scan(&product.UpdatedAt)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"status":  http.StatusOK,
		"product": product,
	})

}

func DeleteProduct(c *gin.Context) {
	productId := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM products WHERE id = $1", productId)
	if err != nil {
		log.Println("Database error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
		"status":  http.StatusOK,
	})

}

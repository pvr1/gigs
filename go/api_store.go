package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteOrder - Delete purchase order by ID
func DeleteOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// GetInventory - Returns gig inventories by status
func GetInventory(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gigs)
}

// GetOrderById - Find purchase order by ID
func GetOrderById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PlaceOrder - Place an order for a gig
func PlaceOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

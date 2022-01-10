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
	_, err := c.GetQuery("RegisteredClaims")
	if err {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid claims"})
		return
	}

	c.JSON(http.StatusOK, gigs)
}

// GetOrderById - Find purchase order by ID
func GetOrderById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// PlaceOrder - Place an order for a gig
func PlaceOrder(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

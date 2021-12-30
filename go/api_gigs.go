package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UpdateGig - Update an existing gig
func UpdateGig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

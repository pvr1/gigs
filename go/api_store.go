package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

// RemoveTransaction - Helper function to remove a transaction from the slice
func RemoveTransaction(s []transaction, index int) []transaction {
	return append(s[:index], s[index+1:]...)
}

// DeleteTransaction - Delete purchase transaction by ID
func DeleteTransaction(c *gin.Context) {
	id := c.Param("transactionId")
	// Loop over the list of transactions, looking for
	// an transaction whose ID value matches the parameter.
	for i, a := range transactions {
		if a.Id == id {
			transactions = RemoveTransaction(transactions, i)
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
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

// GetTransactionById - Find purchase transaction by ID
func GetTransactionById(c *gin.Context) {
	id := c.Param("transactionId")
	// Loop over the list of transactions, looking for
	// an transaction whose ID value matches the parameter.

	for _, a := range transactions {
		if a.Id == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "transaction not found"})
}

// PlaceTransaction - Place an transaction for a gig
func PlaceTransaction(c *gin.Context) {
	var mypurchase transaction

	// Call BindJSON to bind the received JSON to
	// mypurchase.
	mypurchase.Id = uuid.NewV4().String()
	if err := c.BindJSON(&mypurchase); err != nil {
		return
	}

	// Add the new gig to the slice.
	transactions = append(transactions, mypurchase)
	c.JSON(http.StatusCreated, mypurchase)
}

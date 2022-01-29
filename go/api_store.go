package openapi

/* Example on sanitizer
import (
	"html/template"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func postHandler(c *gin.Context) {
	id := c.Param("id")
	var post Post
	db.Where("id = ?", id).First(&post)

	unsafe := blackfriday.MarkdownCommon([]byte(post.Body))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	c.HTML(http.StatusOK, "post.html", pongo2.Context{
		"Post":     post,
		"Markdown": template.HTML(html),
	})
}
*/

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/twinj/uuid"
)

// RemoveTransaction - Helper function to remove a transaction from the slice
func RemoveTransaction(s []transaction, index int) []transaction {
	return append(s[:index], s[index+1:]...)
}

// DeleteTransaction - Delete purchase transaction by ID
func DeleteTransaction(c *gin.Context) {
	id := c.Param("transactionId")
	//First try at blackfriday and bluemonday
	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))
	// Loop over the list of transactions, looking for
	// an transaction whose ID value matches the parameter.
	for i, a := range transactions {
		if a.Id == html {
			transactions = RemoveTransaction(transactions, i)
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// GetTransactions - Returns gig inventories by status
func GetTransactions(c *gin.Context) {
	status, err := c.GetQuery("RegisteredClaims")
	if err {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid claims"})
		return
	}

	unsafe := blackfriday.SanitizedAnchorName(status)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	tmp := []transaction{}

	for _, a := range transactions {
		if a.Status == html {
			tmp = append(tmp, a)
		}
	}

	c.JSON(http.StatusOK, tmp)
}

// GetTransactionById - Find purchase transaction by ID
func GetTransactionById(c *gin.Context) {
	id := c.Param("transactionId")
	// Loop over the list of transactions, looking for
	// an transaction whose ID value matches the parameter.

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	for _, a := range transactions {
		if a.Id == html {
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

	today := time.Now()
	//Validate the transaction
	if mypurchase.GigId == "" || mypurchase.Status == "" ||
		mypurchase.Price == 0 || mypurchase.Complete ||
		today.After(mypurchase.ShipDate) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction"})
		return
	}

	// Add the new gig to the slice.
	transactions = append(transactions, mypurchase)
	c.JSON(http.StatusCreated, mypurchase)
}

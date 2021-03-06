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
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

			//Delete record in mongodb
			client, ctx := connectMongoDB()
			defer client.Disconnect(ctx)
			tranactionsCollection := getCollectionMongoDB(client, "transactions")
			transE, err := tranactionsCollection.DeleteOne(ctx,
				bson.M{"id": html},
			)
			if err != nil {
				println("add record error")
				log.Fatal(transE, err)
			}

			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// GetTransactions - Returns gig inventories by status
func GetTransactions(c *gin.Context) {
	resyncTransactions()
	/*
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
	*/
	c.JSON(http.StatusOK, transactions)
}

// GetTransactionById - Find purchase transaction by ID
func GetTransactionById(c *gin.Context) {
	resyncTransactions()
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

	//Insert record into mongodb
	client, ctx := connectMongoDB()
	defer client.Disconnect(ctx)
	transCollection := getCollectionMongoDB(client, "transactions")
	transE, err := transCollection.InsertMany(ctx, []interface{}{
		&mypurchase,
	})
	if err != nil {
		println("add record error")
		log.Fatal(transE, err)
	}

	c.JSON(http.StatusCreated, mypurchase)
}

// AcceptTransaction - Accept a transaction (as a employer) that has been placed by a worker
func AcceptTransaction(c *gin.Context) {
	id := c.Param("transactionId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	for _, a := range transactions {
		if a.Id == html {
			a.Status = "Accepted"
			//TODO: Send info about the transaction to the worker
			client, ctx := connectMongoDB()
			defer client.Disconnect(ctx)
			transCollection := getCollectionMongoDB(client, "transactions")
			transE, err := transCollection.UpdateOne(ctx,
				bson.M{"id": html},
				bson.M{"$set": bson.M{"status": "Accepted"}},
			)
			if err != nil {
				println("add record error")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction"})
				return
			}

			c.JSON(http.StatusOK, transE)
			return
		}
	}

}

// AcceptTransaction - Accept a transaction (as a employer) that has been placed by a worker
func RejectTransaction(c *gin.Context) {
	id := c.Param("transactionId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	for _, a := range transactions {
		if a.Id == html {
			a.Status = "Rejected"
			//TODO: Send info about the transaction to the worker
			client, ctx := connectMongoDB()
			defer client.Disconnect(ctx)
			transCollection := getCollectionMongoDB(client, "transactions")
			transE, err := transCollection.UpdateOne(ctx,
				bson.M{"id": html},
				bson.M{"$set": bson.M{"status": "Rejected"}},
			)
			if err != nil {
				println("add record error")
				c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid transaction"})
				return
			}

			c.JSON(http.StatusOK, transE)
			return
		}
	}

}

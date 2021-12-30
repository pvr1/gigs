/*
 * Gigs workforce api
 *
 * This is the Gigs workforce api. Use to search for workforce or search for a gig
 *
 * API version: 1.0.0
 * Contact: per.von.rosen@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddGig - Add a new gig to the store
func AddGig(c *gin.Context) {
	var mygig Gig

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&mygig); err != nil {
		return
	}

	// Add the new album to the slice.
	gigs = append(gigs, mygig)
	c.IndentedJSON(http.StatusCreated, mygig)
}

// DeleteGig - Deletes a gig
func DeleteGig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// FindGigsByStatus - Finds Gigs by status
func FindGigsByStatus(c *gin.Context) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println(c.GetQuery("status"))
	status, _ := c.GetQuery("status")
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	log.Println(status)

	tmp := []Gig{}

	for _, a := range gigs {
		if a.Status == status {
			tmp = append(tmp, a)
		}
	}

	c.IndentedJSON(http.StatusOK, tmp)
}

// GetGigById - Find gig by ID
func GetGigById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("gigId"), 10, 0)
	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.

	for _, a := range gigs {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

}

// UpdateGigWithForm - Updates a gig in the store with form data
func UpdateGigWithForm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// UploadFile - uploads an image
func UploadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

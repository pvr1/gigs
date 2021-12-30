package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

// AddGig - Add a new gig to the store
func AddGig(c *gin.Context) {
	var mygig Gig

	// Call BindJSON to bind the received JSON to
	// newgig.
	mygig.Id = uuid.NewV4().String()
	if err := c.BindJSON(&mygig); err != nil {
		return
	}

	// Add the new gig to the slice.
	gigs = append(gigs, mygig)
	c.IndentedJSON(http.StatusCreated, mygig)
}

// RemoveIndex - Helper function to remove a gig from the slice
func RemoveIndex(s []Gig, index int) []Gig {
	return append(s[:index], s[index+1:]...)
}

// DeleteGig - Deletes a gig
func DeleteGig(c *gin.Context) {
	id := c.Param("gigId")
	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == id {
			gigs = RemoveIndex(gigs, i)
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gig not found"})
}

// FindGigsByStatus - Finds Gigs by status
func FindGigsByStatus(c *gin.Context) {
	status, _ := c.GetQuery("status")
	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

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
	id := c.Param("gigId")
	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	for _, a := range gigs {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gig not found"})

}

// UpdateGigWithForm - Updates a gig in the store with form data
func UpdateGigWithForm(c *gin.Context) {
	var mygig Gig
	id := c.Param("gigId")

	// Call BindJSON to bind the received JSON to
	// newgig.
	mygig.Id = uuid.NewV4().String()
	if err := c.BindJSON(&mygig); err != nil {
		// Loop over the list of gigs, looking for
		// an gig whose ID value matches the parameter.
		for i, a := range gigs {
			if a.Id == id {
				gigs[i] = mygig
				c.IndentedJSON(http.StatusOK, a)
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gig not found"})
}

// UploadFile - uploads an image
func UploadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

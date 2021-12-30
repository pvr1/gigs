package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

// UpdateGig - Update an existing gig
func UpdateGig(c *gin.Context) {
	var mygig Gig
	id := c.Param("gigId")

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	mygig.Id = uuid.NewV4().String()
	if err := c.BindJSON(&mygig); err != nil {
		// Loop over the list of albums, looking for
		// an album whose ID value matches the parameter.
		for i, a := range gigs {
			if a.Id == id {
				gigs[i] = mygig
				c.IndentedJSON(http.StatusOK, a)
			}
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "gig not found"})
}

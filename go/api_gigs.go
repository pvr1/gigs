package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeepCopy deepcopies a to b using json marshaling
func DeepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}

// UpdateGig - Update an existing gig
func UpdateGig(c *gin.Context) {
	/*
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	bodyjson, err := c.GetRawData()
	var body Gig
	json.Unmarshal(bodyjson, &body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == body.Id {
			// Update the gig
			DeepCopy(body, &gigs[i])
			c.IndentedJSON(http.StatusOK, body)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

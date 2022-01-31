package openapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
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

	//Validadte the gig
	if mygig.Name == "" || mygig.Description == nil || mygig.Status == "" || mygig.Measurableoutcome == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing required gig fields"})
		return
	}

	// Add the new gig to the slice.
	gigs = append(gigs, mygig)
	c.JSON(http.StatusCreated, mygig)
}

// RemoveGig - Helper function to remove a gig from the slice
func RemoveGig(s []Gig, index int) []Gig {
	return append(s[:index], s[index+1:]...)
}

// DeleteGig - Deletes a gig
func DeleteGig(c *gin.Context) {
	id := c.Param("gigId")

	/*
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	//First try at blackfriday and bluemonday
	//log.Println("id: ", id)
	unsafe := blackfriday.SanitizedAnchorName(id)
	//log.Println("unsafe: ", unsafe)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))
	//log.Println("html: ", html)

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == html {
			gigs = RemoveGig(gigs, i)
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// FindGigsByStatus - Finds Gigs by status
func FindGigsByStatus(c *gin.Context) {
	status, err := c.GetQuery("status")
	if !err {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	unsafe := blackfriday.SanitizedAnchorName(status)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	tmp := []Gig{}

	for _, a := range gigs {
		if a.Status == html {
			tmp = append(tmp, a)
		}
	}

	c.JSON(http.StatusOK, tmp)
}

// GetGigById - Find gig by ID
func GetGigById(c *gin.Context) {
	id := c.Param("gigId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	for _, a := range gigs {
		if a.Id == html {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// UpdateGigWithForm - Updates a gig in the store with form data
func UpdateGigWithForm(c *gin.Context) {
	// If the file doesn't exist, create it or append to the file
	/*
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	*/
	id := c.Param("gigId")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	bodyjson, err := c.GetRawData()
	var body Gig
	json.Unmarshal(bodyjson, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if body.Id != id {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id mismatch"})
		return
	}

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range gigs {
		if a.Id == html {
			// Update the gig
			DeepCopy(body, &gigs[i])
			c.JSON(http.StatusOK, body)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Gig not found"})
}

// UploadFile - uploads an image
func UploadFile(c *gin.Context) {
	// Single file
	file, err := c.FormFile("file")
	log.Println("file: ", file)
	if err != nil {
		log.Println("Error uploading file - ", err)
		c.String(http.StatusBadRequest, fmt.Sprintf("Error uploading: %s", err.Error()))
		return
	}

	// Retrieve file information
	extension := filepath.Ext(file.Filename)
	// Generate random file name for the new uploaded file so it doesn't override the old file with same name
	newFileName := uuid.NewV4().String() + extension

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, "/tmp/"+newFileName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}

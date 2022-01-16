package openapi

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/twinj/uuid"
)

//TODO: Add OIDC support

// CreateUser - Create user
func CreateUser(c *gin.Context) {
	var tmpUser User

	// Call BindJSON to bind the received JSON to
	// newgig.
	if err := c.BindJSON(&tmpUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}
	tmpUser.Id = uuid.NewV4().String()

	//Validate the user
	if tmpUser.Username == "" || tmpUser.Password == "" || tmpUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": tmpUser.Username + " " + tmpUser.Password + " " + tmpUser.Email})
		return
	}

	// Add the new user to the slice.
	users = append(users, tmpUser)
	c.JSON(http.StatusCreated, tmpUser)
}

// CreateUsersWithArrayInput - Creates list of users with given input array
func CreateUsersWithArrayInput(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// CreateUsersWithListInput - Creates list of users with given input array
func CreateUsersWithListInput(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// RemoveUser - Helper function to remove a user from the slice
func RemoveUser(s []User, index int) []User {
	return append(s[:index], s[index+1:]...)
}

// DeleteUser - Delete user
func DeleteUser(c *gin.Context) {
	id := c.Param("username")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range users {
		if a.Username == html {
			users = RemoveUser(users, i)
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// GetUserByName - Get user by user name
func GetUserByName(c *gin.Context) {
	id := c.Param("username")

	unsafe := blackfriday.SanitizedAnchorName(id)
	html := string(bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe)))

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.

	for _, a := range users {
		if a.Username == html {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

// LoginUser - Logs user into the system
func LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// LogoutUser - Logs out current logged in user session
func LogoutUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// UpdateUser - Updated user
func UpdateUser(c *gin.Context) {
	/*
		// If the file doesn't exist, create it or append to the file
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	*/

	bodyjson, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	var body User

	//unsafe := blackfriday2.Run(bodyjson)
	//html := bluemonday.UGCPolicy().SanitizeBytes([]byte(unsafe))

	json.Unmarshal(bodyjson, &body)

	//Validate the user
	if body.Username == "" || body.Password == "" || body.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed user data"})
		//c.JSON(http.StatusBadRequest, gin.H{"message": "bodyjson: " + string(bodyjson) + "\n unsafe: " + string(unsafe) + "\n html: " + string(html) + " " + body.Username + " " + body.Password + " " + body.Email})
		//c.JSON(http.StatusBadRequest, gin.H{"message": "unsafe: " + string(unsafe)})
		//c.JSON(http.StatusBadRequest, gin.H{"message": "html: " + string(html)})
		//c.JSON(http.StatusBadRequest, gin.H{"message": "Username: " + body.Username + "Password: " + body.Password + "Email: " + body.Email})
		return
	}

	// Loop over the list of gigs, looking for
	// an gig whose ID value matches the parameter.
	for i, a := range users {
		if a.Username == body.Username {
			// Update the gig
			DeepCopy(body, &users[i])
			c.JSON(http.StatusOK, body)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

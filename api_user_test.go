package main

import (
	"bytes"
	"net/http"
	"testing"

	openapi "github.com/pvr1/gigs/go"
	"github.com/stretchr/testify/assert"
)

func TestAdduser(t *testing.T) {
	/*
		body := gin.H{
			"Username":            "Kalle",
			"FirstName":           "Carl",
			"LastName":            "Karlsson",
			"Email":               "carl@Karlsson.se",
			"SocialSecuityNumber": "7001016939",
			"Phone":               "+462120000",
		}
	*/

	body := bytes.NewBufferString("Username=Kalle&FirstName=Carl&LastName=Piper&Email=carl@piper.se&SocialSecuityNumber=7001016939&Phone=+462120000")

	router := openapi.NewRouter()
	w := performRequest(router, "POST", "/v2/user", body)
	assert.Equal(t, http.StatusOK, w.Code)

	/*
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["hello"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, "hello", value)
	*/

	a := w.Body.String()
	assert.Equal(t, "user added\n", a)
}

func TestGetuser(t *testing.T) {
	/*
		body := gin.H{
			"userID": "1",
		}
	*/

	body := bytes.NewBufferString("userID=1")
	router := openapi.NewRouter()
	w := performRequest(router, "GET", "/v2/user/1", body)
	assert.Equal(t, http.StatusOK, w.Code)

	/*
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["hello"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, body["hello"], value)
	*/

	a := w.Body.String()
	assert.Equal(t, "There you got your specific user\n", a)
}

func TestGetusers(t *testing.T) {
	/*
		body := gin.H{
			"": "",
		}
	*/
	router := openapi.NewRouter()
	w := performRequest(router, "GET", "/v2/user", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	a := w.Body.String()
	assert.Equal(t, "Yep. A list of user was delivered. Can you see it?? :-)\n", a)
}

func TestUpdateuser(t *testing.T) {
	body := bytes.NewBufferString("userID=1")
	router := openapi.NewRouter()
	w := performRequest(router, "PUT", "/v2/user", body)
	assert.Equal(t, http.StatusOK, w.Code)

	/*
		var response map[string]stringThere you got your specific transaction
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["hello"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, body["hello"], value)
	*/

	a := w.Body.String()
	assert.Equal(t, "user updated. Now it belongs to me.\n", a)
}

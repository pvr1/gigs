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
					"FirstName":  "Kalle",
					"LastName":   "Carlsson",
					"Password":   "pwd",
					"UserStatus": 0,
					"Phone":      "888-888-8888",
					"Id":         "13",
					"Email":      "a.a@a.com",
					"Username":   "kallekula",
		}
	*/

	str := "{\"Id\":\"13\",\"Username\":\"kallekula\",\"FirstName\":\"Kalle\",\"LastName\":\"Kula\",\"Email\":\"a.a@a.com\",\"Password\":\"pwd\",\"Phone\":\"888-888-8888\",\"UserStatus\":1,\"Role\":null}"
	body := bytes.NewBufferString(str)

	router := openapi.NewTestRouter()
	w := performRequest(router, "POST", "/v2/user", body)
	assert.Equal(t, http.StatusCreated, w.Code)

	/*
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["hello"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, "hello", value)
	*/

	//a := w.Body.String()
	//assert.Equal(t, "user added\n", a)
}

func TestGetuser(t *testing.T) {
	/*
		body := gin.H{
			"userID": "1",
		}
	*/

	body := bytes.NewBufferString("")
	router := openapi.NewTestRouter()
	w := performRequest(router, "GET", "/v2/user/aaa", body)
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
	assert.Equal(t, "{\"Id\":\"0\",\"Username\":\"aaa\",\"FirstName\":\"firstName\",\"LastName\":\"lastName\",\"Email\":\"a.a@a.com\",\"Password\":\"password\",\"Phone\":\"888-888-8888\",\"UserStatus\":6,\"Role\":null}", a)
}

func TestGetusers(t *testing.T) {
	/*
		body := gin.H{
			"": "",
		}
	*/
	/*
		auth, err := authenticator.New()
		if err != nil {
			t.Errorf("Failed to initialize the authenticator: %v", err)
		}

		router := openapi.NewRouter(auth)
	*/
	router := openapi.NewRouter()
	w := performRequest(router, "GET", "/v2/user/aaa", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	a := w.Body.String()
	assert.Equal(t, "{\"Id\":\"0\",\"Username\":\"aaa\",\"FirstName\":\"firstName\",\"LastName\":\"lastName\",\"Email\":\"a.a@a.com\",\"Password\":\"password\",\"Phone\":\"888-888-8888\",\"UserStatus\":6,\"Role\":null}", a)
}

func TestUpdateuser(t *testing.T) {
	body := bytes.NewBufferString("{\"Id\":\"0\",\"Username\":\"aaa\",\"FirstName\":\"Lasse\",\"LastName\":\"lastName\",\"Email\":\"a.a@a.com\",\"Password\":\"password\",\"Phone\":\"888-888-8888\",\"UserStatus\":6,\"Role\":null}")
	router := openapi.NewTestRouter()
	w := performRequest(router, "PUT", "/v2/user/aaa", body)
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
	assert.Equal(t, "{\"Id\":\"0\",\"Username\":\"aaa\",\"FirstName\":\"Lasse\",\"LastName\":\"lastName\",\"Email\":\"a.a@a.com\",\"Password\":\"password\",\"Phone\":\"888-888-8888\",\"UserStatus\":6,\"Role\":null}", a)
}

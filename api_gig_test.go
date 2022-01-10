package main

import (
	"bytes"
	"net/http"
	"testing"

	openapi "github.com/pvr1/gigs/go"
	"github.com/stretchr/testify/assert"
)

func TestAddRemoveGig(t *testing.T) {
	/*
		{\"Id\":\"1\",\"Category\":{\"Id\":0,\"Name\":\"\"},\"Name\":\"Hepp\",\"Description\":[\"descr\"],\"Measurableoutcome\":[\"hahaha\"],\"Tags\":null,\"Status\":\"available\"}"
	*/

	body := bytes.NewBufferString("{\"Id\":\"1\",\"Category\":{\"Id\":0,\"Name\":\"\"},\"Name\":\"Hepp\",\"Description\":[\"descr\"],\"Measurableoutcome\":[\"hahaha\"],\"Tags\":null,\"Status\":\"available\"}")
	router := openapi.NewTestRouter()
	w := performRequest(router, "POST", "/v2/gigs", body)
	assert.Equal(t, http.StatusCreated, w.Code)

	/*
		var response map[string]string
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		value, exists := response["hello"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, "hello", value)
	*/

	a := w.Body.String()
	assert.Equal(t, "{\"Id\":\"1\",\"Category\":{\"Id\":0,\"Name\":\"\"},\"Name\":\"Hepp\",\"Description\":[\"descr\"],\"Measurableoutcome\":[\"hahaha\"],\"Tags\":null,\"Status\":\"available\"}", a)

	w = performRequest(router, "DELETE", "/v2/gig/1", body)
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestGetGig(t *testing.T) {
	/*
		body := gin.H{
			"userID": "1",
		}
	*/

	body := bytes.NewBufferString("userID=1")
	router := openapi.NewTestRouter()
	w := performRequest(router, "GET", "/v2/gig/1", body)
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
	assert.Equal(t, "{\"Id\":\"1\",\"Category\":{\"Id\":0,\"Name\":\"\"},\"Name\":\"Hepp\",\"Description\":[\"descr\"],\"Measurableoutcome\":[\"hahaha\"],\"Tags\":null,\"Status\":\"available\"}", a)
}

func TestGetGigs(t *testing.T) {
	/*
		body := gin.H{
			"": "",
		}
	*/
	router := openapi.NewTestRouter()
	w := performRequest(router, "GET", "/v2/store/inventory", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	//a := w.Body.String()
	//assert.Equal(t, "Yep. A list of gig was delivered. Can you see it?? :-)\n", a)
}

func TestUpdategig(t *testing.T) {
	body := bytes.NewBufferString("{\"id\": \"1\",\"category\":{},\"name\":\"Prutt\",\"description\":[\"descr\"],\"measurableoutcome\":[\"hahaha\"],\"status\":\"available\"}")
	router := openapi.NewTestRouter()
	w := performRequest(router, "PUT", "/v2/gigs", body)
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
	assert.Equal(t, "{\"Id\":\"1\",\"Category\":{\"Id\":0,\"Name\":\"\"},\"Name\":\"Prutt\",\"Description\":[\"descr\"],\"Measurableoutcome\":[\"hahaha\"],\"Tags\":null,\"Status\":\"available\"}", a)
}

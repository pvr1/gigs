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
		{"Id": "q","Name": "Fix","Description": ["jobba"],"Measurableoutcome": ["pengar"],"Tags": null,"Status": "available"}
	*/

	body := bytes.NewBufferString("{\n    \"Id\": \"q\",\n    \"Category\": {\n        \"Id\": 0,\n        \"Name\": \"\"\n    },\n    \"Name\": \"Fix\",\n    \"Description\": [\n        \"jobba\"\n    ],\n    \"Measurableoutcome\": [\n        \"pengar\"\n    ],\n    \"Tags\": null,\n    \"Status\": \"available\"\n}")
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
	assert.Equal(t, "{\n    \"Id\": \"q\",\n    \"Category\": {\n        \"Id\": 0,\n        \"Name\": \"\"\n    },\n    \"Name\": \"Fix\",\n    \"Description\": [\n        \"jobba\"\n    ],\n    \"Measurableoutcome\": [\n        \"pengar\"\n    ],\n    \"Tags\": null,\n    \"Status\": \"available\"\n}", a)

	w = performRequest(router, "DELETE", "/v2/gig/q", body)
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
	assert.Equal(t, "{\n    \"Id\": \"1\",\n    \"Category\": {\n        \"Id\": 0,\n        \"Name\": \"\"\n    },\n    \"Name\": \"Gig 1\",\n    \"Description\": [\n        \"description 1\"\n    ],\n    \"Measurableoutcome\": [\n        \"measurableoutcome 1\"\n    ],\n    \"Tags\": null,\n    \"Status\": \"available\"\n}", a)
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
	assert.Equal(t, "{\n    \"Id\": \"1\",\n    \"Category\": {\n        \"Id\": 0,\n        \"Name\": \"\"\n    },\n    \"Name\": \"Prutt\",\n    \"Description\": [\n        \"descr\"\n    ],\n    \"Measurableoutcome\": [\n        \"hahaha\"\n    ],\n    \"Tags\": null,\n    \"Status\": \"available\"\n}", a)
}

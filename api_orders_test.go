package main

import (
	"bytes"
	"net/http"
	"testing"

	openapi "github.com/pvr1/gigs/go"
	"github.com/stretchr/testify/assert"
)

func _TestGetTransaction(t *testing.T) {
	body := bytes.NewBufferString("userID=1")
	router := openapi.NewRouter()
	w := performRequest(router, "GET", "/v2/store/order/1", body)
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
	assert.Equal(t, "There you got your specific transaction\n", a)
}

func _TestGetTransactions(t *testing.T) {
	body := bytes.NewBufferString("userID=1")
	router := openapi.NewRouter()
	w := performRequest(router, "GET", "/v2/store/order/", body)
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
	assert.Equal(t, "Get list of Transactions originating from FX Orders and Payments\n", a)
}

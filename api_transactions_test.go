package main

import (
	"bytes"
	"net/http"
	"testing"

	openapi "github.com/pvr1/gigs/go"
	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	body := bytes.NewBufferString("")
	router := openapi.NewTestRouter()
	w := performRequest(router, "GET", "/v2/store/transaction/1", body)
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
	assert.Equal(t, "{\"Id\":\"1\",\"GigId\":\"1\",\"Price\":100,\"ShipDate\":\"2012-11-01T22:08:41Z\",\"Status\":\"pending\",\"Complete\":false}", a)
}

func _TestGetTransactions(t *testing.T) {
	body := bytes.NewBufferString("")
	router := openapi.NewTestRouter()
	w := performRequest(router, "GET", "/v2/store/transaction/", body)
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
	assert.Equal(t, "Get list of Transactions originating from FX transactions and Payments\n", a)
}

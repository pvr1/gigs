package main

import (
	"testing"

	"github.com/gin-gonic/gin"
	openapi "github.com/pvr1/gigs/go"
)

func TestOpenAPIjson(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openapi.OpenAPIjson(tt.args.c)
		})
	}
}

func TestOpenAPIyaml(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			openapi.OpenAPIyaml(tt.args.c)
		})
	}
}

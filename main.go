/*
 * Gigs workforce api
 *
 * This is the Gigs workforce api. Use to search for workforce or search for a gig
 *
 * API version: 1.0.0
 * Contact: per.von.rosen@gmail.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"

	// WARNING!
	// Change this to a fully-qualified import path
	// once you place this file into your project.
	// For example,
	//
	//sw "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	//
	sw "./go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(router.Run(":8080"))
}

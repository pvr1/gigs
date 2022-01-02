package main

import (
	"log"
	"github.com/pvr1/gigs/gin-oidc/gin_oidc"
	sw "github.com/pvr1/gigs/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

  	//middleware params
 	initParams := gin_oidc.InitParams{
 		Router:       router,
 		ClientId:     "xx-xxx-xxx",
 		ClientSecret: "xx-xxx-xxx",
 		Issuer:       "https://accounts.google.com/", //add '.well-known/openid-configuration' to see it's a good link
 		ClientUrl:    "http://example.domain/", //your website's url
 		Scopes:       ["openid"],
 		ErrorHandler: func(c *gin.Context) {
 			//gin_oidc pushes a new error before any "ErrorHandler" invocation
 			message := c.Errors.Last().Error()
 			//redirect to ErrorEndpoint with error message
 			redirectToErrorPage(c, "http://example2.domain/error", message)
 			//when "ErrorHandler" ends "c.Abort()" is invoked - no further handlers will be invoked
 		},
 		PostLogoutUrl: "http://example2.domain/",
 	}
  
 	//protect all endpoint below this line
 	router.Use(gin_oidc.Init(initParams))

	log.Fatal(router.Run(":8080"))
}

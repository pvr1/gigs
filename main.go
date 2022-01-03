package main

import (
	"github.com/gin-gonic/gin"
	sw "github.com/pvr1/gigs/go"
	"github.com/pvr1/gin-oidc"
	"log"
	"net/url"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	issuer, _ := url.Parse("https://accounts.google.com")
	clienturl, _ := url.Parse("http://localhost:8080/login")
	postlogouturl, _ := url.Parse("http://localhost:8080/logout")

	//middleware params
	initParams := gin_oidc.InitParams{
		Router:       router,
		ClientId:     "xx-xxx-xxx",
		ClientSecret: "xx-xxx-xxx",
		Issuer:       *issuer,
		ClientUrl:    *clienturl,
		Scopes:       []string{"openid", "profile", "email"},
		ErrorHandler: func(c *gin.Context) {
			message := c.Errors.Last().Error()
			c.String(302, message)
		},
		PostLogoutUrl: *postlogouturl,
	}

	//protect all endpoint below this line
	router.Use(gin_oidc.Init(initParams))

	log.Fatal(router.Run(":8080"))
}

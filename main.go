package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sw "github.com/pvr1/gigs/go"
	gin_oidc "github.com/pvr1/gin-oidc"
)

func main() {
	log.Printf("Server started")

	gin.SetMode(gin.ReleaseMode)

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
	router.Use(gin_oidc.Init(initParams), cors.New(cors.Config{
		AllowOrigins:     []string{"https://foo.com", "*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	log.Fatal(router.Run(":8080"))
}

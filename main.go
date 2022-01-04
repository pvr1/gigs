package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sw "github.com/pvr1/gigs/go"
)

func main() {
	log.Printf("Server started")

	gin.SetMode(gin.ReleaseMode)

	router := sw.NewRouter()

	//protect all endpoint below this line
	router.Use(cors.New(cors.Config{
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

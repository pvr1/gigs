package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sw "github.com/pvr1/gigs/go"
	"github.com/pvr1/gigs/go/platform/authenticator"
)

var (
	clientID     = os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
)

func main() {
	log.Printf("Server started")

	gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	router := sw.NewRouter(auth)

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

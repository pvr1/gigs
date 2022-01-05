package main

import (
	"encoding/json"
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/joho/godotenv"

	sw "github.com/pvr1/gigs/go"
)

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

	payload, err := json.Marshal(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
})

func testJWT(middleware *jwtmiddleware.JWTMiddleware, handler http.HandlerFunc) http.Handler {
	log.Println("CheckJWT")
	return middleware.CheckJWT(handler)
}

func main() {
	log.Printf("Server started")

	//gin.SetMode(gin.ReleaseMode)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	/*
		auth, err := authenticator.New()
		if err != nil {
			log.Fatalf("Failed to initialize the authenticator: %v", err)
		}
	*/

	/*
		issuerURL, err := url.Parse("http://localhost:8080/v2/")
		if err != nil {
			log.Fatalf("failed to parse the issuer url: %v", err)
		}

		provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

		// Set up the validator.
		jwtValidator, err := validator.New(
			provider.KeyFunc,
			validator.RS256,
			issuerURL.String(),
			[]string{"https://dev-4du4iqv3.eu.auth0.com/api/v2/"},
		)
		if err != nil {
			log.Fatalf("failed to set up the validator: %v", err)
		}

		// Set up the middleware.

		middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	*/
	router := sw.NewRouter()
	//router.Use(gin.WrapH(testJWT(middleware, handler)))

	//protect all endpoint below this line
	/*
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
	*/

	log.Fatal(router.Run(":8080"))
}

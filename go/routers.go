package openapi

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pvr1/gigs/go/platform/authenticator"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	log.Println("Claims: {}", claims)

	payload, err := json.Marshal(claims)
	log.Println("Payload: {}", payload)
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

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))
	issuerURL, err := url.Parse("http://localhost:8080/v2/")
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{"https://dev-4du4iqv3.eu.auth0.com/api/v2/"},
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}
	middleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, gin.WrapH(middleware.CheckJWT(handler)), route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, gin.WrapH(middleware.CheckJWT(handler)), route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, gin.WrapH(middleware.CheckJWT(handler)), route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, gin.WrapH(middleware.CheckJWT(handler)), route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, gin.WrapH(middleware.CheckJWT(handler)), route.HandlerFunc)
		}
	}

	/*
		router.GET("/login", LoginHandler(auth))
		router.GET("/callback", CallbackHandler(auth))
		router.GET("/user", IsAuthenticated, UserHandler)
		router.GET("/logout", LogoutHandler)
	*/

	return router
}

// NewRouter returns a new router.
func NewTestRouter() *gin.Engine {
	router := gin.Default()

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// LoginHandler for our login.
func LoginHandler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

// LogoutHandler for our logout
func LogoutHandler(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

// CallbackHandler for our callback.
func CallbackHandler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to convert an authorization code into a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}

// UserHandler for our logged-in user page.
func UserHandler(ctx *gin.Context) {
	//session := sessions.Default(ctx)
	//profile := session.Get("profile")

	myResponse := "Try localhost:8080/v2/"

	ctx.String(http.StatusOK, myResponse)
	//ctx.HTML(http.StatusOK, "landingpage.html", profile)
}

// generateRandomState generates a random string suitable for CSRF protection.
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// Index is the index handler.
func Index(c *gin.Context) {
	myResponse := "Hello Gigs World!\n"
	myResponse += "\nThis is the GigsAPI complying to the OpenAPI3.0 specification\n\n"
	myResponse += "The endpoint for the JSON file is:\n"
	myResponse += "/v2/openapi/json\n\n"
	myResponse += "The endpoint for the Yaml file is:\n"
	myResponse += "/v2/openapi/yaml\n"

	c.String(http.StatusOK, myResponse)
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/v2/",
		Index,
	},

	{
		"AddGig",
		http.MethodPost,
		"/v2/gigs",
		AddGig,
	},

	{
		"DeleteGig",
		http.MethodDelete,
		"/v2/gig/:gigId",
		DeleteGig,
	},

	{
		"FindGigsByStatus",
		http.MethodGet,
		"/v2/gig/findByStatus",
		FindGigsByStatus,
	},

	{
		"GetGigById",
		http.MethodGet,
		"/v2/gig/:gigId",
		GetGigById,
	},

	{
		"UpdateGigWithForm",
		http.MethodPost,
		"/v2/gig/:gigId",
		UpdateGigWithForm,
	},

	{
		"UploadFile",
		http.MethodPost,
		"/v2/gig/:gigId/uploadImage",
		UploadFile,
	},

	{
		"UpdateGig",
		http.MethodPut,
		"/v2/gigs",
		UpdateGig,
	},

	{
		"DeleteOrder",
		http.MethodDelete,
		"/v2/store/order/:orderId",
		DeleteOrder,
	},

	{
		"GetInventory",
		http.MethodGet,
		"/v2/store/inventory",
		GetInventory,
	},

	{
		"GetOrderById",
		http.MethodGet,
		"/v2/store/order/:orderId",
		GetOrderById,
	},

	{
		"PlaceOrder",
		http.MethodPost,
		"/v2/store/order",
		PlaceOrder,
	},

	{
		"CreateUser",
		http.MethodPost,
		"/v2/user",
		CreateUser,
	},

	{
		"CreateUsersWithArrayInput",
		http.MethodPost,
		"/v2/user/createWithArray",
		CreateUsersWithArrayInput,
	},

	{
		"CreateUsersWithListInput",
		http.MethodPost,
		"/v2/user/createWithList",
		CreateUsersWithListInput,
	},

	{
		"DeleteUser",
		http.MethodDelete,
		"/v2/user/:username",
		DeleteUser,
	},

	{
		"GetUserByName",
		http.MethodGet,
		"/v2/user/:username",
		GetUserByName,
	},

	{
		"LoginUser",
		http.MethodGet,
		"/v2/user/login",
		LoginUser,
	},

	{
		"LogoutUser",
		http.MethodGet,
		"/v2/user/logout",
		LogoutUser,
	},

	{
		"UpdateUser",
		http.MethodPut,
		"/v2/user/:username",
		UpdateUser,
	},

	{
		"OpenAPIjson",
		http.MethodGet,
		"/v2/openapi/json",
		OpenAPIjson,
	},

	{
		"OpenAPIyaml",
		http.MethodGet,
		"/v2/openapi/yaml",
		OpenAPIyaml,
	},
}

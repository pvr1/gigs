package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
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

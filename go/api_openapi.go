package openapi

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// OpenAPIjson - Get OpenAPI 3.0 JSON
func OpenAPIjson(c *gin.Context) {
	file, err := os.Open("./api/openapi.json")
	if err != nil {
		c.String(http.StatusOK, "Could not find openapi file")
		return
	} else {
		defer file.Close()
	}

	filecontent, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusOK, "Could not read openapi file after loading it")
		return
	}

	c.String(http.StatusOK, string(filecontent))
}

// OpenAPIyaml - Get OpenAPI 3.0 YAML
func OpenAPIyaml(c *gin.Context) {
	file, err := os.Open("./api/openapi.yaml")
	if err != nil {
		c.String(http.StatusOK, "Could not find openapi file")
		return
	} else {
		defer file.Close()
	}

	filecontent, err := ioutil.ReadAll(file)
	if err != nil {
		c.String(http.StatusOK, "Could not read openapi file after loading it")
		return
	}
	c.String(http.StatusOK, string(filecontent))

}

package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const swaggerFilePath = "api/openapi.bundled.json"

// GetSwaggerJSON serves the combined OpenAPI specification from the generated static file
func GetSwaggerJSON(c *gin.Context) {
	path := viper.GetString("http.swagger.path")
	if path == "" {
		path = swaggerFilePath
	}
	// Read the pre-generated swagger.json file
	data, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read OpenAPI spec. Please run 'make generate' first.",
		})
		return
	}

	// Parse JSON to validate it
	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid OpenAPI spec format",
		})
		return
	}

	c.JSON(http.StatusOK, spec)
}

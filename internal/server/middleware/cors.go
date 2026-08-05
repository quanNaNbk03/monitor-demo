package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func CORS() gin.HandlerFunc {
	methods := viper.GetStringSlice("http.server.cors.allowMethod")
	domains := viper.GetStringSlice("http.server.cors.allowDomain")

	config := cors.Config{
		AllowMethods: methods,
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
	if domains[0] == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = domains
	}

	return cors.New(config)
}

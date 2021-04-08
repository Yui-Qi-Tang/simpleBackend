package cors

import (
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Default is a default setting of cors for go-gin
func Default() gin.HandlerFunc {

	corsConfig := ginCors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"DNT", "X-CustomHeader", "Keep-Alive", "User-Agent,X-Requested-With", "If-Modified-Since,Cache-Control", "Content-Type"},
	}
	return ginCors.New(corsConfig)
}

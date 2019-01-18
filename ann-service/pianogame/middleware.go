package pianogame

import (
	"log"

	"github.com/gin-gonic/gin"
)

func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Middle First!!")
	}
}

func TestMiddleware2(c *gin.Context) {
	log.Println("Hi M2")
	c.Next()
}

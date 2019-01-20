package pianogame

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AuthCheck(c *gin.Context) {
	// TODO: check JWT here
	/*
		    if valid {
				pass
			} else {
				redirect to sign up URL by pass http redirect code for browser
			}
	*/
	log.Println("Hi, I am auth checker!")
}

func MiddlewareForMysqlTest(c *gin.Context) {
	// Just demo
	log.Println("Hi, I am middleware for mysql")
	//c.Next()
}

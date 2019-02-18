package pianogame

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticationCheck(c *gin.Context) {
	// TODO: check JWT here
	/*
		    if valid {
				pass
			} else {
				redirect to sign up URL by pass http redirect code for browser
			}
	*/
	if cookie, err := c.Cookie("token"); err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		log.Println("Hi, I am auth checker!", cookie)
	}
}

func MiddlewareForMysqlTest(c *gin.Context) {
	// Just demo
	log.Println("Hi, I am middleware for mysql")
	//c.Next()
}

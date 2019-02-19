package pianogame

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthenticationCheck check if JWT token is valid
func AuthenticationCheck(c *gin.Context) {
	// TODO: do not check the request from websocket
	// log.Println("enter auth check")
	if cookie, err := c.Cookie("token"); err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		if !IsMemberJWTValid(cookie) || IsMemberJWTExpired(cookie) {
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
	}
	log.Println("Authentication pass")
	c.Next()
}

func MiddlewareForMysqlTest(c *gin.Context) {
	// Just demo
	log.Println("Hi, I am middleware for mysql")
	//c.Next()
}

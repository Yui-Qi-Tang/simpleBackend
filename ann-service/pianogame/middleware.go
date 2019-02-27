package pianogame

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthenticationCheck check if JWT token is valid
func AuthenticationCheck(c *gin.Context) {
	if cookie, err := c.Cookie("token"); err != nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
	} else {
		if !IsMemberJWTValid(cookie) || IsMemberJWTExpired(cookie) {
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
	}
	c.Next()
}

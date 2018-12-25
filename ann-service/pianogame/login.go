package pianogame


import (
	"net/http"	
	"github.com/gin-gonic/gin"
)

// LoginPage show login UI from template page
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{
		"loginURL": "http://127.0.0.1:8080/user/login",
	})
}
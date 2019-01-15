package pianogame

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage show login UI from template page
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Login.html", gin.H{
		"loginURL":   "http://127.0.0.1:8080/user/login",
		"signupPage": "http://127.0.0.1:8080/signup",
		"loginPage":  "http://127.0.0.1:8080/login",
	})
}

// SignupPage show sign-up UI from template page
func SignupPage(c *gin.Context) {
	c.HTML(http.StatusOK, "Signup.html", gin.H{
		"registerURL": "http://127.0.0.1:8080/user/register",
	})
}

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
	host := getURLInfo(c)
	registerURL := "/user/register"
	c.HTML(http.StatusOK, "Signup.html", gin.H{
		"registerURL": strConcate(host.String(), registerURL),
	})
}

// GamePage game page
func GamePage(c *gin.Context) {
	c.HTML(http.StatusOK, "Gamepage.html", gin.H{})
}

// IndexPage index page
func IndexPage(c *gin.Context) {
	resources := [...]string{
		"/js/annPage.js",
		"/css/annPage.css",
		"/images/piano_background.jpg",
		"/js/jquery-3.3.1.min.js",
		"/images/piano_2.jpg",
		"/images/Piano.jpg",
		"/music/music.mp3",
	}
	for _, v := range resources {
		go webPusher(c, v)
	}
	c.HTML(http.StatusOK, "AnnPage.html", gin.H{})
}

package pianogame

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoginPage show login UI from template page
func LoginPage(c *gin.Context) {
	host := getURLInfo(c)
	// Bad design for concate website and API service...TODO, use API gateway
	signinAPI := strConcate("https://", UserAPIConfig.User.Network[1].Name, ":", strconv.Itoa(UserAPIConfig.User.Network[1].Port), "/member/v2/user/validation")
	c.HTML(http.StatusOK, "Login.html", gin.H{
		"loginURL":   signinAPI,
		"signupPage": strConcate(host.String(), "/signup"),
		"loginPage":  strConcate(host.String(), "/login"),
	})
}

// SignupPage show sign-up UI from template page
func SignupPage(c *gin.Context) {
	// Bad design for concate website and API service...TODO, use API gateway
	signupAPI := strConcate("https://", UserAPIConfig.User.Network[1].Name, ":", strconv.Itoa(UserAPIConfig.User.Network[1].Port), "/member/v2/user")
	c.HTML(http.StatusOK, "Signup.html", gin.H{
		"registerURL": signupAPI,
	})
}

// GamePage game page
func GamePage(c *gin.Context) {
	c.HTML(http.StatusOK, "Gamepage.html", gin.H{})
}

// IndexPage index page; just demo
func IndexPage(c *gin.Context) {
	resources := []string{
		"/js/annPage.js",
		"/css/annPage.css",
		"/images/piano_background.jpg",
		"/js/jquery-3.3.1.min.js",
		"/images/piano_2.jpg",
		"/images/Piano.jpg",
		"/music/music.mp3",
	}
	pusher := c.Writer.Pusher()
	for _, v := range resources {
		webPusher(pusher, v)
	}
	c.HTML(http.StatusOK, "AnnPage.html", gin.H{})
}

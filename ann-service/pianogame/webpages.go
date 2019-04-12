package pianogame

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginPage show login UI from template page
func LoginPage(c *gin.Context) {
	resources := []string{
		"/css/login.css",
		"/js/login.js",
	}
	pusher := c.Writer.Pusher()
	for _, v := range resources {
		webPusher(pusher, v)
	}
	host := getURLInfo(c)
	// Bad design for concate website and API service...TODO, use API gateway
	c.HTML(http.StatusOK, "Login.html", gin.H{
		"loginURL":   strConcate(host.String(), "/login"),
		"signupPage": strConcate(host.String(), "/signup"),
		"loginPage":  strConcate(host.String(), "/login"),
	})
}

// SignupPage show sign-up UI from template page
func SignupPage(c *gin.Context) {
	// Bad design for concate website and API service...TODO, use API gateway
	c.HTML(http.StatusOK, "Signup.html", gin.H{
		"registerURL": "/signup",
	})
}

// GamePage game page
func GamePage(c *gin.Context) {
	resources := []string{
		"/js/websocket.js",
		"/js/getCookie.js",
		"/js/jquery-3.3.1.min.js",
		"/css/game.css",
	}
	pusher := c.Writer.Pusher()
	for _, v := range resources {
		webPusher(pusher, v)
	}
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

// for http2 push
func webPusher(p http.Pusher, resource string) {
	if p != nil {
		if err := p.Push(resource, nil); err != nil {
			log.Printf("Failed to push: %v", err)
		}
	}
}

package main

import (
	"log"
	"net/http"
	"simpleBackend/ann-service/pianogame"

	"github.com/gin-gonic/gin"
)


func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection?
			2. bind mongodb and go gin api together
			3. refactor
	*/
	log.Println("Hello world, SimpleBackend!!")

	// Go Gin
	gin.SetMode(gin.TestMode) // enable server on localhost:8080
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	// ann-service here
	router.LoadHTMLGlob("ann-service/pianogame/templates/*") // load tempates
	router.POST("user/login", pianogame.UserLogin)
	// register account
	router.POST("user/register", pianogame.UserRegister)

	router.GET("/login", pianogame.LoginPage)
	router.GET("/signup", pianogame.SignupPage)
	

	router.Run() // listen and serve on 127.0.0.1:8080 in gin.TestMode
}

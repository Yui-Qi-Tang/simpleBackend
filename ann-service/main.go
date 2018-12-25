package main

import (
	"log"
	"simpleBackend/ann-service/pianogame"

	"github.com/gin-gonic/gin"
)

// main ann-service entry point */
func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection, in mongdb.go??

				mode: variable is denoted the status of gin(test/production)
			2. add JWT for auth
	*/
	/* Go-Gin setup */
	gin.SetMode(gin.TestMode) // enable server on localhost:8080
	router := gin.Default()
	router.LoadHTMLGlob("pianogame/templates/*") // load tempates

	/* APIs */
	router.POST("user/login", pianogame.UserLogin)  // login
	router.POST("user/register", pianogame.UserRegister) // signup

    /* Web page */
	router.GET("/login", pianogame.LoginPage) // login page
	router.GET("/signup", pianogame.SignupPage)  // signup page
	
	/* Run server */
	log.Println("Start server")
	router.Run() // listen and serve on 127.0.0.1:8080 in gin.TestMode
}

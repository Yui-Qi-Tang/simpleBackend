package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"simpleBackend/ann-service/pianogame"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func runserverTLS(server *http.Server, cert string, key string) {
	// Start HTTPS server by net/http
	if err := server.ListenAndServeTLS(cert, key); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func shutDownGraceful(server *http.Server) {
	log.Println("Server graceful exiting")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

func waitQuitSignal(hint string) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println(hint)
}

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
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	/* Use middleware */
	router.Use(gin.Recovery())
	router.Use(location.New(location.DefaultConfig()))
	router.LoadHTMLFiles(pianogame.SysConfig.HTMLTemplates...) // load tempates (Parameters is variadic), ref: https://golang.org/ref/spec#Passing_arguments_to_..._parameters

	// set static files
	router.Static("/js", pianogame.SysConfig.Static.Js)
	router.Static("/css", pianogame.SysConfig.Static.CSS)
	router.Static("/images", pianogame.SysConfig.Static.Images)
	router.Static("/music", pianogame.SysConfig.Static.Music)

	userRoute := router.Group("user")
	mysqlRoute := router.Group("mysql")

	/* APIs */
	userRoute.POST("/login", pianogame.UserLogin)       // login
	userRoute.POST("/register", pianogame.UserRegister) // signup

	mysqlRoute.POST("/test", pianogame.MysqlCheckTable)        // just test
	mysqlRoute.POST("/user/test", pianogame.InsertUserToMysql) // just test
	mysqlRoute.GET("/user", pianogame.GetUsers)                // just test
	mysqlRoute.DELETE("/user", pianogame.DeleteUser)           // just test

	router.POST("/upload", pianogame.UploadFileSample) // file upload demo

	router.POST("/parsejwt", pianogame.DecodeJwt)

	/* Web page */
	router.GET("/login", pianogame.LoginPage)   // login page
	router.GET("/signup", pianogame.SignupPage) // signup page
	router.GET("/game", pianogame.GamePage)     // game page
	router.GET("/", pianogame.IndexPage)        // index page

	/* Run server */
	log.Println("Start server")

	srv := &http.Server{
		Addr:    pianogame.SysConfig.Port,
		Handler: router,
	}

	go runserverTLS(srv, pianogame.SysConfig.Ssl.Cert, pianogame.SysConfig.Ssl.Key)

	/* Graceful shotdown */
	waitQuitSignal("Receive Quit server Signal") // block until receive quit signal
	shutDownGraceful(srv)

	defer pianogame.MysqlDB.Close()
}

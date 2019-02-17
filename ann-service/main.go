package main

import (
	"simpleBackend/ann-service/pianogame"

	"simpleBackend/ann-service/pianogame/clientapi"

	"simpleBackend/ann-service/pianogame/protocol-buffer/pbserver"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

// go:generate protoc -I pianogame/grpc/ pianogame/grpc/user_service.proto --go_out=plugins=grpc:pianogame/grpc

// main ann-service entry point */
func main() {
	/*
		TO-DO:
			1. load config from file for mongodb
			    Base data: host/port/connect protocol
				Q: specify db and collection, in mongdb.go??

				mode: variable is denoted the status of gin(test/production)
			2. add JWT for auth
			3. Website <-gRPC-> api auth
				 need to create an API as a wrapper for internal API
				 An front-API in website to receive data;
				 a 'middler' receives the data from front api and push data to back-API-service
				 Fig.
					user request -HTTP-> front-API on website -gRPC-> back-API-service
			4. Use NewSQL server lol
	*/

	/* Go-Gin setup */
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	/* Use middleware */
	router.Use(gin.Recovery())
	router.Use(location.New(location.DefaultConfig()))
	router.Use(pianogame.AuthCheck)
	router.LoadHTMLFiles(pianogame.WebConfig.Settings.HTMLTemplates...) // load tempates (Parameters is variadic), ref: https://golang.org/ref/spec#Passing_arguments_to_..._parameters

	// set static files
	router.Static("/js", pianogame.WebConfig.Settings.Static.Js)
	router.Static("/css", pianogame.WebConfig.Settings.Static.CSS)
	router.Static("/images", pianogame.WebConfig.Settings.Static.Images)
	router.Static("/music", pianogame.WebConfig.Settings.Static.Music)

	userRoute := router.Group("user")
	// mysqlRoute := router.Group("mysql")
	// mysqlRoute.Use(pianogame.MiddlewareForMysqlTest) // my first middle for auth

	/* Front APIs */
	userRoute.POST("/login", pianogame.UserLogin)       // login
	userRoute.POST("/register", pianogame.UserRegister) // signup
	router.POST("/upload", pianogame.UploadFileSample)  // file upload demo
	router.POST("/parsejwt", pianogame.DecodeJwt)
	router.POST("/parse-cookie-jwt", pianogame.DecodeJwtFromCookie)
	router.GET("/game/socket", pianogame.GameWebSocketHandler)

	router.POST("/fake/login", clientapi.Login)

	/* Web page */
	router.GET("/login", pianogame.LoginPage)   // login page
	router.GET("/signup", pianogame.SignupPage) // signup page
	router.GET("/game", pianogame.GamePage)     // game page
	router.GET("/", pianogame.IndexPage)        // index page

	/* Start servers  */
	pianogame.ServiceInstances = append(
		pianogame.StartServers(router, pianogame.WebConfig.Settings.Network, pianogame.WebConfig.Settings.Meta),
		pianogame.StartServers(pianogame.UserServiceRouter(), pianogame.UserAPIConfig.User.Network, pianogame.UserAPIConfig.User.Meta)...,
	)
	/* gRPC server */
	go pbserver.StartAuthenticationService()
	/*
		HINT: if there does exist another serivce, please append http instances again:

		pianogame.ServiceInstances = append(
			pianogame.ServiceInstances,
			pianogame.StartServers(routerForService, pianogame.Service.Network)...
		)
		...

		another again? please append pianogame.ServiceInstances again an so on...
		For now, I think it's a bad idea to set multiple service, lol
	*/

	pianogame.WaitQuitSignal("Receive Quit server Signal") // block until receive quit signal from system

	// stop servers
	for _, v := range pianogame.ServiceInstances {
		pianogame.ShutDownGraceful(v) // terminate each server
	} // for

	defer pianogame.MysqlDB.Close()
}

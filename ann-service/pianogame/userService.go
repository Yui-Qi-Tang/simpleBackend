package pianogame

import (
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func init() {
	// This is blank
}

// UserServiceRouter Simple way to create an gin router for user micro-service
func UserServiceRouter() *gin.Engine {
	/*
	   TODO: Let api function as private in pianogame package
	*/
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(location.New(location.DefaultConfig()))
	// router.Use(AuthCheck)
	// router.Use(MiddlewareForMysqlTest) // my first middle for auth

	authRoute := router.Group("member/v2/")
	authRoute.POST("/user", AddUser)                   // api for sign-up
	authRoute.POST("/user/validation", UserValidation) // api for sign-in
	authRoute.GET("/user/:token/", GetUserInfoByToken) // just test

	authRoute.POST("/test", MysqlCheckTable) // just test
	authRoute.GET("/user", GetUsers)         // just test
	authRoute.DELETE("/user", DeleteUser)    // just test

	return router
}

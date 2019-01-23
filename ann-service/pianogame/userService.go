package pianogame

import (
	"fmt"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
)

func init() {
	// This is blank
	// fmt.Println(SysConfig)
	// fmt.Println(Ssl.Path.Cert)
	fmt.Println(authSettings.Secret.Jwt)
}

// UserServiceRouter Simple way to create an gin router for user micro-service
func UserServiceRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(location.New(location.DefaultConfig()))
	router.Use(AuthCheck)
	router.Use(MiddlewareForMysqlTest) // my first middle for auth

	authRoute := router.Group("mysql")

	authRoute.POST("/test", MysqlCheckTable)        // just test
	authRoute.POST("/user/test", InsertUserToMysql) // just test
	authRoute.GET("/user", GetUsers)                // just test
	authRoute.DELETE("/user", DeleteUser)           // just test

	return router
}

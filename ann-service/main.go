package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"simpleBackend/ann-service/pianogame"

	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
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
	/* Load config by yaml format */
	var config pianogame.Config
	configFile, configErr := ioutil.ReadFile("config/api/config.yaml") // open file and read
	if configErr != nil {
		log.Panicf("read service config file error: %v", configErr)
	} // fi
	configUnmarshalError := yaml.Unmarshal([]byte(configFile), &config)
	if configUnmarshalError != nil {
		log.Panicf("error while unmarshal from db config: %v", configUnmarshalError)
	} // fi
	log.Println("Load config file finished")

	/* Go-Gin setup */
	gin.SetMode(gin.TestMode) // enable server on localhost:8080
	router := gin.Default()
	router.LoadHTMLFiles(config.HTMLTemplates...) // load tempates (Parameters is variadic), ref: https://golang.org/ref/spec#Passing_arguments_to_..._parameters

	// set static files
	//router.Static("/js", "./js")
	//router.Static("/css", "./css")
	userRoute := router.Group("user")
	mysqlRoute := router.Group("mysql")

	/* APIs */
	userRoute.POST("/login", pianogame.UserLogin)              // login
	userRoute.POST("/register", pianogame.UserRegister)        // signup
	mysqlRoute.POST("/test", pianogame.MysqlCheckTable)        // just test
	mysqlRoute.POST("/user/test", pianogame.InsertUserToMysql) // just test

	mysqlRoute.GET("/user", pianogame.GetUsers) // just test

	mysqlRoute.DELETE("/user", pianogame.DeleteUser) // just test

	router.POST("/upload", func(c *gin.Context) {
		// single file
		savePlace := "/tmp"

		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusOK, "File upload error!!")
		}

		// log.Println(file.Filename)

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", savePlace, file.Filename))

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	/* Web page */
	router.GET("/login", pianogame.LoginPage)   // login page
	router.GET("/signup", pianogame.SignupPage) // signup page

	/* Run server */
	log.Println("Start server")
	router.Run()
	// Close Mysql DB client
	defer pianogame.MysqlDB.Close()
}

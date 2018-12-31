package main

import (
	"io/ioutil"
	"log"
	"simpleBackend/ann-service/pianogame"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
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

	/* APIs */
	router.POST("user/login", pianogame.UserLogin)              // login
	router.POST("user/register", pianogame.UserRegister)        // signup
	router.POST("mysql/test", pianogame.MysqlCheckTable)        // just test
	router.POST("mysql/user/test", pianogame.InsertUserToMysql) // just test

	router.GET("mysql/user/", pianogame.GetUsers) // just test

	router.DELETE("mysql/user/", pianogame.DeleteUser) // just test

	/* Web page */
	router.GET("/login", pianogame.LoginPage)   // login page
	router.GET("/signup", pianogame.SignupPage) // signup page

	/* Run server */
	log.Println("Start server")

	/* set TLS with LetsEncrypt
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("myserveryuki1.com", "myserveryuki2.com"),
		Cache:      autocert.DirCache("/tmp/.cache"),
	}
	log.Println(m)
	log.Fatal(autotls.RunWithManager(router, &m))
	log.Fatal(autotls.Run(router, "myserveryuki.com"))
	*/
	router.Run() // listen and serve on 127.0.0.1:8080 in gin.TestMode

	// Close Mysql DB client
	defer pianogame.MysqlDB.Close()
}

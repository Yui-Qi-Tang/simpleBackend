package main

import (
	"io/ioutil"
	"simpleBackend/http2server"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func main() {
	serverConfigYaml := http2server.ServerConfigYaml{}
	//serverConfig.Port = ":8081"
	//serverConfig.Crt = "https_keys/server.crt"
	//serverConfig.Key = "https_keys/server.key"
	//serverConfig.Static = "static"

	// http2server.CreateServer(
	// 	serverConfig,
	// )
	// resultMap := make(map[string]interface{})
	if content, ioErr := ioutil.ReadFile("config/server.yaml"); ioErr != nil {
		logrus.Fatalf("read service config file error: %v", ioErr)
	} else {
		if ymlErr := yaml.Unmarshal(content, &serverConfigYaml); ymlErr != nil {
			logrus.Fatalf("error while unmarshal from db config: %v", ymlErr)
		}
	}

	// fmt.Println(resultMap)
}

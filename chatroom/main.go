package main

import (
	"io/ioutil"
	"log"
	"simpleBackend/http2server" // "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

func main() {
	serverConfigYaml := http2server.ServerConfigYaml{}

	// read yaml file
	yamlContent, ioErr := ioutil.ReadFile("config/server.yaml") // open file and read
	if ioErr != nil {
		log.Fatalf("read service config file error: %v", ioErr)
		return
	} // fi

	// save yamlContent to structure
	/*
		// another way to save data that decoded by yaml.Unmarshal without pre-defined structure
		resultMap := make(map[string]interface{})  // a way to get data from yaml file without defined structure
		yaml.Unmarshal(yamlContent, resultMap)
	*/
	yamlError := yaml.Unmarshal(yamlContent, &serverConfigYaml)
	if yamlError != nil {
		log.Fatalf("error while unmarshal from db config: %v", yamlError)
		return
	}

	log.Println("Config ok!")

	http2server.CreateServerByYaml(&serverConfigYaml)
}

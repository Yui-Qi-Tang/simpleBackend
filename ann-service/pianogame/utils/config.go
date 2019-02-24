package utils

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// ReadFile read file by ioutil
func ReadFile(filePath string) []byte {
	fileBytes, err := ioutil.ReadFile(filePath) // open file and read
	ErrorCheck(err, "readFile Error")
	return fileBytes
}

// LoadYAMLConfig load cofig file with yaml
func LoadYAMLConfig(configFilePath, errMsg, successMsg string, configStructure interface{}) {
	// This is blank
	bytesData := ReadFile(configFilePath)
	configUnmarshalError := yaml.Unmarshal(bytesData, configStructure)
	ErrorCheck(configUnmarshalError, errMsg)
	log.Println(configFilePath, successMsg)
}

// ErrorCheck ErrorCheck function, just dump error during develop
func ErrorCheck(e error, msg ...string) {
	// TO-DO: better to logging error
	if e != nil {
		errorMsg := ""
		for _, v := range msg {
			errorMsg += v
		}
		log.Panicf("%s => %v", errorMsg, e)
	} // fi
}

package pianogame

import (
	"log"

	yaml "gopkg.in/yaml.v2"
)

// SysConfig TO-DO: split some data field out
var SysConfig Config

func init() {
	/* Load API config data */
	bytesData := readFile("config/api/config.yaml")
	configUnmarshalError := yaml.Unmarshal(bytesData, &SysConfig) // TO-DO: a data formater(yaml or json)
	errorCheck(configUnmarshalError, "error while unmarshal from db config")
	log.Println("Load API config file finished")
}

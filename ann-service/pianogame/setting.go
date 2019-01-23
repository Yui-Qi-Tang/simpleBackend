package pianogame

import "log"

// SysConfig TO-DO: split some data field out
var SysConfig Config

// Ssl ssl settings
var Ssl SSLPath

var authSettings Auth

var apigw apiGW

func init() {
	/* Load API config data */
	loadYAMLConfig(
		"config/api/config.yaml",
		"error while unmarshal from API config",
		"Load API config file finished",
		&SysConfig,
	)
	loadYAMLConfig(
		"config/api/config.yaml",
		"error while unmarshal from API config",
		"Load API config file finished",
		&apigw,
	)
	log.Println(apigw)
	/* Load SSL config */
	loadYAMLConfig(
		"config/ssl/config.yaml",
		"error while unmarshal from ssl config",
		"Load SSL config file finished",
		&Ssl,
	)
	/* Load auth secret */
	loadYAMLConfig(
		"config/auth/config.yaml",
		"error while unmarshal from auth config",
		"Load Auth config file finished",
		&authSettings,
	)

}

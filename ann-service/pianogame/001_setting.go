package pianogame

/*
   HINT: init order of Golang is by file name
*/

import (
	"simpleBackend/ann-service/pianogame/datastructure"
)

// Ssl ssl settings
var Ssl SSLPath

var authSettings Auth

// UserAPIConfig setting
var UserAPIConfig userAPI

// WebConfig settings for website
var WebConfig WebSiteConfig

// MongoConfig mongo db config
var MongoConfig datastructure.MongoDBSetting

// GrpcConfig grpc config
var GrpcConfig datastructure.GRPCSetting

func init() {
	/* Load Web site config */
	loadYAMLConfig(
		"config/website/config.yaml",
		"error while unmarshal from website config",
		"Load website config file finished",
		&WebConfig,
	)
	/* Load API config */
	loadYAMLConfig(
		"config/api/config.yaml",
		"error while unmarshal from API config",
		"Load API config file finished",
		&UserAPIConfig,
	)
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
	/* Load mong db config */
	loadYAMLConfig(
		"config/database/mongo/config.yaml",
		"error while unmarshal from mongo db config",
		"Load mongo db config file finished",
		&MongoConfig,
	)
	/* Load grpc config */
	loadYAMLConfig(
		"config/grpc/config.yaml",
		"error while unmarshal from grpc config",
		"Load grpc config file finished",
		&GrpcConfig,
	)
}

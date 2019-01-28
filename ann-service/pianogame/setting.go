package pianogame

// Ssl ssl settings
var Ssl SSLPath

var authSettings Auth

// UserAPIConfig setting
var UserAPIConfig userAPI

// WebConfig settings for website
var WebConfig WebSiteConfig

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
}

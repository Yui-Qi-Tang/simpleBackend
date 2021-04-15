package load

import (
	"time"
)

// Config is the configuration for system
type Config struct {
	HTTPSrv    HTTPServer    `yaml:"httpserver"`
	ThirdParty ThirdPartyAPI `yaml:"third-party-api"`
	DB         Database      `yaml:"databases"`
}

// HTTPServer is the configuration about http server
type HTTPServer struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

// ThirdPartyAPI saves the information of 3-rd api
type ThirdPartyAPI struct {
	Nasa ThirdPartyAPINasa `yaml:"nasa"`
}

// ThirdPartyAPINasa save the nasa api key
type ThirdPartyAPINasa struct {
	Key string `yaml:"api_key"`
}

// Database saves the connection info of different databases
type Database struct {
	Main  DatabaseMain  `yaml:"main"`
	Redis DataBaseRedis `yaml:"redis"`
}

// DatabaseMain saves the main database information
type DatabaseMain struct {
	Type string `yaml:"type"`
	DSN  string `yaml:"dsn"`

	MaxOpenConns int           `yaml:"max_open_conns"`
	MaxIdleConns int           `yaml:"max_idle_conns"`
	MaxLifeTime  time.Duration `yaml:"conn_max_life_time_second"`
}

// DatabaseRedis save the configurations for Redis server
type DataBaseRedis struct {
	Password string `yaml:"password"`
	Addr     string `yaml:"addr"`
}

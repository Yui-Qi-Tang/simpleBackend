package pianogame

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// ServiceInstances  an array to store the http server instance
var ServiceInstances []*http.Server

// Login structure for user login
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type serviceMeta struct {
	Name string `yaml:"name"`
}

// WebSiteConfig for config website
type WebSiteConfig struct {
	Settings webSiteSettings `yaml:"web_service"`
}

type webSiteSettings struct {
	Network       []host      `yaml:"addrs"`
	HTMLTemplates []string    `yaml:"html_templates"`
	Static        staticPath  `yaml:"static"`
	Meta          serviceMeta `yaml:"meta"`
}

// MysqlConfig structure for set mysql db
type MysqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
	DBName   string `yaml:"database_name"`
}

type staticPath struct {
	CSS    string `yaml:"css"`
	Js     string `yaml:"js"`
	Images string `yaml:"images"`
	Music  string `yaml:"music"`
}

// SSLPath ssl file path
type SSLPath struct {
	Path ssl `yaml:"ssl"`
}

type ssl struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type jwtClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

/* API */
type userAPI struct {
	User apiUserService `yaml:"user_api"`
}

type apiUserService struct {
	Network []host      `yaml:"addrs"`
	Meta    serviceMeta `yaml:"meta"`
}

type host struct {
	Name string `yaml:"hostname"`
	Port int    `yaml:"port"`
}

/* For Auth */

// AuthData authorization
type authData struct {
	Token string `json:"token"` // JWT
}

type Auth struct {
	Secret authSecret `yaml:"secret"`
}

type authSecret struct {
	Jwt string `yaml:"jwt"`
}

type AuthMemberClaim struct {
	ID uint `json:"ID"`
	jwt.StandardClaims
}

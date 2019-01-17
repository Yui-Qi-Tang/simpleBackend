package pianogame

import (
	jwt "github.com/dgrijalva/jwt-go"
)

// Login structure for user login
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// Config structure for set API server
type Config struct {
	Debug         bool       `yaml:"debug"`
	HTMLTemplates []string   `yaml:"html_templates"`
	Static        staticPath `yaml:"static"`
	JwtSec        string     `yaml:"jwtSecret"`
	Ssl           ssl        `yaml:"ssl"`
	Port          int        `yaml:"port"`
	IP            string     `yaml:"ip"`
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

type ssl struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// AuthData authorization
type authData struct {
	Token string `json:"token"` // JWT
}

type jwtClaim struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

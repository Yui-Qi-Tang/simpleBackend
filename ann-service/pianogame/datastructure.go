package pianogame

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

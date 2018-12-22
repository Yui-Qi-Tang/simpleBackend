package http2server

// ServerConfig comment...
type ServerConfig struct {
	Port   string
	Crt    string
	Key    string
	Static string
}

// ServerConfigYaml A structure for yaml config
type ServerConfigYaml struct {
	Port        string `yaml:"port"`
	SslFilePath struct {
		Crt string `yaml:"crt"`
		Key string `yaml:"key"`
	} `yaml:"sslFilePath"`
	Static string `yaml:"staticFile"`
}

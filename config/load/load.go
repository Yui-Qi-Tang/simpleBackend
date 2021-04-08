package load

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadFromFile returns Config from yaml file
func LoadFromFile(file string) (*Config, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load config")
	}

	return yamlUnmarshal(b)
}

func yamlUnmarshal(b []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(b, config); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}
	return config, nil
}

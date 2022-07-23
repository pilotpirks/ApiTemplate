package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// config.yaml
type Config struct {
	DbSettings struct {
		DbType    string `yaml:"dbtype"`    // db type
		DbConnect string `yaml:"dbconnect"` // db connect params string
	} `yaml:"dbsettings"`
	Settings struct {
		Debug bool `yaml:"debug"`
	} `yaml:"settings"`
}

func Get() (*Config, error) {
	conf := &Config{}

	f, err := os.Open("./config.yaml")
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) GetDBType() string {
	return c.DbSettings.DbType
}

func (c *Config) GetDBConnStr() string {
	return c.DbSettings.DbConnect
}

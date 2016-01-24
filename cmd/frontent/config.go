package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config structure
type Config struct {
	Token struct {
		ClientID string `yaml:"clientid"`
	} `yaml:"token"`

	Server struct {
		Listen string `yaml:"listen"`
	} `yaml:"server"`
}

// OpenConfig opens a config file
func OpenConfig(file string) (*Config, error) {
	c := Config{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

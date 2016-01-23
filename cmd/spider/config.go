package main

import (
	"io/ioutil"
	"time"

	"github.com/gophergala2016/blogalert/repositories/rethink"
	"gopkg.in/yaml.v2"
)

// Config structure
type Config struct {
	RethinkDB *rethink.Config `yaml:"rethinkdb"`
	Refresh   time.Duration   `yaml:"refresh"`
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

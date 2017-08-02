package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const fileName = ".flowverk.yaml"

//Config struct
type Config struct {
	JiraURL     string `yaml:"jiraURL"`
	ProjectName string `yaml:"projectName"`
	User        string
	Pass        string
	Transitions struct {
		Todo       string
		InProgress string `yaml:"inprogress"`
		InReview   string `yaml:"inreview"`
		Done       string
	}
}

// GetConfig ...
func GetConfig() *Config {
	config := new(Config)

	confFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(confFile, &config)

	return config
}

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf `yaml:"logger"`
	HTTP     HTTPConf   `yaml:"http"`
	Database DBConf     `yaml:"database"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DBConf struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"`
}

func NewConfig(configPath string) Config {
	cnf := Config{}

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cnf)
	if err != nil {
		panic(err)
	}

	return cnf
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.HTTP.Host, c.HTTP.Port)
}

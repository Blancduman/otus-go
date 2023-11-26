package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger   LoggerConf     `yaml:"logger"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type RabbitMQConfig struct {
	DSN          string `yaml:"dsn"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	Queue        string `yaml:"queue"`
	Key          string `yaml:"key"`
	ConsumerTag  string `yaml:"consumerTag"`
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

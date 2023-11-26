package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger    LoggerConf     `yaml:"logger"`
	Database  DBConf         `yaml:"database"`
	RabbitMQ  RabbitMQConfig `yaml:"rabbitmq"`
	Scheduler SchedulerConf  `yaml:"scheduler"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type DBConf struct {
	URL  string `yaml:"url"`
	Type string `yaml:"type"`
}

type RabbitMQConfig struct {
	DSN          string `yaml:"dsn"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchangeType"`
	Queue        string `yaml:"queue"`
	Key          string `yaml:"key"`
	ConsumerTag  string `yaml:"consumerTag"`
}

type SchedulerConf struct {
	Period int64 `yaml:"period"`
	Mark   int64 `yaml:"mark"`
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

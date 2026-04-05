package main

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	Logger    LoggerConf
	Storage   StorageConf
	RabbitMQ  RabbitMQConf
	Scheduler SchedulerConf
}

type SchedulerConf struct {
	Interval string `yaml:"interval"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type StorageType string

const (
	Memory StorageType = "MEMORY"
	SQL    StorageType = "SQL"
)

type StorageConf struct {
	Type StorageType `yaml:"type"`
	Addr string      `yaml:"addr"`
}

type RabbitMQConf struct {
	Addr     string `yaml:"addr"`
	Queue    string `yaml:"queue"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func NewConfig(configFile string) Config {
	yamlBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	conf := Config{}
	err = yaml.Unmarshal(yamlBytes, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return conf
}

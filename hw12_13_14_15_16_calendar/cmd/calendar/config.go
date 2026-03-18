package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger LoggerConf
	Server ServerConf
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type ServerConf struct {
	Addr string `yaml:"addr"`
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

package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Conf struct {
	DBUsername string `yaml:"db_username"`
	DBPassword string `yaml:"db_password"`
	DBName     string `yaml:"db_name"`
	DBHost     string `yaml:"db_host"`
	DBPort     string `yaml:"db_port"`
	Port       string `yaml:"port"`
}

func GetConf() *Conf {

	cfg := Conf{}

	yamlFile, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &cfg
}

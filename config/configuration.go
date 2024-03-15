package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		SqlitePath       string `yaml:"sqlite_path"`
		PostgresHost     string `yaml:"postgres_host"`
		PostgresPort     int    `yaml:"postgrest_port"`
		PostgresUser     string `yaml:"postgres_user"`
		PostgresPassword string `yaml:"postgres_password"`
		PostgresDBName   string `yaml:"postgres_db_name"`
	} `yaml:"database"`
	Worker struct {
		Identifier string `yaml:"identifier"`
	} `yaml:"worker"`
}

const ConfigPath = "config.yml"

var Configuration Config

func init() {
	file, err := os.Open(ConfigPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(data, &Configuration)
	if err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}

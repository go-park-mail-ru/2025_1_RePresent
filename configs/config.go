package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Host     string `yaml:"HOST"`
	Port     int    `yaml:"PORT"`
	Username string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Dbname   string `yaml:"POSTGRES_DB"`
	Sslmode  string `yaml:"SSLMODE"`
}

type Config struct {
	Database DatabaseConfig `yaml:"connect_database_in_container"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

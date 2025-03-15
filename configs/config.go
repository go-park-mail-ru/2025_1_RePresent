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

type MailConfig struct {
	SmtpServer string `yaml:"SMTP_SERVER"`
	Port       string `yaml:"PORT"`
	Username   string `yaml:"USERNAME"`
	Password   string `yaml:"PASSWORD"`
	Sender     string `yaml:"SENDER"`
}

type Config struct {
	Database DatabaseConfig `yaml:"connect_database_in_container"`
	Email    MailConfig     `yaml:"connect_smtp_server"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Database DatabaseConfig `yaml:"connect_database_in_container"`
		Email    MailConfig     `yaml:"connect_smtp_server"`
	}
	if err := unmarshal(&aux); err != nil {
		return err
	}
	c.Database = aux.Database
	c.Email = aux.Email
	return nil
}

func LoadConfigs(paths ...string) (*Config, error) {
	var cfg Config
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

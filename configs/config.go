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

type AuthRedisConfig struct {
	EndPoint string `yaml:"ENDPOINT"`
	Password string `yaml:"PASSWORD"`
	Database int    `yaml:"DB_NUMBER"`
}

type Config struct {
	Database  DatabaseConfig  `yaml:"database_postgresql"`
	Email     MailConfig      `yaml:"smtp_server"`
	AuthRedis AuthRedisConfig `yaml:"auth_redis"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux struct {
		Database  DatabaseConfig  `yaml:"database_postgresql"`
		Email     MailConfig      `yaml:"smtp_server"`
		AuthRedis AuthRedisConfig `yaml:"auth_redis"`
	}
	if err := unmarshal(&aux); err != nil {
		return err
	}
	c.Database = aux.Database
	c.Email = aux.Email
	c.AuthRedis = aux.AuthRedis
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

package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Dbname   string
	Sslmode  string
}

type MailConfig struct {
	SmtpServer string
	Port       string
	Username   string
	Password   string
	Sender     string
}

type AuthRedisConfig struct {
	EndPoint string
	Password string
	Database int
}

type AttemptRedisConfig struct {
	EndPoint string
	Password string
	Database int
	Attempts int
}

type MinioConfig struct {
	EndPoint       string
	AccessKeyID    string
	SecretAccesKey string
	Token          string
	UseSSL         string
}

type Config struct {
	Database     DatabaseConfig
	Email        MailConfig
	AuthRedis    AuthRedisConfig
	AttemptRedis AttemptRedisConfig
	Minio        MinioConfig
}


func LoadConfig() (Config, error) {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		return Config{}, fmt.Errorf("Error loading .env file")
	}

	config := Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("HOST"),
			Port:     parseEnvInt("PORT"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Dbname:   os.Getenv("POSTGRES_DB"),
			Sslmode:  os.Getenv("SSLMODE"),
		},
		Email: MailConfig{
			SmtpServer: os.Getenv("SMTP_SERVER"),
			Port:       os.Getenv("PORT"),
			Username:   os.Getenv("USERNAME"),
			Password:   os.Getenv("PASSWORD"),
			Sender:     os.Getenv("SENDER"),
		},
		AuthRedis: AuthRedisConfig{
			EndPoint: os.Getenv("ENDPOINT"),
			Password: os.Getenv("PASSWORD"),
			Database: parseEnvInt("DB_NUMBER"),
		},
		AttemptRedis: AttemptRedisConfig{
			EndPoint: os.Getenv("ENDPOINT"),
			Password: os.Getenv("PASSWORD"),
			Database: parseEnvInt("DB_NUMBER"),
			Attempts: parseEnvInt("ATTEMPTS"),
		},
		Minio: MinioConfig{
			EndPoint:       os.Getenv("ENDPOINT"),
			AccessKeyID:    os.Getenv("ACCESS_KEY_ID"),
			SecretAccesKey: os.Getenv("SECRET_ACCESS_KEY"),
			Token:          os.Getenv("TOKEN"),
			UseSSL:         os.Getenv("USE_SSL"),
		},
	}
	return config, nil
}
// func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	var aux struct {
// 		Database     DatabaseConfig     `yaml:"database"`
// 		Email        MailConfig         `yaml:"smtp_server"`
// 		AuthRedis    AuthRedisConfig    `yaml:"auth_redis"`
// 		AttemptRedis AttemptRedisConfig `yaml:"attempt_redis"`
// 		Minio        MinioConfig        `yaml:"object_storage"`
// 	}
// 	if err := unmarshal(&aux); err != nil {
// 		return err
// 	}
// 	c.Database = aux.Database
// 	c.Email = aux.Email
// 	c.AuthRedis = aux.AuthRedis
// 	c.AttemptRedis = aux.AttemptRedis
// 	c.Minio = aux.Minio
// 	return nil
// }

// func LoadConfigs(paths ...string) (*Config, error) {
// 	var cfg Config
// 	var total []byte
// 	for _, path := range paths {
// 		data, err := os.ReadFile(path)
// 		if err != nil {
// 			return nil, err
// 		}
// 		total = append(total, data...)
// 	}

// 	err := yaml.Unmarshal(total, &cfg)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &cfg, nil
// }

// func (d DatabaseConfig) ConnectionString() string {
// 	return fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
// 		d.Host, d.Port, d.Username, d.Password, d.Dbname, d.Sslmode,
// 	)
// }

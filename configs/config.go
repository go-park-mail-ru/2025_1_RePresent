package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	UsernameDefault string
	UsernameBanner  string
	UsernameAuth    string
	Password        string
	Dbname          string
	Sslmode         string
}

type ScyllaConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	SlotKeyspace string
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

type YooConfig struct {
	ShopID        string
	SecretKey     string
	UUID          string // Добавим UUID для идентификации мерчанта в запросах выплат
	AccountNumber string // номер кошелька для автоматических выплат
}

type GigaChatConfig struct {
	AuthKey  string
	ClientID string
}

type Config struct {
	Database     DatabaseConfig
	Email        MailConfig
	AuthRedis    AuthRedisConfig
	AttemptRedis AttemptRedisConfig
	Minio        MinioConfig
	Scylla       ScyllaConfig
	Yoo          YooConfig
	GigaChat     GigaChatConfig
}

func LoadConfigs() (*Config, error) {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		return &Config{}, fmt.Errorf("Error loading .env file")
	}

	config := Config{
		Database: DatabaseConfig{
			Host:            os.Getenv("PSQL_HOST"),
			Port:            parseEnvInt("PSQL_INSIDE_PORT"),
			UsernameDefault: "postgres",
			UsernameBanner:  os.Getenv("PSQL_USER_BANNER"),
			UsernameAuth:    os.Getenv("PSQL_USER_AUTH"),
			Password:        os.Getenv("PSQL_PASSWORD"),
			Dbname:          os.Getenv("PSQL_DB_NAME"),
			Sslmode:         os.Getenv("PSQL_SSLMODE"),
		},
		Email: MailConfig{
			SmtpServer: os.Getenv("SMTP_SERVER"),
			Port:       os.Getenv("SMTP_PORT"),
			Username:   os.Getenv("SMTP_USERNAME"),
			Password:   os.Getenv("SMTP_PASSWORD"),
			Sender:     os.Getenv("SMTP_SENDER"),
		},
		AuthRedis: AuthRedisConfig{
			EndPoint: os.Getenv("REDIS_ENDPOINT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: parseEnvInt("REDIS_DB_NUMBER"),
		},
		AttemptRedis: AttemptRedisConfig{
			EndPoint: os.Getenv("REDIS_ENDPOINT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: parseEnvInt("REDIS_DB_NUMBER"),
			Attempts: parseEnvInt("REDIS_ATTEMPTS"),
		},
		Minio: MinioConfig{
			EndPoint:       os.Getenv("MINIO_ENDPOINT"),
			AccessKeyID:    os.Getenv("MINIO_ACCESS_KEY_ID"),
			SecretAccesKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
			Token:          os.Getenv("MINIO_TOKEN"),
			UseSSL:         os.Getenv("MINIO_USE_SSL"),
		},

		Scylla: ScyllaConfig{
			Host:         os.Getenv("SCYLLA_HOST"),
			Port:         parseEnvInt("SCYLLA_PORT"),
			Username:     os.Getenv("SCYLLA_USERNAME"),
			Password:     os.Getenv("SCYLLA_PASSWORD"),
			SlotKeyspace: os.Getenv("SCYLLA_SLOT_KEYSPACE"),
		},
		Yoo: YooConfig{
			ShopID:        os.Getenv("YOO_SHOP_ID"),
			SecretKey:     os.Getenv("YOO_SECRET_KEY"),
			UUID:          os.Getenv("YOO_UUID"),
			AccountNumber: os.Getenv("YOO_ACCOUNT_NUMBER"), // добавили кошелёк для выплат
		},
		GigaChat: GigaChatConfig{
			AuthKey:  os.Getenv("GIGACHAT_AUTH_KEY"),
			ClientID: os.Getenv("GIGACHAT_CLIENT_ID"),
		},
	}
	return &config, nil
}

func parseEnvInt(key string) int {
	value := os.Getenv(key)
	// Преобразуем строку в int, если нужно. Если ошибка, возвращаем 0
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return 0
	}
	return result
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

func (d DatabaseConfig) ConnectionString(login string) string {
	var username string
	switch login {
	case "auth":
		username = d.UsernameAuth
	case "banner":
		username = d.UsernameBanner
	default:
		username = "postgres"
	}
	conString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, username, d.Password, d.Dbname, d.Sslmode,
	)
	log.Println("Generated connection string:", conString)
	return conString
}

package config

type GigaChatConfig struct {
	AuthKey  string `env:"GIGACHAT_AUTH_KEY" envDefault:""`
	ClientID string `env:"GIGACHAT_CLIENT_ID" envDefault:""`
}

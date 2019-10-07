package config

type ServiceConfiguration struct {
	ListenAddress string `env:"SERVICE_PORT" envDefault:"127.0.0.1:8000"`
}

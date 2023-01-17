package config

type Config struct {
	GrpcServer GrpcServer
}

type GrpcServer struct {
}

func ReadConfig() *Config {
	var cfg Config

	return &cfg
}

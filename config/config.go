package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GroupAddr string `envconfig:"GROUP_ADDRESS"`
}

func Init(c *Config) error {
	return envconfig.Process("", c)
}

package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Port        string     `yaml:"port,omitempty"`
	Timeout     int        `yaml:"timeout,omitempty"`
	Environment string     `yaml:"mode,omitempty"`
	Core        CoreConfig `yaml:"core,omitempty"`
}

type CoreConfig struct {
	BlackFridayDate string `yaml:"black_friday_date,omitempty"`
}

type DiscountConfig struct {
	Host    string `yaml:"host,omitempty"`
	Timeout int    `yaml:"timeout"`
}

type RepositoryConfig struct {
	Host string `yaml:"host,omitempty"`
}

type Config struct {
	Service    ServiceConfig    `yaml:"service,omitempty"`
	Discount   DiscountConfig   `yaml:"discount,omitempty"`
	Repository RepositoryConfig `yaml:"repository,omitempty"`
}

const DEVELOPMENT = "DEVELOPMENT"
const STAGING = "STAGING"
const PRODUCTION = "PRODUCTION"

func NewConfig() Config {
	return Config{
		Service: ServiceConfig{
			Timeout:     1,
			Environment: DEVELOPMENT,
		},
		Discount: DiscountConfig{
			Timeout: 1,
		},
	}
}

func (c *Config) LoadFromEnv(env, path string) error {
	envPath := os.Getenv(env)
	if envPath == "" {
		envPath = path
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed reading config file: %v", err)
	}

	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		return fmt.Errorf("failed loading config file: %v", err)
	}

	return nil
}

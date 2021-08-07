package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type ServiceConfig struct {
	Host        string `yaml:"host"`
	Timeout     int    `yaml:"timeout"`
	Environment string `yaml:"mode"`
}

type DiscountConfig struct {
	Insecure bool   `yaml:"insecure"`
	CertPath string `yaml:"cert_path"`
	KeyPath  string `yaml:"key_path"`
}

type RulesConfig struct {
	BlackFridayDate string `yaml:"black_friday_date"`
}

type Config struct {
	Service  ServiceConfig
	Discount DiscountConfig
	Rules    RulesConfig
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
}

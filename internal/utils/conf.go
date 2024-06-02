package utils

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Redis struct {
		Addr         string `yaml:"addr"`
		Port         int    `yaml:"port"`
		DB           int    `yaml:"db"`
		Password     string `yaml:"password"`
		PoolSize     int    `yaml:"poolSize"`
		MinIdleConns int    `yaml:"minIdleConns"`
	} `yaml:"redis"`
	Email struct {
		Value    string `yaml:"value"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Addr     string `yaml:"addr"`
	} `yaml:"email"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

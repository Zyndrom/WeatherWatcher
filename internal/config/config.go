package config

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	OpenMapKey   string `yaml:"openWeatherApiKey"`
	JwtSecretKey string `yaml:"jwtSecretKey"`
	Postgres     struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"postgres"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
}

func ReadConfig() *Config {
	config := &Config{}

	file, err := os.Open("./config/config.yml")
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	defer file.Close()

	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err = decoder.Decode(config); err != nil {
			logrus.Fatalf(err.Error())
		}
	}

	return config
}

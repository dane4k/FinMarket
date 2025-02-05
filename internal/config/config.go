package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
)

type Server struct {
	Port int `yaml:"port" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Name     string `yaml:"name" env-required:"true"`
}

type Auth struct {
	JWTSecret string `yaml:"jwt_secret" env-required:"true"`
}

type Telegram struct {
	Token string `yaml:"bot_token" env-required:"true"`
}

type Imgur struct {
	DefaultAvatar string `yaml:"default_avatar" env-required:"true"`
	ClientID      string `yaml:"client_id" env-required:"true"`
	ClientSecret  string `yaml:"client_secret" env-required:"true"`
	AccessToken   string `yaml:"access_token" env-required:"true"`
}

type Config struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Auth     Auth     `yaml:"auth"`
	Telegram Telegram `yaml:"telegram"`
	Imgur    Imgur    `yaml:"imgur"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Fatal(err.Error())
		return nil, err
	}

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	return &config, nil
}

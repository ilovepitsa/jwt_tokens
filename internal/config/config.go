package config

import (
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Application struct {
		Name string `yaml:"name"`
	}

	Network struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	Postgres struct {
		User               string        `yaml:"user"`
		Password           string        `yaml:"password"`
		Host               string        `yaml:"host"`
		Port               string        `yaml:"port"`
		DB                 string        `yaml:"database"`
		ConnectionAttempts int           `yaml:"attempts"`
		ConntectTimeout    time.Duration `yaml:"timeout"`
	}

	Tokens struct {
		AccessTTL  time.Duration `yaml:"accessTTL"`
		RefreshTTL time.Duration `yaml:"refreshTTL"`
	}

	Config struct {
		AppSetting       Application `yaml:"app"`
		NetworkSettings  Network     `yaml:"network"`
		PostgresSettings Postgres    `yaml:"postgres"`
		TokensSettings   Tokens      `yaml:"tokens"`
	}
)

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	err := cleanenv.ReadConfig(path.Join("./", configPath), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

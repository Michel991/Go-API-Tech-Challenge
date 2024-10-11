package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	HTTPPort string `mapstructure:"HTTP_PORT"`
}

type DatabaseConfig struct {
	Name     string `mapstructure:"DATABASE_NAME"`
	User     string `mapstructure:"DATABASE_USER"`
	Password string `mapstructure:"DATABASE_PASSWORD"`
	Host     string `mapstructure:"DATABASE_HOST"`
	Port     string `mapstructure:"DATABASE_PORT"`
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// Fallback to environment variables if .env file is not found
	if config.Database.Host == "" {
		config.Database.Host = os.Getenv("DATABASE_HOST")
	}
	if config.Database.Port == "" {
		config.Database.Port = os.Getenv("DATABASE_PORT")
	}
	if config.Database.User == "" {
		config.Database.User = os.Getenv("DATABASE_USER")
	}
	if config.Database.Password == "" {
		config.Database.Password = os.Getenv("DATABASE_PASSWORD")
	}
	if config.Database.Name == "" {
		config.Database.Name = os.Getenv("DATABASE_NAME")
	}
	if config.HTTPPort == "" {
		config.HTTPPort = os.Getenv("HTTP_PORT")
	}

	return &config, nil
}

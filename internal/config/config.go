package config

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ErrReadingConfig = errors.New("error reading config file")
var ErrUnmarshallingConfig = errors.New("error unmarshalling config")
var ErrLoadConfig = errors.New("failed to load configuration")

type RedisConfig struct {
	Address  string `mapstructure:"REDIS_ADDRESS"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DB       string `mapstructure:"POSTGRES_DB"`
	SslMode  string `mapstructure:"POSTGRES_SSL_MODE"`
}

type Config struct {
	RedisConfig    `mapstructure:",squash"`
	PostgresConfig `mapstructure:",squash"`
}

func GetConfig() (*Config, error) {
	start := time.Now()
	cfg, err := loadConfig()
	if err != nil {
		logrus.Errorf("failed to load configuration: %v", err)
		return nil, ErrLoadConfig
	}
	logrus.Infof("successfully initialized the config file in %v", time.Since(start))
	return cfg, nil
}

func loadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			logrus.Warn("config file not found")
		}
		return nil, ErrReadingConfig
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Errorf("error unmarshalling config: %v", err)
		return nil, ErrUnmarshallingConfig
	}

	return &cfg, nil
}

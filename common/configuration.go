package common

import (
	"time"

	"github.com/spf13/viper"
)

type RDBConfig struct {
	Host            string        `mapstructure:"DB_HOST"`
	Port            int           `mapstructure:"DB_PORT"`
	Database        string        `mapstructure:"DB_NAME"`
	Username        string        `mapstructure:"DB_USER"`
	Password        string        `mapstructure:"DB_PASS"`
	MaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONN"`
	MaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONN"`
	ConnMaxLifetime time.Duration `mapstructure:"DB_CONN_LIFETIME"`
}

type AppConfig struct {
	Database RDBConfig
}

func LoadConfig() (*AppConfig, error) {
	var config RDBConfig
	viper.AddConfigPath("/etc/photostorage/")
	viper.AddConfigPath("$HOME/.photostorage")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName("local")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	result := AppConfig{Database: config}
	return &result, nil
}

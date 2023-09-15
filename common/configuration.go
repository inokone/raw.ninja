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

type ImageStoreConfig struct {
	Type string `mapstructure:"IMG_STORE_TYPE"`
	Path string `mapstructure:"IMG_STORE_PATH"`
}

type AppConfig struct {
	Database RDBConfig
	Store    ImageStoreConfig
}

func LoadConfig() (*AppConfig, error) {
	var db RDBConfig
	var is ImageStoreConfig
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/photostorage/")
	viper.AddConfigPath("$HOME/.photostorage")
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&is)
	if err != nil {
		return nil, err
	}
	result := AppConfig{Database: db, Store: is}
	return &result, nil
}

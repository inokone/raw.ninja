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

type AuthConfig struct {
	JWTSecret string `mapstructure:"JWT_SIGN_SECRET"`
	JWTExp    int    `mapstructure:"JWT_EXPIRATION_HOURS"`
	JWTSecure bool   `mapstructure:"JWT_COOKIE_SECURE"`
}

type AppConfig struct {
	Database RDBConfig
	Store    ImageStoreConfig
	Auth     AuthConfig
}

func LoadConfig() (*AppConfig, error) {
	var db RDBConfig
	var is ImageStoreConfig
	var au AuthConfig
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/photostorage/")
	viper.AddConfigPath("$HOME/.photostorage")
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.SetDefault("JWT_COOKIE_SECURE", true)
	viper.SetDefault("JWT_EXPIRATION_HOURS", 24)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	if err = viper.Unmarshal(&db); err != nil {
		return nil, err
	}
	if err = viper.Unmarshal(&is); err != nil {
		return nil, err
	}
	if err = viper.Unmarshal(&au); err != nil {
		return nil, err
	}
	return &AppConfig{Database: db, Store: is, Auth: au}, nil
}

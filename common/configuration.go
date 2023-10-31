package common

import (
	"time"

	"github.com/spf13/viper"
)

// RDBConfig is a configuration of the relational database.
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

// ImageStoreConfig is a configuration of the image store.
type ImageStoreConfig struct {
	Type  string `mapstructure:"IMG_STORE_TYPE"`
	Path  string `mapstructure:"IMG_STORE_PATH"`
	Quota int64  `mapstructure:"IMG_STORE_QUOTA"`
}

// AuthConfig is a configuration of the authentication.
type AuthConfig struct {
	JWTSecret    string `mapstructure:"JWT_SIGN_SECRET"`
	JWTExp       int    `mapstructure:"JWT_EXPIRATION_HOURS"`
	JWTSecure    bool   `mapstructure:"JWT_COOKIE_SECURE"`
	FrontendRoot string `mapstructure:"FRONTEND_ROOT"`
}

// MailConfig is a configuration of e-mail massaging.
type MailConfig struct {
	NoReplyAddress string `mapstructure:"MAIL_NO_REPLY_ADDRESS"`
	SMTPAddress    string `mapstructure:"MAIL_SMTP_ADDRESS"`
	SMTPUser       string `mapstructure:"MAIL_SMTP_USER"`
	SMTPPassword   string `mapstructure:"MAIL_SMTP_PASSWORD"`
	SMTPPort       int    `mapstructure:"MAIL_SMTP_PORT"`
}

// LogConfig is a configuration of the logging.
type LogConfig struct {
	LogLevel  string `mapstructure:"LOG_LEVEL"`
	PrettyLog bool   `mapstructure:"PRETTY_LOG"`
}

// AppConfig is the holder of all configurations for the application
type AppConfig struct {
	Database RDBConfig
	Store    ImageStoreConfig
	Auth     AuthConfig
	Log      LogConfig
	Mail     MailConfig
}

// LoadConfig is a function loading the configuration from app.env file in the runtime directory or environment variables.
// As a fallback `$HOME/.photostorage` directory also can be used for the .evn file.
func LoadConfig() (*AppConfig, error) {
	var db RDBConfig
	var is ImageStoreConfig
	var au AuthConfig
	var lg LogConfig
	var ml MailConfig
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
	if err = viper.Unmarshal(&lg); err != nil {
		return nil, err
	}
	if err = viper.Unmarshal(&ml); err != nil {
		return nil, err
	}
	return &AppConfig{Database: db, Store: is, Auth: au, Log: lg, Mail: ml}, nil
}

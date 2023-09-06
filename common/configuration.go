package common

import (
	"log"
	"os"
	"strconv"
	"time"
)

type RDBConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type AppConfig struct {
	Database RDBConfig
}

func (r *RDBConfig) New() {
	r.Host = getEnv("DB_HOST", "localhost")
	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		port = 5432
		log.Printf("Environment variable DB_PORT (%v) can not be converted to int. Using default value [%v]", os.Getenv("DB_PORT"), port)
	}
	r.Port = port
	r.Database = getEnv("DB_NAME", "photostorage")
	r.Username = os.Getenv("DB_USER")
	r.Password = os.Getenv("DB_PASS")
	duration, err := time.ParseDuration(getEnv("DB_CONN_LIFETIME", "1h"))
	if err != nil {
		duration = time.Hour
		log.Printf("Environment variable DB_CONN_LIFETIME (%v) can not be converted to time. Using default value [%v]", os.Getenv("DB_CONN_LIFETIME"), string(r.ConnMaxLifetime))
	}
	r.ConnMaxLifetime = duration
	count, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONN", "10"))
	if err != nil {
		count = 10
		log.Printf("Environment variable DB_MAX_IDLE_CONN (%v) can not be converted to int. Using default value [%v]", os.Getenv("DB_MAX_IDLE_CONN"), count)
	}
	r.MaxIdleConns = count
	count, err = strconv.Atoi(getEnv("DB_MAX_OPEN_CONN", "100"))
	if err != nil {
		count = 100
		log.Printf("Environment variable DB_MAX_OPEN_CONN (%v) can not be converted to int. Using default value [%v]", os.Getenv("DB_MAX_OPEN_CONN"), count)
	}
	r.MaxOpenConns = count
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (a *AppConfig) New() {
	a.Database.New()
}

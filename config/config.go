package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server            `mapstructure:",squash"`
	Logger            `mapstructure:",squash"`
	Redis             `mapstructure:",squash"`
	UserAndAdminSetup `mapstructure:",squash"`
	DbDriver          string `mapstructure:"DB_DRIVER"`
	DbUrl             string `mapstructure:"DB_URL"`
	BaseApiUrl        string `mapstructure:"BASE_API_URL"`
}

type Server struct {
	ApiPort int    `mapstructure:"API_PORT"`
	Mode    string `mapstructure:"SERVER_MODE"`
}

type Logger struct {
	LogLevel string `mapstructure:"LOG_LEVEL"`
	Encoding string `mapstructure:"LOG_ENCODING"`
}

type Redis struct {
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

type UserAndAdminSetup struct {
	DefaultEmailSendEmail       string        `mapstructure:"DEFAULT_EMAIL_SEND_EMAIL"`
	DefaultServiceNameSendEmail string        `mapstructure:"DEFAULT_SERVICE_NAME_SEND_EMAIL"`
	DefaultTimeSendEmail        time.Duration `mapstructure:"DEFAULT_TIME_SEND_EMAIL"`
}

func Load(path string) *Config {
	var c Config

	v := viper.New()

	v.AddConfigPath(path)
	v.SetConfigName("simple-api-gateway")
	v.SetConfigType("env")
	v.SetConfigFile(".env")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("config.ReadInConfig: %v", err)
	}

	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("config.Load.Unmarshal: %v", err)
	}

	return &c
}
